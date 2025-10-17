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
	format  string
	noColor bool
)

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint Terraform files",
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(_ *cobra.Command, _ []string) error {
		color.NoColor = noColor

		if slices.Contains(config.GetSupportedOutputFormats(), format) {
			return nil
		}
		return fmt.Errorf("invalid --format: %s (want %s)", format, strings.Join(config.GetSupportedOutputFormats(), "|"))
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		src := engine.FileSystem{SkipDirs: []string{".git", ".terraform"}}
		code := runner.Lint(target, src, core.All(), cmd.OutOrStdout(), format)
		os.Exit(code)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)

	defaultOutputConfig, err := config.GetDefaultOutput()
	if err != nil {
		panic(err)
	}

	formatUsageHelp := fmt.Sprintf("Output format. Supported: %s", strings.Join(config.GetSupportedOutputFormats(), "|"))
	lintCmd.Flags().StringVarP(&format, "format", "f", defaultOutputConfig.Format, formatUsageHelp)

	lintCmd.Flags().BoolVar(&noColor, "no-color", !defaultOutputConfig.Color, "Disable color output")
}
