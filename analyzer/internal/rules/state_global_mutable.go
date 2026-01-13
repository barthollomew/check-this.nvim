package rules

import (
	"strings"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

type StateGlobalMutable struct{}

// newstateglobalmutable builds rule.
func NewStateGlobalMutable() Rule { return StateGlobalMutable{} }

func (StateGlobalMutable) ID() string { return "state.global_mutable" }

func (StateGlobalMutable) Meta() Meta {
	return Meta{
		DefaultSeverity: "info",
		Tags:            []string{"state"},
		Short:           "Global mutable state",
		Long:            "Global mutable state can lead to hidden coupling and race conditions.",
	}
}

func (StateGlobalMutable) Supports(language string) bool {
	switch strings.ToLower(language) {
	case "python", "javascript", "typescript":
		return true
	}
	return false
}

func (r StateGlobalMutable) Run(ctx Context) ([]diagnostic.Diagnostic, error) {
	switch strings.ToLower(ctx.Language) {
	case "python":
		return r.runPython(ctx), nil
	case "javascript", "typescript":
		return r.runJS(ctx), nil
	default:
		return nil, nil
	}
}

func (r StateGlobalMutable) runPython(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	for i := 0; i < int(ctx.Root.NamedChildCount()); i++ {
		stmt := ctx.Root.NamedChild(i)
		if stmt == nil {
			continue
		}
		if stmt.Type() == "expression_statement" {
			text := strings.TrimSpace(content(ctx.Source, stmt))
			if strings.Contains(text, "=") && (strings.Contains(text, "[") || strings.Contains(text, "{")) {
				diags = append(diags, diagnostic.Diagnostic{
					RuleID:      r.ID(),
					Message:     "Module-level mutable state",
					Explanation: "Global mutable collections can be shared implicitly across imports. Consider scoping within functions or using immutables.",
					Severity:    "info",
					Range:       rangeFromNode(stmt),
				})
			}
		}
	}
	return diags
}

func (r StateGlobalMutable) runJS(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	for i := 0; i < int(ctx.Root.NamedChildCount()); i++ {
		stmt := ctx.Root.NamedChild(i)
		if stmt == nil {
			continue
		}
		switch stmt.Type() {
		case "lexical_declaration", "variable_declaration":
			for j := 0; j < int(stmt.NamedChildCount()); j++ {
				decl := stmt.NamedChild(j)
				if decl == nil {
					continue
				}
				init := decl.ChildByFieldName("value")
				if init != nil && isMutableLiteral(init, ctx.Source) {
					diags = append(diags, diagnostic.Diagnostic{
						RuleID:      r.ID(),
						Message:     "Module-level mutable state",
						Explanation: "Globals that hold mutable objects are easily shared across imports; prefer local scopes or factories.",
						Severity:    "info",
						Range:       rangeFromNode(decl),
						Tags:        []string{"state"},
					})
				}
			}
		}
	}
	return diags
}

func isMutableLiteral(n *sitter.Node, source []byte) bool {
	t := n.Type()
	if t == "object" || t == "array" || t == "dictionary" || t == "list" {
		return true
	}
	text := strings.TrimSpace(content(source, n))
	return strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[")
}
