package config_test

import (
	"errors"
	"testing"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/spf13/cobra"
)

func TestAddStandardFlags(t *testing.T) {
	cmd := &cobra.Command{}
	config.AddStandardFlags(cmd)

	expectedFlags := []string{"format", "no-color", "no-emojis", "config"}
	for _, flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("flag %s not found", flagName)
		}
	}
}

func TestParseStandardFlags_ShouldOverrideVariables(t *testing.T) {
	cmd := &cobra.Command{}
	config.AddStandardFlags(cmd)

	setupErr1 := cmd.Flags().Set("format", "json")
	setupErr2 := cmd.Flags().Set("no-color", "true")
	setupErr3 := cmd.Flags().Set("no-emojis", "false")
	if setupErr1 != nil || setupErr2 != nil || setupErr3 != nil {
		t.Errorf("setup error: %v", errors.Join(setupErr1, setupErr2, setupErr3))
	}

	err := config.ParseStandardFlags(cmd)
	if err != nil {
		t.Errorf("ParseStandardFlags() error = %v", err)
	}

	outputConfiguration := config.GetOutputConfiguration()
	if outputConfiguration.Format != "json" {
		t.Errorf("formatFlag = %v, want %v", outputConfiguration.Format, "json")
	}
	if outputConfiguration.Color.IsTrue {
		t.Errorf("ParseStandardFlags() did not set --no-color")
	}
	if !outputConfiguration.Emojis.IsTrue {
		t.Errorf("ParseStandardFlags() unexpectedly changed value of --no-emojis")
	}
}

func TestParseStandardFlags_ShouldFailOnInvalidFormat(t *testing.T) {
	cmd := &cobra.Command{}
	config.AddStandardFlags(cmd)

	setupErr := cmd.Flags().Set("format", "abcd")
	if setupErr != nil {
		t.Errorf("setup error: %v", setupErr)
	}

	err := config.ParseStandardFlags(cmd)
	if err == nil {
		t.Errorf("expected error, got none")
	}
}
