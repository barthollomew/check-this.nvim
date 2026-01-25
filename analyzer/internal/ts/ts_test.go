package ts

import "testing"

func TestDetectLanguageTrimsAndLowercases(t *testing.T) {
	got := DetectLanguage(" PYTHON ", "example.js")
	if got != LangPython {
		t.Fatalf("expected %s, got %s", LangPython, got)
	}
}
