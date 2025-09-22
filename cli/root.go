package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfgFile  string
	logLevel string
)

// rootCmd is the base command for the Synnergy CLI.
var rootCmd = &cobra.Command{
	Use:           "synnergy",
	Short:         "Synnergy blockchain CLI",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		switch logLevel {
		case "info", "debug":
		default:
			return fmt.Errorf("invalid log level: %s", logLevel)
		}
		if _, err := stage73State(); err != nil {
			return fmt.Errorf("load Stage 73 state: %w", err)
		}
		return nil
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		return saveStage73State()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log verbosity: info or debug")
	rootCmd.PersistentFlags().StringVar(&stage73StatePath, "stage73-state", "", "Path to Stage 73 state snapshot")
	rootCmd.AddCommand(guiCmd)
}

// Execute runs the root command.
func Execute() error { return rootCmd.Execute() }

// RootCmd exposes the root command for documentation generation and tooling.
func RootCmd() *cobra.Command { return rootCmd }
