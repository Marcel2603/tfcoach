/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"strings"

	"github.com/Marcel2603/tfcoach/cmd"
	"github.com/Marcel2603/tfcoach/internal/logging"
)

func main() {
	logLevel := strings.ToUpper(strings.TrimSpace(os.Getenv("TFCOACH_LOG")))
	logging.SetupLogger(logLevel)
	cmd.Execute()
}
