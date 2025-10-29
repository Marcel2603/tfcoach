package config

import "os"

type Navigator interface {
	GetHomeDir() (string, error)
	GetCustomConfigPath() (string, error)
}

type DefaultNavigator struct {
	CustomConfigPath string
}

func (*DefaultNavigator) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (n *DefaultNavigator) GetCustomConfigPath() (string, error) {
	if n.CustomConfigPath != "" {
		return n.CustomConfigPath, nil
	}
	return os.Getwd()
}
