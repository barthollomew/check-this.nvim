package diagnostic

// Position represents a point in the source file. Zero-based indexing is used to
// align with Tree-sitter nodes.
type Position struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

// Range represents a span in the source file.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Diagnostic describes a single finding produced by the analyzer.
type Diagnostic struct {
	RuleID      string   `json:"rule_id"`
	Severity    string   `json:"severity"`
	Message     string   `json:"message"`
	Explanation string   `json:"explanation,omitempty"`
	Range       Range    `json:"range"`
	Tags        []string `json:"tags,omitempty"`
	DocsURL     string   `json:"docs_url,omitempty"`
}

// Stats contains runtime metrics for transparency.
type Stats struct {
	ParseMS   int `json:"parse_ms"`
	AnalyzeMS int `json:"analyze_ms"`
	RulesRun  int `json:"rules_run"`
}

// Output is the JSON envelope emitted by the analyzer.
type Output struct {
	Version     string       `json:"version"`
	Path        string       `json:"path,omitempty"`
	Language    string       `json:"language"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Stats       Stats        `json:"stats"`
}
