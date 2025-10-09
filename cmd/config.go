package cmd

import (
	"fmt"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the current config",
	Run: func(cmd *cobra.Command, _ []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%+v \n", config.Configuration)
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Ruleconfig %+v \n", config.GetConfigByRuleID("core.naming_conv"))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
