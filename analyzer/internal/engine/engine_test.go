package engine

import (
	"strings"
	"testing"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/config"
)

func TestAnalyzeSuppression(t *testing.T) {
	src := "# check-this: disable=errors.swallowed\ntry:\n    risky()\nexcept Exception:\n    pass\n"
	input := AnalyzeInput{
		Path:    "example.py",
		Lang:    "python",
		Source:  []byte(src),
		Config:  config.Config{},
		Version: "1.0",
	}
	if err := ValidateInput(input); err != nil {
		t.Fatalf("validate: %v", err)
	}
	out, err := NewEngine().Analyze(input)
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	for _, d := range out.Diagnostics {
		if strings.EqualFold(d.RuleID, "errors.swallowed") {
			t.Fatalf("expected suppression to drop errors.swallowed")
		}
	}
}
