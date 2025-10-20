// in order to keep the prod package clean
// need to embed this test inside the config package (access on private methods)
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var invalidDefaultConfigsYAML = []string{
	// invalid YAML
	`rules: {::: {"enabled": false}}`,
	// invalid output format
	`rules: {}
output:
  format: abcd
  color: false`,
	// incomplete config,
	`rules: {}`,
}

var invalidConfigsYAML = []string{
	// invalid YAML
	`rules: {::: {"enabled": false}}`,
	// invalid output format
	`output:
  format: abcd
  color: false`,
}

var invalidConfigsJSON = []string{
	// invalid JSON
	`{"rules": {4}}`,
	// invalid output format
	`{"output": {"format": "abcd", "color": false}}`,
}

func resetYamlDefaultData() {
	yamlDefaultData = []byte(`rules: {}`)
}

func TestLoadDefaultConfig(t *testing.T) {
	configData, err := loadConfig()
	if err != nil {
		t.Errorf("loadConfig() error = %v", err)
	}

	if len(configData.Rules) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(configData.Rules))
	}
}

func TestMustLoadConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("mustLoadConfig() did panic on default config")
		}
	}()

	_ = mustLoadConfig()
}

func TestLoadDefaultConfig_Invalid(t *testing.T) {
	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			defer resetYamlDefaultData()
			yamlDefaultData = []byte(invalidConfig)

			_, err := loadConfig()
			if err == nil {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestMustLoadConfig_Invalid(t *testing.T) {
	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			defer resetYamlDefaultData()
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("mustLoadConfig() did not panic on invalid default config")
				}
			}()

			yamlDefaultData = []byte(invalidConfig)

			_ = mustLoadConfig()
		})
	}
}

func TestLoadConfig_OverriddenByFile(t *testing.T) {
	contentYAML := []byte(`rules:
  RULE_1:
    "enabled": false
output:
  format: compact
  color: false
`)
	contentJSON := []byte(`{"rules": {"RULE_1": {"enabled": false}}, "output": {"format": "compact", "color": false}}`)

	want := config{
		Rules:  map[string]RuleConfiguration{"RULE_1": {Enabled: true}},
		Output: OutputConfiguration{Format: "compact", Color: NullableBool{HasValue: true, IsTrue: false}},
	}

	tests := []struct {
		filename string
		content  []byte
	}{
		{
			filename: ".tfcoach.yml",
			content:  contentYAML,
		},
		{
			filename: ".tfcoach.yaml",
			content:  contentYAML,
		},
		{
			filename: ".tfcoach.json",
			content:  contentJSON,
		},
		{
			filename: ".tfcoach",
			content:  contentJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			_ = os.Chdir(dir)
			_ = os.WriteFile(filepath.Join(dir, tt.filename), tt.content, 0644)
			configData, err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}

			if reflect.DeepEqual(configData, want) {
				t.Errorf("Wanted %v, got %v", want, configData)
			}
		})
	}
}

func TestLoadConfig_InvalidRulesOverride(t *testing.T) {
	runTest := func(t *testing.T, fileName string, invalidConfig []byte) {
		dir := t.TempDir()
		_ = os.Chdir(dir)
		_ = os.WriteFile(filepath.Join(dir, fileName), invalidConfig, 0644)
		_, err := loadConfig()
		if err == nil {
			t.Errorf("Expected error, got none")
		}
	}

	for _, fileName := range []string{".tfcoach.yml", ".tfcoach.yaml"} {
		for _, invalidConfig := range invalidConfigsYAML {
			t.Run(fileName+"_"+invalidConfig, func(t *testing.T) {
				runTest(t, fileName, []byte(invalidConfig))
			})
		}
	}

	for _, fileName := range []string{".tfcoach.json", ".tfcoach"} {
		for _, invalidConfig := range invalidConfigsJSON {
			t.Run(fileName+"_"+invalidConfig, func(t *testing.T) {
				runTest(t, fileName, []byte(invalidConfig))
			})
		}
	}
}

func TestGetConfigByRuleId(t *testing.T) {
	content := []byte(`{"rules": {"RULE_1": {"enabled": false, "spec": {"foo":"bar"}}}, "output": {"format": "compact", "color": false}}`)

	tests := []struct {
		ruleID   string
		expected RuleConfiguration
	}{
		{
			ruleID:   "not_found",
			expected: RuleConfiguration{Enabled: true},
		},
		{
			ruleID:   "RULE_1",
			expected: RuleConfiguration{Enabled: false, Spec: map[string]string{"foo": "bar"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.ruleID, func(t *testing.T) {
			dir := t.TempDir()
			_ = os.Chdir(dir)
			_ = os.WriteFile(filepath.Join(dir, ".tfcoach.json"), content, 0644)
			configData, err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}

			configuration = configData
			ruleConfig := GetConfigByRuleID(tt.ruleID)
			if ruleConfig.Enabled != tt.expected.Enabled {
				t.Errorf("Expected %+v, got %+v", tt.expected, ruleConfig)
			}
			if !reflect.DeepEqual(ruleConfig.Spec, tt.expected.Spec) {
				t.Errorf("Expected %+v, got %+v", tt.expected.Spec, ruleConfig.Spec)
			}
		})
	}
}

func TestGetOutputConfiguration(t *testing.T) {
	configCompactFalseYAML := []byte(`output:
  format: compact
  color: false
`)
	configCompactFalseJSON := []byte(`{"output": {"format": "compact", "color": false}}`)

	want := OutputConfiguration{Format: "compact", Color: NullableBool{HasValue: true, IsTrue: false}}

	tests := []struct {
		fileName string
		content  []byte
	}{
		{
			fileName: ".tfcoach.yaml",
			content:  configCompactFalseYAML,
		},
		{
			fileName: ".tfcoach.yml",
			content:  configCompactFalseYAML,
		},
		{
			fileName: ".tfcoach",
			content:  configCompactFalseJSON,
		},
		{
			fileName: ".tfcoach.json",
			content:  configCompactFalseJSON,
		},
	}
	for _, tt := range tests {
		t.Run(tt.fileName, func(t *testing.T) {
			dir := t.TempDir()
			_ = os.Chdir(dir)
			_ = os.WriteFile(filepath.Join(dir, tt.fileName), tt.content, 0644)
			configData, err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}

			configuration = configData
			var got OutputConfiguration
			got = GetOutputConfiguration()
			if got != want {
				t.Errorf("Expected %+v, got %+v", want, got)
			}
		})
	}
}
