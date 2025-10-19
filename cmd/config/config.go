package config

import (
	"fmt"
	"slices"
	"strings"
)

var (
	// TODO later: educational
	supportedOutputFormats = []string{"json", "compact", "pretty"}
)

type config struct {
	Rules  map[string]RuleConfiguration `json:"rules" yaml:"rules"`
	Output OutputConfiguration          `json:"output" yaml:"output"`
}

type RuleConfiguration struct {
	Enabled bool              `json:"enabled" yaml:"enabled" default:"true"`
	Spec    map[string]string `json:"spec" yaml:"spec"`
}

type OutputConfiguration struct {
	Format    string `json:"format" yaml:"format"`
	Color     string `json:"color" yaml:"color"`
	colorBool bool
}

func (o *OutputConfiguration) Validate() error {
	if !slices.Contains(supportedOutputFormats, o.Format) {
		return fmt.Errorf("invalid format: %q (supported: %v)", o.Format, supportedOutputFormats)
	}

	parsedColor, err := o.parseColor()
	if err != nil {
		return err
	}
	o.colorBool = parsedColor
	return nil
}

func (o *OutputConfiguration) parseColor() (bool, error) {
	colorStr := strings.ToLower(o.Color)
	if !slices.Contains([]string{"true", "false"}, colorStr) {
		return false, fmt.Errorf("could not parse color %s", o.Color)
	}
	return colorStr == "true", nil
}

func (o *OutputConfiguration) SupportedFormats() []string {
	return slices.Clone(supportedOutputFormats)
}

func (o *OutputConfiguration) ParsedColor() bool {
	return o.colorBool
}
