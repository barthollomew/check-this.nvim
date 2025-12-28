package rules

import (
	"strings"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

func rangeFromNode(n *sitter.Node) diagnostic.Range {
	if n == nil {
		return diagnostic.Range{
			Start: diagnostic.Position{Line: 0, Col: 0},
			End:   diagnostic.Position{Line: 0, Col: 1},
		}
	}
	start := n.StartPoint()
	end := n.EndPoint()
	return diagnostic.Range{
		Start: diagnostic.Position{Line: int(start.Row), Col: int(start.Column)},
		End:   diagnostic.Position{Line: int(end.Row), Col: int(end.Column)},
	}
}

func matchesAny(s string, options ...string) bool {
	for _, opt := range options {
		if strings.EqualFold(s, opt) {
			return true
		}
	}
	return false
}
