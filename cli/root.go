package cli

import "github.com/spf13/cobra"

// rootCmd is the base command for the Synnergy CLI.
var rootCmd = &cobra.Command{
	Use:   "synnergy",
	Short: "Synnergy blockchain CLI",
}

// Execute runs the root command.
func Execute() error { return rootCmd.Execute() }
