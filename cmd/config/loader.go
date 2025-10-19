package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

var (
	// ship the default config with the app
	//
	//go:embed .tfcoach.default.yml
	yamlDefaultData []byte

	configuration = mustLoadConfig()
)

func GetConfigByRuleID(ruleID string) RuleConfiguration {
	ruleConfiguration, ok := configuration.Rules[ruleID]

	if ok {
		return ruleConfiguration
	}
	return RuleConfiguration{Enabled: true}
}

func GetOutputConfiguration() OutputConfiguration {
	return configuration.Output
}

func mustLoadConfig() config {
	configData, err := loadConfig()
	if err != nil {
		panic("Could not load config: " + err.Error())
	}
	return configData
}

func loadConfig() (config, error) {
	var configData config
	err := loadConfigFromYaml(yamlDefaultData, &configData)
	if err != nil {
		return config{}, err
	}

	customConfigPath, found := getCustomConfigPath()
	var appData config
	if found {
		appData, err = loadCustomConfigFromFile(customConfigPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Could not load config from custom config file %s: %s\n", customConfigPath, err.Error())
		} else {
			mergeErr := mergo.Merge(&configData, appData, mergo.WithOverride, mergo.WithTransformers(NullableBoolTransformer{}))
			if mergeErr != nil {
				return config{}, mergeErr
			}
		}
	}

	var envData config
	err = loadConfigFromEnv(&envData)
	if err != nil {
		return config{}, err
	}
	mergeErr := mergo.Merge(&configData, envData, mergo.WithOverride, mergo.WithTransformers(NullableBoolTransformer{}))
	if mergeErr != nil {
		return config{}, mergeErr
	}

	validationErr := configData.Output.Validate()
	if validationErr != nil {
		return config{}, validationErr
	}

	return configData, nil
}

func loadCustomConfigFromFile(configPath string) (config, error) {
	var appData config
	extension := filepath.Ext(configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return appData, err
	}
	switch extension {
	case ".tfcoach", ".json":
		err = loadConfigFromJSON(configData, &appData)
		return appData, err
	case ".yaml", ".yml":
		err = loadConfigFromYaml(configData, &appData)
		return appData, err
	default:
		return appData, os.ErrNotExist
	}
}

func loadConfigFromEnv(mapData *config) error {
	return envconfig.Process("tfcoach", mapData)
}

func loadConfigFromYaml(data []byte, mapData *config) error {
	return yaml.Unmarshal(data, mapData)
}

func loadConfigFromJSON(data []byte, mapData *config) error {
	return json.Unmarshal(data, mapData)
}

func getCustomConfigPath() (string, bool) {
	files := []string{
		".tfcoach.yml",
		".tfcoach.yaml",
		".tfcoach.json",
		".tfcoach",
	}

	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			return f, true
		}
	}

	return "", false
}
