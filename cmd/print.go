package cmd

import (
	"os"

	"github.com/Marcel2603/tfcoach/cmd/config"
	"github.com/Marcel2603/tfcoach/internal/runner"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print [path]",
	Short: "Print tfcoach JSON report in another format",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		return config.ParseStandardFlags(cmd)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srcReportPath := args[0]
		finalOutputConfig := config.GetOutputConfiguration()

		code := runner.Print(srcReportPath, cmd.OutOrStdout(), finalOutputConfig.Format, finalOutputConfig.Emojis.IsTrue)
		os.Exit(code)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)
	config.AddStandardFlags(printCmd)

	printCmd.Annotations = map[string]string{
		"exitCodes": "0:No issues found,1:Read error,2:Conversion error",
	}
}
