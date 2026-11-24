package rules

import (
	"github.com/barthollomew/check-this.nvim/analyzer/internal/diagnostic"
	sitter "github.com/smacker/go-tree-sitter"
)

type Meta struct {
	DefaultSeverity string
	Tags            []string
	Short           string
	Long            string
}

type Context struct {
	Language string
	Root     *sitter.Node
	Source   []byte
}

type Rule interface {
	ID() string
	Meta() Meta
	Supports(language string) bool
	Run(ctx Context) ([]diagnostic.Diagnostic, error)
}

func content(source []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}
	return string(source[n.StartByte():n.EndByte()])
}
