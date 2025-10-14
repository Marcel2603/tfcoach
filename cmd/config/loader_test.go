// in order to keep the prod package clean
// need to embed this test inside the config package (access on private methods)
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func resetYamlDefaultData() {
	yamlDefaultData = []byte(`rules: {}`)
}

func TestLoadDefaultConfig(t *testing.T) {
	err := loadConfig()
	if err != nil {
		t.Errorf("loadConfig() error = %v", err)
	}
	if len(Configuration.Rules) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(Configuration.Rules))
	}
}

func TestLoadDefaultConfig_Invalid(t *testing.T) {
	defer resetYamlDefaultData()
	yamlDefaultData = []byte(`rules: {::: {"enabled": false}}`)

	err := loadConfig()
	if err == nil {
		t.Errorf("expected error, got none")
	}
}

func TestLoadConfig_OverriddenByFile(t *testing.T) {
	contentYAML := []byte(`rules: {"RULE_1": {"enabled": false}}`)
	contentJSON := []byte(`{"rules": {"RULE_1": {"enabled": false}}}`)

	tests := []struct {
		filename string
		content  []byte
		expected config
	}{
		{
			filename: ".tfcoach.yml",
			content:  contentYAML,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		},
		{
			filename: ".tfcoach.yaml",
			content:  contentYAML,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		},
		{
			filename: ".tfcoach.json",
			content:  contentJSON,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		},
		{
			filename: ".tfcoach",
			content:  contentJSON,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			_ = os.Chdir(dir)
			_ = os.WriteFile(filepath.Join(dir, tt.filename), tt.content, 0644)
			err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}
			if reflect.DeepEqual(Configuration, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, Configuration)
			}
		})
	}
}

func TestLoadConfig_InvalidOverride(t *testing.T) {
	contentYAML := []byte(`rules: {::: {"enabled": false}}`)
	contentJSON := []byte(`{"rules": {4}}`)
	expected := config{}

	tests := []struct {
		filename string
		content  []byte
		expected config
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
			err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}
			if reflect.DeepEqual(Configuration, expected) {
				t.Errorf("Expected %v, got %v", expected, Configuration)
			}
		})
	}
}

func TestGetConfigByRuleId(t *testing.T) {
	content := []byte(`{"rules": {"RULE_1": {"enabled": false, "spec": {"foo":"bar"}}}}`)

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
			err := loadConfig()
			if err != nil {
				t.Errorf("loadConfig() error = %v", err)
			}
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
