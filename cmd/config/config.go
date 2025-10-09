package config

type config struct {
	Rules map[string]RuleConfiguration `json:"rules" yaml:"rules"`
}

type RuleConfiguration struct {
	Enabled bool              `json:"enabled" yaml:"enabled" default:"true"`
	Spec    map[string]string `json:"spec" yaml:"spec"`
}
