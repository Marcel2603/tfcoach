package config

import (
	"fmt"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	formatFlag     string
	noColorFlag    bool
	noEmojisFlag   bool
	configPathFlag string

	defaultOutputConfig OutputConfiguration
	finalOutputConfig   OutputConfiguration
)

func AddStandardFlags(cmd *cobra.Command) {
	defaultOutputConfig = GetOutputConfiguration()

	formatUsageHelp := fmt.Sprintf("Output format. Supported: %s", strings.Join(SupportedFormats(), "|"))
	cmd.Flags().StringVarP(&formatFlag, "format", "f", defaultOutputConfig.Format, formatUsageHelp)

	cmd.Flags().BoolVar(&noColorFlag, "no-color", !defaultOutputConfig.Color.IsTrue, "Disable color output")

	cmd.Flags().BoolVar(&noEmojisFlag, "no-emojis", !defaultOutputConfig.Emojis.IsTrue, "Prevent emojis in output")

	cmd.Flags().StringVarP(&configPathFlag, "config", "c", "", "Custom config file path (default current directory)")
}

func ParseStandardFlags(cmd *cobra.Command) error {
	err := LoadConfig(&DefaultNavigator{CustomConfigPath: configPathFlag})
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("format") {
		OverrideFormat(formatFlag)
	}
	if cmd.Flags().Changed("no-color") {
		OverrideColor(!noColorFlag)
	}
	if cmd.Flags().Changed("no-emojis") {
		OverrideEmojis(!noEmojisFlag)
	}

	finalOutputConfig = GetOutputConfiguration()
	color.NoColor = !finalOutputConfig.Color.IsTrue

	if slices.Contains(SupportedFormats(), finalOutputConfig.Format) {
		return nil
	}
	return fmt.Errorf("invalid --format: %s (want %s)", finalOutputConfig.Format, strings.Join(SupportedFormats(), "|"))
}
