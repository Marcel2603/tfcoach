// in order to keep the prod package clean
// need to embed this test inside the config package (access on private methods)
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type navigatorMock struct {
	tempDir string
}

func (n *navigatorMock) HomeDir() (string, error) {
	return n.tempDir, nil
}

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

func resetYamlDefaultData() {
	yamlDefaultData = []byte(`rules: {}
output:
  format: educational
  color: true
  emojis: true
`)
}

func resetNavigator() {
	navig = &DefaultNavigator{}
}

func TestLoadDefaultConfig(t *testing.T) {
	t.Cleanup(resetNavigator)
	navig = &navigatorMock{tempDir: t.TempDir()}

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
			t.Errorf("MustLoadConfig() did panic on default config")
		}
	}()

	_ = MustLoadConfig()
}

func TestLoadDefaultConfig_Invalid(t *testing.T) {
	t.Cleanup(resetYamlDefaultData)
	t.Cleanup(resetNavigator)
	navig = &navigatorMock{tempDir: t.TempDir()}

	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			yamlDefaultData = []byte(invalidConfig)

			_, err := loadConfig()
			if err == nil {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestMustLoadConfig_Invalid(t *testing.T) {
	t.Cleanup(resetYamlDefaultData)
	t.Cleanup(resetNavigator)
	navig = &navigatorMock{tempDir: t.TempDir()}

	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustLoadConfig() did not panic on invalid default config")
				}
			}()

			yamlDefaultData = []byte(invalidConfig)

			_ = MustLoadConfig()
		})
	}
}

func TestLoadConfig_OverriddenByFile(t *testing.T) {
	t.Cleanup(resetNavigator)
	navig = &navigatorMock{tempDir: t.TempDir()}

	contentYAML := []byte(`rules:
  RULE_1:
    "enabled": false
output:
  format: compact
  color: false
  emojis: false
`)
	contentJSON := []byte(`{"rules": {"RULE_1": {"enabled": false}}, "output": {"format": "compact", "color": false, "emojis": false}}`)

	want := config{
		Rules:  map[string]RuleConfiguration{"RULE_1": {Enabled: false}},
		Output: OutputConfiguration{Format: "compact", Color: NullableBool{HasValue: true, IsTrue: false}, Emojis: NullableBool{HasValue: true, IsTrue: false}},
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

			if !reflect.DeepEqual(configData, want) {
				t.Errorf("Wanted %v, got %v", want, configData)
			}
		})
	}
}

func TestLoadConfig_DoubleOverrideInHomeThenCustom(t *testing.T) {
	t.Cleanup(resetNavigator)
	homeDir := t.TempDir()
	navig = &navigatorMock{tempDir: homeDir}

	contentHomeYAML := []byte(`rules:
  RULE_1:
    enabled: false
output:
  format: compact
  color: false
  emojis: false
`)
	contentHomeJSON := []byte(`{"rules": {"RULE_1": {"enabled": false}}, "output": {"format": "compact", "color": false, "emojis": false}}`)

	contentCustomYAML := []byte(`rules:
  RULE_2:
    enabled: false
output:
  format: pretty
  emojis: true
`)
	contentCustomJSON := []byte(`{"rules": {"RULE_2": {"enabled": false}}, "output": {"format": "pretty", "emojis": true}}`)

	want := config{
		Rules:  map[string]RuleConfiguration{"RULE_1": {Enabled: false}, "RULE_2": {Enabled: false}},
		Output: OutputConfiguration{Format: "pretty", Color: NullableBool{HasValue: true, IsTrue: false}, Emojis: NullableBool{HasValue: true, IsTrue: true}},
	}

	tests := []struct {
		filename      string
		contentHome   []byte
		contentCustom []byte
	}{
		{
			filename:      ".tfcoach.yml",
			contentHome:   contentHomeYAML,
			contentCustom: contentCustomYAML,
		},
		{
			filename:      ".tfcoach.yaml",
			contentHome:   contentHomeYAML,
			contentCustom: contentCustomYAML,
		},
		{
			filename:      ".tfcoach.json",
			contentHome:   contentHomeJSON,
			contentCustom: contentCustomJSON,
		},
		{
			filename:      ".tfcoach",
			contentHome:   contentHomeJSON,
			contentCustom: contentCustomJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			homeConfigDir := filepath.Join(homeDir, ".config", "tfcoach")
			_ = os.MkdirAll(homeConfigDir, 0777)
			err := os.WriteFile(filepath.Join(homeConfigDir, tt.filename), tt.contentHome, 0644)
			dir := t.TempDir()
			_ = os.Chdir(dir)
			_ = os.WriteFile(filepath.Join(dir, tt.filename), tt.contentCustom, 0644)
			configData, err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}

			if !reflect.DeepEqual(configData, want) {
				t.Errorf("Wanted %v, got %v", want, configData)
			}
		})
	}
}

func TestLoadConfig_InvalidOverride(t *testing.T) {
	contentYAML := []byte(`output:
  format: abcd
  color: false`)
	contentJSON := []byte(`{"output": {"format": "abcd", "color": false}}`)

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
			_, err := loadConfig()
			if err == nil {
				t.Errorf("Expected error, got none")
			}
		})
	}
}

func TestLoadConfig_InvalidCustomFileIsIgnored(t *testing.T) {
	t.Cleanup(resetNavigator)
	navig = &navigatorMock{tempDir: t.TempDir()}

	contentYAML := []byte(`rules: {::: {"enabled": false}}`)
	contentJSON := []byte(`{"rules": {4}}`)

	want := config{
		Rules:  make(map[string]RuleConfiguration),
		Output: OutputConfiguration{Format: "educational", Color: NullableBool{HasValue: true, IsTrue: true}, Emojis: NullableBool{HasValue: true, IsTrue: true}},
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

			if !reflect.DeepEqual(configData, want) {
				t.Errorf("Wanted %v, got %v", want, configData)
			}
		})
	}
}

func TestGetConfigByRuleId(t *testing.T) {
	content := []byte(`{"rules": {"RULE_1": {"enabled": false, "spec": {"foo":"bar"}}}, "output": {"format": "compact", "color": false, "emojis": true}}`)

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
  emojis: true
`)
	configCompactFalseJSON := []byte(`{"output": {"format": "compact", "color": false, "emojis": true}}`)

	want := OutputConfiguration{Format: "compact", Color: NullableBool{HasValue: true, IsTrue: false}, Emojis: NullableBool{HasValue: true, IsTrue: true}}

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
