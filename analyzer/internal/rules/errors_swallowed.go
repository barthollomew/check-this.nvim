package rules

import (
	"strings"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

type ErrorsSwallowed struct{}

// newerrorsswallowed builds rule.
func NewErrorsSwallowed() Rule { return ErrorsSwallowed{} }

func (ErrorsSwallowed) ID() string { return "errors.swallowed" }

func (ErrorsSwallowed) Meta() Meta {
	return Meta{
		DefaultSeverity: "warning",
		Tags:            []string{"reliability", "errors"},
		Short:           "Exceptions caught and ignored",
		Long:            "Empty error handlers swallow failures and hide outages.",
	}
}

func (ErrorsSwallowed) Supports(language string) bool {
	switch strings.ToLower(language) {
	case "python", "javascript", "typescript":
		return true
	}
	return false
}

func (r ErrorsSwallowed) Run(ctx Context) ([]diagnostic.Diagnostic, error) {
	switch strings.ToLower(ctx.Language) {
	case "python":
		return r.runPython(ctx), nil
	case "javascript", "typescript":
		return r.runJS(ctx), nil
	default:
		return nil, nil
	}
}

func (r ErrorsSwallowed) runPython(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Type() == "except_clause" {
			block := firstChildOfType(n, "block")
			if block == nil {
				block = firstChildOfType(n, "suite")
			}
			if block == nil || isEmptyBlock(block) || isPassOnly(block) {
				diags = append(diags, diagnostic.Diagnostic{
					RuleID:      r.ID(),
					Message:     "Exception handled but nothing done",
					Explanation: "Swallowing exceptions makes outages harder to detect; log or re-raise instead.",
					Range:       rangeFromNode(n),
				})
			}
		}
		for i := 0; i < int(n.NamedChildCount()); i++ {
			walk(n.NamedChild(i))
		}
	}
	walk(ctx.Root)
	return diags
}

func (r ErrorsSwallowed) runJS(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Type() == "catch_clause" {
			body := n.ChildByFieldName("body")
			if body == nil {
				body = firstChildOfType(n, "statement_block")
			}
			if body == nil || body.NamedChildCount() == 0 || isEmptyBlock(body) {
				diags = append(diags, diagnostic.Diagnostic{
					RuleID:      r.ID(),
					Message:     "Empty catch block swallows errors",
					Explanation: "Unhandled errors disappear silently; handle or log the failure path.",
					Range:       rangeFromNode(n),
					Tags:        []string{"reliability", "errors"},
				})
			}
		}
		for i := 0; i < int(n.NamedChildCount()); i++ {
			walk(n.NamedChild(i))
		}
	}
	walk(ctx.Root)
	return diags
}

func isEmptyBlock(block *sitter.Node) bool {
	return block != nil && block.NamedChildCount() == 0
}

func isPassOnly(block *sitter.Node) bool {
	if block == nil || block.NamedChildCount() != 1 {
		return false
	}
	child := block.NamedChild(0)
	return child != nil && child.Type() == "pass_statement"
}

func firstChildOfType(n *sitter.Node, t string) *sitter.Node {
	if n == nil {
		return nil
	}
	for i := 0; i < int(n.NamedChildCount()); i++ {
		child := n.NamedChild(i)
		if child != nil && child.Type() == t {
			return child
		}
	}
	return nil
}
