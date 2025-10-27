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

	navig         navigator = &defaultNavigator{}
	configuration           = mustLoadConfig()
)

// TODO #36: use dependency injection here to make testing easier

type navigator interface {
	GetHomeDir() (string, error)
}

type defaultNavigator struct{}

func (*defaultNavigator) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

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
	// 1. default config from repo
	var configData config
	err := loadConfigFromYaml(yamlDefaultData, &configData)
	if err != nil {
		return config{}, err
	}

	// 2. config from home dir
	var homeDir string
	homeDir, err = navig.GetHomeDir()
	if err != nil {
		// TODO later: add debug log, home dir not defined
		_, _ = fmt.Fprintf(os.Stderr, "Could not get home directory: %s\n", err.Error())
	} else {
		homeConfigPath, found := getHomeConfigPath(homeDir)
		var homeConfigData config
		if found {
			homeConfigData, err = loadCustomConfigFromFile(homeConfigPath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Could not load config from home config file %s: %s\n", homeConfigPath, err.Error())
			} else {
				mergeErr := mergo.Merge(&configData, homeConfigData, mergo.WithOverride, mergo.WithTransformers(NullableBoolTransformer{}))
				if mergeErr != nil {
					return config{}, mergeErr
				}
			}
		}
	}

	// 3. config from current dir
	customConfigPath, found := getCustomConfigPath()
	var customConfigData config
	if found {
		customConfigData, err = loadCustomConfigFromFile(customConfigPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Could not load config from custom config file %s: %s\n", customConfigPath, err.Error())
		} else {
			mergeErr := mergo.Merge(&configData, customConfigData, mergo.WithOverride, mergo.WithTransformers(NullableBoolTransformer{}))
			if mergeErr != nil {
				return config{}, mergeErr
			}
		}
	}

	// 4. config from env
	var envData config
	err = loadConfigFromEnv(&envData)
	if err != nil {
		return config{}, err
	}
	mergeErr := mergo.Merge(&configData, envData, mergo.WithOverride, mergo.WithTransformers(NullableBoolTransformer{}))
	if mergeErr != nil {
		return config{}, mergeErr
	}

	validationErr := configData.Validate()
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

func getHomeConfigPath(homeDir string) (string, bool) {
	baseDir := filepath.Join(homeDir, ".config", "tfcoach")

	return getFirstMatchingPath(baseDir, []string{
		".tfcoach.yml",
		".tfcoach.yaml",
		".tfcoach.json",
		".tfcoach",
	})
}

func getCustomConfigPath() (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}

	return getFirstMatchingPath(cwd, []string{
		".tfcoach.yml",
		".tfcoach.yaml",
		".tfcoach.json",
		".tfcoach",
	})
}

func getFirstMatchingPath(baseDir string, paths []string) (string, bool) {
	for _, path := range paths {
		fullPath := filepath.Join(baseDir, path)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, true
		}
	}

	return "", false
}
