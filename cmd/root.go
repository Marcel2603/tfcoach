package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tfcoach",
	Short: "Tiny Terraform coach",
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

func GetRootCommand() *cobra.Command {
	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Annotations = map[string]string{
		"exitCodes": "0:OK",
	}
}
