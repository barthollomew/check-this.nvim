package rules

import (
	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

// Meta provides metadata about a rule.
type Meta struct {
	DefaultSeverity string
	Tags            []string
	Short           string
	Long            string
}

// Context contains the data a rule needs to perform analysis.
type Context struct {
	Language string
	Root     *sitter.Node
	Source   []byte
}

// Rule represents a single analysis rule.
type Rule interface {
	ID() string
	Meta() Meta
	Supports(language string) bool
	Run(ctx Context) ([]diagnostic.Diagnostic, error)
}

// content returns the string for a node, guarding against nil.
func content(source []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return string(source[n.StartByte():n.EndByte()])
}
