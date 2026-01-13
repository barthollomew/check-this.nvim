package rules

import (
	"strings"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

type NetNoTimeout struct{}

// newnetnotimeout builds rule.
func NewNetNoTimeout() Rule { return NetNoTimeout{} }

func (NetNoTimeout) ID() string { return "net.no_timeout" }

func (NetNoTimeout) Meta() Meta {
	return Meta{
		DefaultSeverity: "warning",
		Tags:            []string{"reliability", "network"},
		Short:           "Network call without timeout",
		Long:            "Network calls without timeouts can hang and block resources during outages.",
	}
}

func (NetNoTimeout) Supports(language string) bool {
	switch strings.ToLower(language) {
	case "python", "javascript", "typescript":
		return true
	}
	return false
}

func (r NetNoTimeout) Run(ctx Context) ([]diagnostic.Diagnostic, error) {
	switch strings.ToLower(ctx.Language) {
	case "python":
		return r.runPython(ctx), nil
	case "javascript", "typescript":
		return r.runJS(ctx), nil
	default:
		return nil, nil
	}
}

func (r NetNoTimeout) runPython(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Type() == "call" {
			fn := n.ChildByFieldName("function")
			if fn != nil {
				name := strings.TrimSpace(content(ctx.Source, fn))
				if isRequestsFunction(name) && !hasKeywordArgument(n, "timeout", ctx.Source) {
					diags = append(diags, diagnostic.Diagnostic{
						RuleID:      r.ID(),
						Message:     "Network call without timeout",
						Explanation: "HTTP calls should specify a timeout to avoid hanging during partial outages.",
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

func (r NetNoTimeout) runJS(ctx Context) []diagnostic.Diagnostic {
	var diags []diagnostic.Diagnostic
	var walk func(n *sitter.Node)
	walk = func(n *sitter.Node) {
		if n == nil {
			return
		}
		if n.Type() == "call_expression" {
			fn := n.ChildByFieldName("function")
			name := strings.TrimSpace(content(ctx.Source, fn))
			if name == "fetch" && !hasFetchTimeout(n, ctx.Source) {
				diags = append(diags, diagnostic.Diagnostic{
					RuleID:      r.ID(),
					Severity:    "info",
					Message:     "fetch call without AbortController/timeout",
					Explanation: "Provide an AbortController or timeout so fetch calls do not hang indefinitely.",
					Range:       rangeFromNode(n),
					Tags:        []string{"reliability", "network"},
				})
			}
			if strings.HasPrefix(name, "axios") && !argumentContains(n, "timeout", ctx.Source) {
				diags = append(diags, diagnostic.Diagnostic{
					RuleID:      r.ID(),
					Severity:    "info",
					Message:     "axios call without timeout option",
					Explanation: "Set axios timeouts to avoid hanging requests during outages.",
					Range:       rangeFromNode(n),
					Tags:        []string{"reliability", "network"},
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

func isRequestsFunction(name string) bool {
	prefixes := []string{"requests.", "httpx."}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) {
			return true
		}
	}
	return false
}

func hasKeywordArgument(call *sitter.Node, name string, source []byte) bool {
	if call == nil {
		return false
	}
	args := call.ChildByFieldName("arguments")
	if args == nil {
		return false
	}
	for i := 0; i < int(args.NamedChildCount()); i++ {
		arg := args.NamedChild(i)
		if arg == nil {
			continue
		}
		if arg.Type() == "keyword_argument" {
			id := arg.ChildByFieldName("name")
			if id != nil && strings.TrimSpace(content(source, id)) == name {
				return true
			}
			// fallback to text check
			if strings.HasPrefix(strings.TrimSpace(content(source, arg)), name) {
				return true
			}
		}
	}
	return false
}

func hasFetchTimeout(call *sitter.Node, source []byte) bool {
	args := call.ChildByFieldName("arguments")
	if args == nil {
		return false
	}
	if args.NamedChildCount() < 2 {
		return false
	}
	for i := 0; i < int(args.NamedChildCount()); i++ {
		arg := args.NamedChild(i)
		if arg == nil {
			continue
		}
		text := strings.TrimSpace(content(source, arg))
		if strings.Contains(text, "timeout") || strings.Contains(text, "AbortController") || strings.Contains(text, "signal") {
			return true
		}
	}
	return false
}

func argumentContains(call *sitter.Node, needle string, source []byte) bool {
	args := call.ChildByFieldName("arguments")
	if args == nil {
		return false
	}
	for i := 0; i < int(args.NamedChildCount()); i++ {
		arg := args.NamedChild(i)
		if arg == nil {
			continue
		}
		if strings.Contains(strings.ToLower(content(source, arg)), strings.ToLower(needle)) {
			return true
		}
	}
	return false
}
