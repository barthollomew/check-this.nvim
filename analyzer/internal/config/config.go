package config

type RuleSetting struct {
	Enabled  *bool  `json:"enabled,omitempty"`
	Severity string `json:"severity,omitempty"`
}

type Config struct {
	Rules map[string]RuleSetting `json:"rules,omitempty"`
}

func (c Config) RuleEnabled(ruleID string) bool {
	if cfg, ok := c.Rules[ruleID]; ok && cfg.Enabled != nil {
		return *cfg.Enabled
	}
	return true
}

func (c Config) RuleSeverity(ruleID, defaultSeverity string) string {
	if cfg, ok := c.Rules[ruleID]; ok && cfg.Severity != "" {
		return cfg.Severity
	}
	return defaultSeverity
}
