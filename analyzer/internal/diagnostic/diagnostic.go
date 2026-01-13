package diagnostic

// position is a zero-based point in source.
type Position struct {
	Line int `json:"line"`
	Col  int `json:"col"`
}

// range is a span in source.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// diagnostic is one finding.
type Diagnostic struct {
	RuleID      string   `json:"rule_id"`
	Severity    string   `json:"severity"`
	Message     string   `json:"message"`
	Explanation string   `json:"explanation,omitempty"`
	Range       Range    `json:"range"`
	Tags        []string `json:"tags,omitempty"`
	DocsURL     string   `json:"docs_url,omitempty"`
}

// stats holds runtime metrics.
type Stats struct {
	ParseMS   int `json:"parse_ms"`
	AnalyzeMS int `json:"analyze_ms"`
	RulesRun  int `json:"rules_run"`
}

// output is the json envelope.
type Output struct {
	Version     string       `json:"version"`
	Path        string       `json:"path,omitempty"`
	Language    string       `json:"language"`
	Diagnostics []Diagnostic `json:"diagnostics"`
	Stats       Stats        `json:"stats"`
}
