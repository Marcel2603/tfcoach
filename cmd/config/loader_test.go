// in oder to keep the prod package clean
// need to embed this test inside the config package (access on private methods)
package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	loadConfig()
	if len(Configuration.Rules) != 0 {
		t.Errorf("Expected 0 rules, got %d", len(Configuration.Rules))
	}
}

func TestLoadConfigOverriddenByFile(t *testing.T) {
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
		}, {
			filename: ".tfcoach.yaml",
			content:  contentYAML,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		},
		{
			filename: ".tfcoach.json",
			content:  contentJSON,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		}, {
			filename: ".tfcoach",
			content:  contentJSON,
			expected: config{Rules: map[string]RuleConfiguration{"RULE_1": {Enabled: true}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			dir := t.TempDir()
			os.Chdir(dir)
			os.WriteFile(filepath.Join(dir, tt.filename), tt.content, 0644)
			loadConfig()
			if reflect.DeepEqual(Configuration, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, Configuration)
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
			os.Chdir(dir)
			os.WriteFile(filepath.Join(dir, ".tfcoach.json"), content, 0644)
			loadConfig()
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
