package config

// RuleSetting captures per-rule configuration such as enablement and severity.
type RuleSetting struct {
	Enabled  *bool  `json:"enabled,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// Config holds analyzer configuration. In v1 the analyzer focuses on
// rule-level toggles and severity overrides.
type Config struct {
	Rules map[string]RuleSetting `json:"rules,omitempty"`
}

// RuleEnabled returns whether a rule should run. Defaults to true when unset.
func (c Config) RuleEnabled(ruleID string) bool {
	if cfg, ok := c.Rules[ruleID]; ok && cfg.Enabled != nil {
		return *cfg.Enabled
	}
	return true
}

// RuleSeverity resolves the severity for a rule, allowing overrides.
func (c Config) RuleSeverity(ruleID, defaultSeverity string) string {
	if cfg, ok := c.Rules[ruleID]; ok && cfg.Severity != "" {
		return cfg.Severity
	}
	return defaultSeverity
}

// Merge applies overrides on top of the current config, returning a new copy.
func (c Config) Merge(override Config) Config {
	out := Config{Rules: map[string]RuleSetting{}}
	for k, v := range c.Rules {
		out.Rules[k] = v
	}
	for k, v := range override.Rules {
		out.Rules[k] = v
	}
	return out
}
