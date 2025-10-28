package config

import "os"

type Navigator interface {
	HomeDir() (string, error)
}

type DefaultNavigator struct{}

func (*DefaultNavigator) HomeDir() (string, error) {
	return os.UserHomeDir()
}
