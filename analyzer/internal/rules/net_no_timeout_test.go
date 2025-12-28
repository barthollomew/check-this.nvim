package rules

import (
	"testing"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/ts"
)

func TestNetNoTimeoutPython(t *testing.T) {
	src := []byte(`requests.get("https://service")`)
	root, err := ts.Parse("python", src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rule := NewNetNoTimeout()
	diags, err := rule.Run(Context{Language: "python", Root: root, Source: src})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
}

func TestNetNoTimeoutFetch(t *testing.T) {
	src := []byte(`fetch("/api/data")`)
	root, err := ts.Parse("javascript", src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	rule := NewNetNoTimeout()
	diags, err := rule.Run(Context{Language: "javascript", Root: root, Source: src})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if len(diags) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(diags))
	}
}
