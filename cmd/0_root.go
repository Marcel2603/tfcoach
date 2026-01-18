package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Marcel2603/tfcoach/cmd/config"
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
	// this init function needs to run first in the module to ensure the default config is loaded, because the
	// subcommands depend on it, hence the file name "0_root"

	err := config.LoadDefaultConfig()
	if err != nil {
		panic(err)
	}

	rootCmd.Annotations = map[string]string{
		"exitCodes": "0:OK",
	}
}

func addStandardFlags(cmd *cobra.Command) {
	defaultOutputConfig = config.GetOutputConfiguration()

	formatUsageHelp := fmt.Sprintf("Output format. Supported: %s", strings.Join(config.SupportedFormats(), "|"))
	cmd.Flags().StringVarP(&formatFlag, "format", "f", defaultOutputConfig.Format, formatUsageHelp)

	cmd.Flags().BoolVar(&noColorFlag, "no-color", !defaultOutputConfig.Color.IsTrue, "Disable color output")

	cmd.Flags().BoolVar(&noEmojisFlag, "no-emojis", !defaultOutputConfig.Emojis.IsTrue, "Prevent emojis in output")

	cmd.Flags().StringVarP(&configPathFlag, "config", "c", "", "Custom config file path (default current directory)")
}

func parseStandardFlags(cmd *cobra.Command) error {
	err := config.LoadConfig(&config.DefaultNavigator{CustomConfigPath: configPathFlag})
	if err != nil {
		return err
	}

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
}
