package config

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// ship the default config with the app
//
//go:embed .tfcoach.default.yml
var yamlDefaultData []byte
var Configuration config

// load config automatically
func init() {
	loadConfig()
}

func GetConfigByRuleId(ruleID string) RuleConfiguration {
	ruleConfiguration, ok := Configuration.Rules[ruleID]

	if ok {
		return ruleConfiguration
	}
	return RuleConfiguration{Enabled: true}
}

func loadConfig() {
	var defaultData config

	loadConfigFromYaml(yamlDefaultData, &defaultData)

	customConfigPath, found := getCustomConfigPath()
	if found == nil {
		appData, err := loadCustomConfigFromFile(customConfigPath)
		if err == nil {
			mergo.Merge(&defaultData, appData, mergo.WithOverride)
		}
	}

	var envData config
	loadConfigFromEnv(&envData)

	mergo.Merge(defaultData, envData, mergo.WithOverride)
	Configuration = defaultData
}

func loadCustomConfigFromFile(configPath string) (config, error) {
	var appData config
	extension := filepath.Ext(configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return appData, err
	}
	if extension == ".tfcoach" || extension == ".json" {
		loadConfigFromYaml(configData, &appData)
		return appData, nil
	}
	if extension == ".yaml" || extension == ".yml" {
		loadConfigFromYaml(configData, &appData)
		return appData, nil
	}
	return appData, os.ErrNotExist

}

func loadConfigFromEnv(mapData *config) {
	err := envconfig.Process("", mapData)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func loadConfigFromYaml(data []byte, mapData *config) {
	err := yaml.Unmarshal(data, &mapData)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func loadConfigFromJson(data []byte, mapData *config) {
	err := json.Unmarshal(data, &mapData)
	if err != nil {
		log.Fatalf("error: %v", err)
		return
	}
}

func getCustomConfigPath() (string, error) {
	files := []string{
		".tfcoach.yml",
		".tfcoach.yaml",
		".tfcoach.json",
		".tfcoach",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			return f, nil
		}
	}

	return "", os.ErrNotExist
}
