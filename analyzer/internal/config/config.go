package config

// rulesetting holds per-rule config.
type RuleSetting struct {
	Enabled  *bool  `json:"enabled,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// config holds analyzer settings.
type Config struct {
	Rules map[string]RuleSetting `json:"rules,omitempty"`
}

// ruleenabled reports if rule runs.
func (c Config) RuleEnabled(ruleID string) bool {
	if cfg, ok := c.Rules[ruleID]; ok && cfg.Enabled != nil {
		return *cfg.Enabled
	}
	return true
}

// ruleseverity picks rule severity.
func (c Config) RuleSeverity(ruleID, defaultSeverity string) string {
	if cfg, ok := c.Rules[ruleID]; ok && cfg.Severity != "" {
		return cfg.Severity
	}
	return defaultSeverity
}

// merge applies overrides.
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
