package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"

	"gopkg.in/yaml.v3"
)

var (
	supportedOutputFormats = []string{"json", "compact", "pretty", "educational"}
)

type NullableBool struct {
	IsTrue   bool
	HasValue bool
}

type config struct {
	Rules  map[string]RuleConfiguration `json:"rules" yaml:"rules"`
	Output OutputConfiguration          `json:"output" yaml:"output"`
}

type RuleConfiguration struct {
	Enabled bool              `json:"enabled" yaml:"enabled" default:"true"`
	Spec    map[string]string `json:"spec" yaml:"spec"`
}

type OutputConfiguration struct {
	Format string       `json:"format" yaml:"format"`
	Color  NullableBool `json:"color" yaml:"color"`
	Emojis NullableBool `json:"emojis" yaml:"emojis"`
}

func (c *config) Validate() error {
	var errs []error
	if !slices.Contains(supportedOutputFormats, c.Output.Format) {
		errs = append(errs, fmt.Errorf("invalid format: %q (supported: %v)", c.Output.Format, supportedOutputFormats))
	}

	if !c.Output.Color.HasValue {
		errs = append(errs, fmt.Errorf("invalid color: never set"))
	}

	if !c.Output.Emojis.HasValue {
		errs = append(errs, fmt.Errorf("invalid emojis config: never set"))
	}

	return errors.Join(errs...)
}

func (o *OutputConfiguration) SupportedFormats() []string {
	return slices.Clone(supportedOutputFormats)
}

func (nullableBool *NullableBool) UnmarshalJSON(b []byte) error {
	var unmarshalledJSON bool

	if err := json.Unmarshal(b, &unmarshalledJSON); err != nil {
		return err
	}

	nullableBool.IsTrue = unmarshalledJSON
	nullableBool.HasValue = true

	return nil
}

func (nullableBool *NullableBool) UnmarshalYAML(value *yaml.Node) error {
	var unmarshalledYAML bool

	if err := yaml.Unmarshal([]byte(value.Value), &unmarshalledYAML); err != nil {
		return err
	}

	nullableBool.IsTrue = unmarshalledYAML
	nullableBool.HasValue = true

	return nil
}

type NullableBoolTransformer struct{}

func (t NullableBoolTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(NullableBool{}) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				hasValue := src.FieldByName("HasValue").Bool()
				if hasValue {
					dst.Set(src)
				}
			}
			return nil
		}
	}
	return nil
}
