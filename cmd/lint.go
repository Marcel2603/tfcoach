package cmd

import (
	"os"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/Marcel2603/tfcoach/internal/engine"
	"github.com/Marcel2603/tfcoach/internal/runner"
	"github.com/Marcel2603/tfcoach/rules/core"
	"github.com/spf13/cobra"
)

var (
	includeTgCacheFlag bool
)

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint Terraform files",
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		err := parseStandardFlags(cmd)
		if err != nil {
			return err
		}

		if cmd.Flags().Changed("include-terragrunt-cache") {
			config.OverrideIncludeTgCache(includeTgCacheFlag)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		target := "."
		if len(args) > 0 {
			target = args[0]
		}

		skipDirs := []string{".git", ".terraform"}
		if !finalOutputConfig.IncludeTerragruntCache.IsTrue {
			skipDirs = append(skipDirs, ".terragrunt-cache")
		}

		src := engine.FileSystem{SkipDirs: skipDirs}
		code := runner.Lint(target, src, core.EnabledRules(), cmd.OutOrStdout(), finalOutputConfig.Format, finalOutputConfig.Emojis.IsTrue)
		os.Exit(code)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)

	addStandardFlags(lintCmd)

	lintCmd.Flags().BoolVar(
		&includeTgCacheFlag,
		"include-terragrunt-cache",
		defaultOutputConfig.IncludeTerragruntCache.IsTrue,
		"Include Terragrunt cache in scanned files",
	)

	lintCmd.Annotations = map[string]string{
		"exitCodes": "0:No issues found,1:Issues found,2:Runtime error",
	}
}
