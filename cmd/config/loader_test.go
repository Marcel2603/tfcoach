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
	homeDir          string
	customConfigPath string
}

func (n *navigatorMock) GetHomeDir() (string, error) {
	return n.homeDir, nil
}

func (n *navigatorMock) GetCustomConfigPath() (string, error) {
	if n.customConfigPath != "" {
		return n.customConfigPath, nil
	}
	return os.Getwd()
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

func TestLoadConfig_NoOverrides(t *testing.T) {
	configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
	if err != nil {
		t.Errorf("loadConfig() error = %v", err)
	}

	if len(configData.Rules) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(configData.Rules))
	}
}

func TestMustLoadDefaultConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustLoadDefaultConfig() did panic on default config")
		}
	}()

	MustLoadDefaultConfig()
}

func TestMustLoadConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustLoadConfig() did panic on default config")
		}
	}()

	MustLoadConfig(&navigatorMock{homeDir: t.TempDir()})
}

func TestLoadConfig_InvalidDefaultConfig(t *testing.T) {
	t.Cleanup(resetYamlDefaultData)

	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			yamlDefaultData = []byte(invalidConfig)

			_, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
			if err == nil {
				t.Errorf("expected error, got none")
			}
		})
	}
}

func TestMustLoadDefaultConfig_Invalid(t *testing.T) {
	t.Cleanup(resetYamlDefaultData)

	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustLoadConfig() did not panic on invalid default config")
				}
			}()

			yamlDefaultData = []byte(invalidConfig)

			MustLoadDefaultConfig()
		})
	}
}

func TestMustLoadConfig_InvalidDefaultConfig(t *testing.T) {
	t.Cleanup(resetYamlDefaultData)

	for _, invalidConfig := range invalidDefaultConfigsYAML {
		t.Run(invalidConfig, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("MustLoadConfig() did not panic on invalid default config")
				}
			}()

			yamlDefaultData = []byte(invalidConfig)

			MustLoadConfig(&navigatorMock{homeDir: t.TempDir()})
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
			configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
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
			homeDir := t.TempDir()
			homeConfigDir := filepath.Join(homeDir, ".config", "tfcoach")
			_ = os.MkdirAll(homeConfigDir, 0777)
			err := os.WriteFile(filepath.Join(homeConfigDir, tt.filename), tt.contentHome, 0644)
			dir := t.TempDir()
			_ = os.Chdir(dir)
			_ = os.WriteFile(filepath.Join(dir, tt.filename), tt.contentCustom, 0644)
			configData, err := loadConfig(&navigatorMock{homeDir: homeDir})
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}

			if !reflect.DeepEqual(configData, want) {
				t.Errorf("Wanted %v, got %v", want, configData)
			}
		})
	}
}

func TestLoadConfig_OverriddenByFileInNonStandardLocation(t *testing.T) {
	contentIgnoredYAML := []byte(`rules:
  RULE_1:
    enabled: false
output:
  format: compact
  emojis: false
`)
	contentIgnoredJSON := []byte(`{"rules": {"RULE_1": {"enabled": false}}, "output": {"format": "compact", "emojis": false}}`)

	contentCustomYAML := []byte(`rules:
  RULE_2:
    enabled: false
output:
  format: pretty
  color: false
`)
	contentCustomJSON := []byte(`{"rules": {"RULE_2": {"enabled": false}}, "output": {"format": "pretty", "color": false}}`)

	want := config{
		Rules:  map[string]RuleConfiguration{"RULE_2": {Enabled: false}},
		Output: OutputConfiguration{Format: "pretty", Color: NullableBool{HasValue: true, IsTrue: false}, Emojis: NullableBool{HasValue: true, IsTrue: true}},
	}

	tests := []struct {
		filename       string
		contentIgnored []byte
		contentCustom  []byte
	}{
		{
			filename:       ".tfcoach.yml",
			contentIgnored: contentIgnoredYAML,
			contentCustom:  contentCustomYAML,
		},
		{
			filename:       ".tfcoach.yaml",
			contentIgnored: contentIgnoredYAML,
			contentCustom:  contentCustomYAML,
		},
		{
			filename:       ".tfcoach.json",
			contentIgnored: contentIgnoredJSON,
			contentCustom:  contentCustomJSON,
		},
		{
			filename:       ".tfcoach",
			contentIgnored: contentIgnoredJSON,
			contentCustom:  contentCustomJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dirIgnored := t.TempDir()
			_ = os.Chdir(dirIgnored)
			_ = os.WriteFile(filepath.Join(dirIgnored, tt.filename), tt.contentIgnored, 0644)

			dirCustom := t.TempDir()
			_ = os.WriteFile(filepath.Join(dirCustom, tt.filename), tt.contentCustom, 0644)

			configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir(), customConfigPath: dirCustom})
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}

			if !reflect.DeepEqual(configData, want) {
				t.Errorf("Wanted %v, got %v", want, configData)
			}
		})
	}
}

func TestLoadConfig_OverriddenByFileWithNonStandardName(t *testing.T) {
	contentYAML := []byte(`rules:
  RULE_1:
    enabled: false
output:
  format: pretty
  color: false
`)
	contentJSON := []byte(`{"rules": {"RULE_1": {"enabled": false}}, "output": {"format": "pretty", "color": false}}`)

	want := config{
		Rules:  map[string]RuleConfiguration{"RULE_1": {Enabled: false}},
		Output: OutputConfiguration{Format: "pretty", Color: NullableBool{HasValue: true, IsTrue: false}, Emojis: NullableBool{HasValue: true, IsTrue: true}},
	}

	tests := []struct {
		filename string
		content  []byte
	}{
		{
			filename: "my-tfcoach-config.yml",
			content:  contentYAML,
		},
		{
			filename: "whatever.yaml",
			content:  contentYAML,
		},
		{
			filename: "test.json",
			content:  contentJSON,
		},
		{
			filename: "config.tfcoach",
			content:  contentJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dirCustom := t.TempDir()
			_ = os.WriteFile(filepath.Join(dirCustom, tt.filename), tt.content, 0644)

			configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir(), customConfigPath: filepath.Join(dirCustom, tt.filename)})
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
			_, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
			if err == nil {
				t.Errorf("Expected error, got none")
			}
		})
	}
}

func TestLoadConfig_InvalidCustomFileIsIgnored(t *testing.T) {
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
			configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
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
			configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
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
			configData, err := loadConfig(&navigatorMock{homeDir: t.TempDir()})
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
