package rules

import (
	"testing"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/ts"
)

func TestErrorsSwallowedPythonPass(t *testing.T) {
	src := []byte(`
try:
    risky()
except Exception:
    pass
`)
	root, err := ts.Parse("python", src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rule := NewErrorsSwallowed()
	diags, err := rule.Run(Context{Language: "python", Root: root, Source: src})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
}

func TestErrorsSwallowedJSEmptyCatch(t *testing.T) {
	src := []byte(`
try {
  risky()
} catch (e) {}
`)
	root, err := ts.Parse("javascript", src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rule := NewErrorsSwallowed()
	diags, err := rule.Run(Context{Language: "javascript", Root: root, Source: src})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
}
