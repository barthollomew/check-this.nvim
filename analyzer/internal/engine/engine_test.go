package engine

import (
	"strings"
	"testing"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/config"
)

func TestAnalyzeSuppression(t *testing.T) {
	src := "# check-this: disable=errors.swallowed\ntry:\n    risky()\nexcept Exception:\n    pass\n"
	input := AnalyzeInput{
		Path:    "example.py",
		Lang:    "python",
		Source:  []byte(src),
		Config:  config.Config{},
		Version: "1.0",
	}
	if err := ValidateInput(input); err != nil {
		t.Fatalf("validate: %v", err)
	}
	out, err := NewEngine().Analyze(input)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	for _, d := range out.Diagnostics {
		if strings.EqualFold(d.RuleID, "errors.swallowed") {
			t.Fatalf("expected suppression to drop errors.swallowed")
		}
	}
}

func TestAnalyzeRuleDisabledByConfig(t *testing.T) {
	disabled := false
	src := "try:\n    risky()\nexcept Exception:\n    pass\n"
	input := AnalyzeInput{
		Path:   "example.py",
		Lang:   "python",
		Source: []byte(src),
		Config: config.Config{
			Rules: map[string]config.RuleSetting{
				"errors.swallowed": {Enabled: &disabled},
			},
		},
		Version: "1.0",
	}
	if err := ValidateInput(input); err != nil {
		t.Fatalf("validate: %v", err)
	}
	out, err := NewEngine().Analyze(input)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if len(out.Diagnostics) != 0 {
		t.Fatalf("expected no diagnostics when rule disabled, got %d", len(out.Diagnostics))
	}
}

func TestAnalyzeSeverityOverride(t *testing.T) {
	src := `fetch("/api/data")`
	input := AnalyzeInput{
		Path:   "example.js",
		Lang:   "javascript",
		Source: []byte(src),
		Config: config.Config{
			Rules: map[string]config.RuleSetting{
				"net.no_timeout": {Severity: "error"},
			},
		},
		Version: "1.0",
	}
	if err := ValidateInput(input); err != nil {
		t.Fatalf("validate: %v", err)
	}
	out, err := NewEngine().Analyze(input)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if len(out.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(out.Diagnostics))
	}
	if out.Diagnostics[0].Severity != "error" {
		t.Fatalf("expected severity override to apply, got %s", out.Diagnostics[0].Severity)
	}
}
