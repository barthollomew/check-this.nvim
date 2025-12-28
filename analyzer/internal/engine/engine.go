package engine

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/config"
	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	"github.com/barthollomew/check-this.nvim/analyzer/internal/rules"
	"github.com/barthollomew/check-this.nvim/analyzer/internal/ts"
)

// Engine coordinates parsing, rule execution, and suppression handling.
type Engine struct {
	rules []rules.Rule
}

// NewEngine returns an Engine populated with the default ruleset.
func NewEngine() Engine {
	return Engine{
		rules: []rules.Rule{
			rules.NewErrorsSwallowed(),
			rules.NewNetNoTimeout(),
			rules.NewRetryUnbounded(),
			rules.NewStateGlobalMutable(),
		},
	}
}

// AnalyzeInput captures the request to analyze a source buffer.
type AnalyzeInput struct {
	Path    string
	Source  []byte
	Lang    string
	Config  config.Config
	Version string
}

// Analyze runs the enabled rules and returns a JSON-friendly output envelope.
func (e Engine) Analyze(input AnalyzeInput) (diagnostic.Output, error) {
	if strings.TrimSpace(string(input.Source)) == "" {
		return diagnostic.Output{
			Version:     versionOrDefault(input.Version),
			Path:        input.Path,
			Language:    input.Lang,
			Diagnostics: nil,
			Stats:       diagnostic.Stats{},
		}, nil
	}

	startParse := time.Now()
	root, parseErr := ts.Parse(input.Lang, input.Source)
	parseMS := int(time.Since(startParse).Milliseconds())

	out := diagnostic.Output{
		Version:  versionOrDefault(input.Version),
		Path:     input.Path,
		Language: input.Lang,
		Stats: diagnostic.Stats{
			ParseMS: parseMS,
		},
	}

	if parseErr != nil {
		out.Diagnostics = append(out.Diagnostics, diagnostic.Diagnostic{
			RuleID:      "internal.parse_error",
			Severity:    "error",
			Message:     fmt.Sprintf("parse error: %v", parseErr),
			Explanation: "The analyzer could not parse this file; results may be incomplete.",
			Range: diagnostic.Range{
				Start: diagnostic.Position{Line: 0, Col: 0},
				End:   diagnostic.Position{Line: 0, Col: 1},
			},
			Tags: []string{"internal"},
		})
		return out, nil
	}

	suppressions := collectSuppressions(input.Source)
	analyzeStart := time.Now()
	for _, rule := range e.rules {
		if !input.Config.RuleEnabled(rule.ID()) {
			continue
		}
		if !rule.Supports(input.Lang) {
			continue
		}

		diags, err := runRule(rule, rules.Context{
			Language: input.Lang,
			Root:     root,
			Source:   input.Source,
		})
		if err != nil {
			out.Diagnostics = append(out.Diagnostics, diagnostic.Diagnostic{
				RuleID:   fmt.Sprintf("internal.%s", rule.ID()),
				Severity: "error",
				Message:  fmt.Sprintf("rule %s failed: %v", rule.ID(), err),
				Range: diagnostic.Range{
					Start: diagnostic.Position{Line: 0, Col: 0},
					End:   diagnostic.Position{Line: 0, Col: 1},
				},
				Tags: []string{"internal"},
			})
			continue
		}

		for _, d := range diags {
			if shouldSuppress(suppressions, d) {
				continue
			}
			if d.RuleID == "" {
				d.RuleID = rule.ID()
			}
			if d.Severity == "" {
				d.Severity = input.Config.RuleSeverity(rule.ID(), rule.Meta().DefaultSeverity)
			}
			if len(d.Tags) == 0 {
				d.Tags = rule.Meta().Tags
			}
			out.Diagnostics = append(out.Diagnostics, d)
		}
		out.Stats.RulesRun++
	}
	out.Stats.AnalyzeMS = int(time.Since(analyzeStart).Milliseconds())
	return out, nil
}

func runRule(rule rules.Rule, ctx rules.Context) ([]diagnostic.Diagnostic, error) {
	diags, err := rule.Run(ctx)
	if err != nil {
		return nil, err
	}
	return diags, nil
}

func shouldSuppress(s map[string]map[int]struct{}, d diagnostic.Diagnostic) bool {
	lines, ok := s[d.RuleID]
	if ok {
		if _, exists := lines[d.Range.Start.Line]; exists {
			return true
		}
		// Presence of the rule in the map indicates a file-wide disable as well.
		if _, exists := lines[-1]; exists {
			return true
		}
	}
	return false
}

func collectSuppressions(source []byte) map[string]map[int]struct{} {
	out := map[string]map[int]struct{}{}
	lines := strings.Split(string(source), "\n")
	for i, line := range lines {
		commentIdx := strings.Index(line, "check-this: disable=")
		if commentIdx == -1 {
			continue
		}
		after := line[commentIdx+len("check-this: disable="):]
		for _, part := range strings.Split(after, ",") {
			ruleID := strings.TrimSpace(part)
			if ruleID == "" {
				continue
			}
			if out[ruleID] == nil {
				out[ruleID] = map[int]struct{}{}
			}
			out[ruleID][i] = struct{}{}
			// Sentinel -1 marks a file-wide suppression.
			out[ruleID][-1] = struct{}{}
		}
	}
	return out
}

func versionOrDefault(v string) string {
	if strings.TrimSpace(v) == "" {
		return "1.0"
	}
	return v
}

// ValidateInput ensures the request contains the minimum fields needed.
func ValidateInput(input AnalyzeInput) error {
	if len(input.Source) == 0 {
		return errors.New("source is empty")
	}
	if input.Lang == "" {
		return errors.New("language is required")
	}
	if !ts.Supported(input.Lang) {
		return fmt.Errorf("language %s not supported", input.Lang)
	}
	return nil
}
