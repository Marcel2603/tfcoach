package cmd

import (
	"os"

	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/runner"
	"github.com/Marcel2603/tfcoach/rules/core"
	"github.com/spf13/cobra"
)

var (
	lintPath string
)

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint Terraform files",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		src := engine.FileSystem{SkipDirs: []string{".git", ".terraform"}}
		code := runner.Lint(target, src, core.All(), cmd.OutOrStdout())
		os.Exit(code)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)
	lintCmd.Flags().StringVarP(&lintPath, "path", "p", ".", "Path to scan (default: current dir)")
}
