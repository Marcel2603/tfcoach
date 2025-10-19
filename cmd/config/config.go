package config

type config struct {
	Rules  map[string]RuleConfiguration `json:"rules" yaml:"rules"`
	Output OutputConfiguration          `json:"output" yaml:"output"`
}

type RuleConfiguration struct {
	Enabled bool              `json:"enabled" yaml:"enabled" default:"true"`
	Spec    map[string]string `json:"spec" yaml:"spec"`
}

type OutputConfiguration struct {
	Format string `json:"format" yaml:"format"`
	Color  bool   `json:"color" yaml:"color"`
}
