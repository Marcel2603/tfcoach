package config

import (
	_ "embed"
	"encoding/json"
	"errors"
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

	configuration config
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

func LoadDefaultConfig() error {
	var configData config
	err := loadConfigFromYaml(yamlDefaultData, &configData)
	if err != nil {
		return fmt.Errorf("could not load default config: %s", err)
	}

	err = configData.Validate()
	if err != nil {
		return fmt.Errorf("invalid default config: %s", err)
	}

	configuration = configData
	return nil
}

func LoadConfig(navigator Navigator) error {
	// 1. default config from ".tfcoach.default.yml"
	var configData config
	err := loadConfigFromYaml(yamlDefaultData, &configData)
	if err != nil {
		return err
	}

	// 2. config from home dir
	var homeConfigData config
	homeConfigData, err = loadConfigFromHomeDir(navigator)
	// TODO later: print error in debug log if err != nil
	if err == nil {
		mergeErr := mergeInto(&configData, homeConfigData)
		if mergeErr != nil {
			return mergeErr
		}
	}

	// 3. config from current dir
	var customConfigData config
	customConfigData, err = loadConfigFromLocalFile(navigator)
	// TODO later: print error in debug log if err != nil
	if err == nil {
		mergeErr := mergeInto(&configData, customConfigData)
		if mergeErr != nil {
			return mergeErr
		}
	}

	// 4. config from env
	var envData config
	err = loadConfigFromEnv(&envData)
	if err != nil {
		return err
	}
	mergeErr := mergeInto(&configData, envData)
	if mergeErr != nil {
		return mergeErr
	}

	// 5. validate
	validationErr := configData.Validate()
	if validationErr != nil {
		return validationErr
	}

	configuration = configData
	return nil
}

func loadConfigFromHomeDir(navigator Navigator) (config, error) {
	homeConfigPath, found := getHomeConfigPath(navigator)
	if !found {
		return config{}, errors.New("no config found in home directory")
	}

	homeConfigData, err := loadCustomConfigFromFile(homeConfigPath)
	if err != nil {
		return config{}, fmt.Errorf("could load config from home directory: %v", err)
	}

	return homeConfigData, nil
}

func loadConfigFromLocalFile(navigator Navigator) (config, error) {
	customConfigPath, found := getCustomConfigPath(navigator)
	if !found {
		return config{}, errors.New("no config found in local directory")
	}

	customConfigData, err := loadCustomConfigFromFile(customConfigPath)
	if err != nil {
		return config{}, fmt.Errorf("could load config from local directory: %v", err)
	}

	return customConfigData, nil
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

func getHomeConfigPath(navigator Navigator) (string, bool) {
	homeDir, err := navigator.GetHomeDir()
	if err != nil {
		return "", false
	}

	candidateBaseDirs := []string{
		filepath.Join(homeDir, ".config", "tfcoach"),
		filepath.Join(homeDir, ".tfcoach"),
	}

	for _, baseDir := range candidateBaseDirs {
		path, found := getFirstMatchingPath(baseDir, []string{
			".tfcoach.yml",
			".tfcoach.yaml",
			".tfcoach.json",
			".tfcoach",
		})
		if found {
			return path, true
		}
	}

	return "", false
}

func getCustomConfigPath(navigator Navigator) (string, bool) {
	path, err := navigator.GetCustomConfigPath()
	if err != nil {
		return "", false
	}

	var fi os.FileInfo
	fi, err = os.Stat(path)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Could not read custom config file at '%s': %s\n", path, err.Error())
		return "", false
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return getFirstMatchingPath(path, []string{
			".tfcoach.yml",
			".tfcoach.yaml",
			".tfcoach.json",
			".tfcoach",
		})
	case mode.IsRegular():
		return path, true
	}

	return "", false
}

func getFirstMatchingPath(baseDir string, fileNames []string) (string, bool) {
	for _, fileName := range fileNames {
		fullPath := filepath.Join(baseDir, fileName)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, true
		}
	}

	return "", false
}

func mergeInto(target *config, updated config) error {
	return mergo.Merge(target, updated, mergo.WithOverride, mergo.WithTransformers(NullableBoolTransformer{}))
}
