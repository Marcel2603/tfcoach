package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/Marcel2603/tfcoach/internal/runner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print [path]",
	Short: "Print tfcoach JSON report in another format",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		err := config.LoadConfig(&config.DefaultNavigator{CustomConfigPath: configPathFlag})
		if err != nil {
			return err
		}

		// TODO #52 avoid so much reuse between commands

		if cmd.Flags().Changed("format") {
			config.OverrideFormat(formatFlag)
		}
		if cmd.Flags().Changed("no-color") {
			config.OverrideColor(!noColorFlag)
		}
		if cmd.Flags().Changed("no-emojis") {
			config.OverrideEmojis(!noEmojisFlag)
		}

		finalOutputConfig = config.GetOutputConfiguration()
		color.NoColor = !finalOutputConfig.Color.IsTrue

		if slices.Contains(config.SupportedFormats(), finalOutputConfig.Format) {
			return nil
		}
		return fmt.Errorf("invalid --format: %s (want %s)", finalOutputConfig.Format, strings.Join(config.SupportedFormats(), "|"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srcReportPath := args[0]
		code := runner.Print(srcReportPath, cmd.OutOrStdout(), finalOutputConfig.Format, finalOutputConfig.Emojis.IsTrue)
		os.Exit(code)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}
	defaultOutputConfig = config.GetOutputConfiguration()

	formatUsageHelp := fmt.Sprintf("Output format. Supported: %s", strings.Join(config.SupportedFormats(), "|"))
	printCmd.Flags().StringVarP(&formatFlag, "format", "f", defaultOutputConfig.Format, formatUsageHelp)

	printCmd.Flags().BoolVar(&noColorFlag, "no-color", !defaultOutputConfig.Color.IsTrue, "Disable color output")
	printCmd.Flags().BoolVar(&noEmojisFlag, "no-emojis", !defaultOutputConfig.Emojis.IsTrue, "Prevent emojis in output")

	printCmd.Flags().StringVarP(&configPathFlag, "config", "c", "", "Custom config file path (default current directory)")

	printCmd.Annotations = map[string]string{
		"exitCodes": "0:No issues found,1:Read error,2:Conversion error",
	}
}
