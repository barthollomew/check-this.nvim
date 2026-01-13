package rules

import (
	"strings"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

type RetryUnbounded struct{}

// newretryunbounded builds rule.
func NewRetryUnbounded() Rule { return RetryUnbounded{} }

func (RetryUnbounded) ID() string { return "retry.unbounded" }

func (RetryUnbounded) Meta() Meta {
	return Meta{
		DefaultSeverity: "warning",
		Tags:            []string{"reliability", "retries"},
		Short:           "Retry loop lacks limits/backoff",
		Long:            "Unbounded retries can overload dependencies during outages.",
	}
}

func (RetryUnbounded) Supports(language string) bool {
	switch strings.ToLower(language) {
	case "python", "javascript", "typescript":
		return true
	}
	return false
}

func (r RetryUnbounded) Run(ctx Context) ([]diagnostic.Diagnostic, error) {
	switch strings.ToLower(ctx.Language) {
	case "python":
		return r.runPython(ctx), nil
	case "javascript", "typescript":
		return r.runJS(ctx), nil
	default:
		return nil, nil
	}
}

func (r RetryUnbounded) runPython(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Type() == "while_statement" {
			cond := n.ChildByFieldName("condition")
			condText := strings.TrimSpace(content(ctx.Source, cond))
			if strings.EqualFold(condText, "true") {
				body := n.ChildByFieldName("body")
				if body == nil {
					body = firstChildOfType(n, "block")
				}
				if body != nil && !hasBackoff(body, ctx.Source) {
					diags = append(diags, diagnostic.Diagnostic{
						RuleID:      r.ID(),
						Message:     "Retry loop without cap or backoff",
						Explanation: "Infinite retries can amplify outages; add max attempts and backoff.",
						Range:       rangeFromNode(n),
					})
				}
			}
		}
		for i := 0; i < int(n.NamedChildCount()); i++ {
			walk(n.NamedChild(i))
		}
	}
	walk(ctx.Root)
	return diags
}

func (r RetryUnbounded) runJS(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Type() == "while_statement" || n.Type() == "for_statement" {
			if isInfiniteLoop(n, ctx.Source) {
				body := n.ChildByFieldName("body")
				if body == nil {
					body = firstChildOfType(n, "statement_block")
				}
				if body != nil && !hasBackoff(body, ctx.Source) {
					diags = append(diags, diagnostic.Diagnostic{
						RuleID:      r.ID(),
						Message:     "Potential unbounded retry loop",
						Explanation: "Add max attempts or backoff to avoid hammering dependencies during failures.",
						Severity:    "warning",
						Range:       rangeFromNode(n),
						Tags:        []string{"reliability", "retries"},
					})
				}
			}
		}
		for i := 0; i < int(n.NamedChildCount()); i++ {
			walk(n.NamedChild(i))
		}
	}
	walk(ctx.Root)
	return diags
}

func hasBackoff(body *sitter.Node, source []byte) bool {
	if body == nil {
		return false
	}
	var found bool
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil || found {
			return
		}
		if n.Type() == "call" || n.Type() == "call_expression" {
			fn := n.ChildByFieldName("function")
			name := strings.ToLower(strings.TrimSpace(content(source, fn)))
			if strings.Contains(name, "sleep") || strings.Contains(name, "backoff") || strings.Contains(name, "delay") || strings.Contains(name, "settimeout") {
				found = true
				return
			}
		}
		if n.Type() == "break_statement" || n.Type() == "return_statement" {
			// break/return usually bound the loop.
			found = true
			return
		}
		for i := 0; i < int(n.NamedChildCount()); i++ {
			walk(n.NamedChild(i))
		}
	}
	walk(body)
	return found
}

func isInfiniteLoop(n *sitter.Node, source []byte) bool {
	switch n.Type() {
	case "while_statement":
		cond := n.ChildByFieldName("condition")
		condText := strings.TrimSpace(strings.ToLower(content(source, cond)))
		return condText == "true"
	case "for_statement":
		// for(;;) in js has nil condition
		cond := n.ChildByFieldName("condition")
		return cond == nil || strings.TrimSpace(content(source, cond)) == ""
	default:
		return false
	}
}
