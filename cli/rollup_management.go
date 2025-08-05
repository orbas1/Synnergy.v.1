package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rmCmd := &cobra.Command{
		Use:   "rollup_management",
		Short: "Control rollup aggregator",
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause the rollup aggregator",
		Run: func(cmd *cobra.Command, args []string) {
			rollupMgr.Pause()
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume the rollup aggregator",
		Run: func(cmd *cobra.Command, args []string) {
			rollupMgr.Resume()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show aggregator status",
		Run: func(cmd *cobra.Command, args []string) {
			if rollupMgr.Status() {
				fmt.Println("paused")
			} else {
				fmt.Println("running")
			}
		},
	}

	rmCmd.AddCommand(pauseCmd, resumeCmd, statusCmd)
	rootCmd.AddCommand(rmCmd)
}
