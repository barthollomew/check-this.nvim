package rules

import (
	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

// meta holds rule info.
type Meta struct {
	DefaultSeverity string
	Tags            []string
	Short           string
	Long            string
}

// context holds rule input.
type Context struct {
	Language string
	Root     *sitter.Node
	Source   []byte
}

// rule is one check.
type Rule interface {
	ID() string
	Meta() Meta
	Supports(language string) bool
	Run(ctx Context) ([]diagnostic.Diagnostic, error)
}

// content returns node text, nil-safe.
func content(source []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return string(source[n.StartByte():n.EndByte()])
}
