package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/runner"
	"github.com/Marcel2603/tfcoach/rules/core"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	formatFlag     string
	noColorFlag    bool
	noEmojisFlag   bool
	configPathFlag string

	defaultOutputConfig config.OutputConfiguration
	finalOutputConfig   config.OutputConfiguration
)

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint Terraform files",
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		// TODO #16: use configPathFlag to influence custom config loading
		config.MustLoadConfig(&config.DefaultNavigator{})

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
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		src := engine.FileSystem{SkipDirs: []string{".git", ".terraform"}}
		code := runner.Lint(target, src, core.EnabledRules(), cmd.OutOrStdout(), finalOutputConfig.Format, finalOutputConfig.Emojis.IsTrue)
		os.Exit(code)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)

	config.MustLoadDefaultConfig()
	defaultOutputConfig = config.GetOutputConfiguration()

	formatUsageHelp := fmt.Sprintf("Output format. Supported: %s", strings.Join(config.SupportedFormats(), "|"))
	lintCmd.Flags().StringVarP(&formatFlag, "format", "f", defaultOutputConfig.Format, formatUsageHelp)

	lintCmd.Flags().BoolVar(&noColorFlag, "no-color", !defaultOutputConfig.Color.IsTrue, "Disable color output")
	lintCmd.Flags().BoolVar(&noEmojisFlag, "no-emojis", !defaultOutputConfig.Emojis.IsTrue, "Prevent emojis in output")

	lintCmd.Flags().StringVarP(&configPathFlag, "config", "c", "", "Custom config file path (default current directory)")

	lintCmd.Annotations = map[string]string{
		"exitCodes": "0:No issues found,1:Issues found,2:Runtime error",
	}
}
