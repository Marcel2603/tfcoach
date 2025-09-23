package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	BuildVersion = "dev"
	BuildCommit  = "none"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "tfcoach %s (%s)\n", BuildVersion, BuildCommit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
