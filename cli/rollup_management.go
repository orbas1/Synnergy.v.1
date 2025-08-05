package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "rollupmgr",
		Short: "Control the rollup aggregator",
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause aggregator",
		Run: func(cmd *cobra.Command, args []string) {
			rollupMgr.Pause()
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume aggregator",
		Run: func(cmd *cobra.Command, args []string) {
			rollupMgr.Resume()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show pause status",
		Run: func(cmd *cobra.Command, args []string) {
			if rollupMgr.Status() {
				fmt.Println("paused")
			} else {
				fmt.Println("running")
			}
		},
	}

	cmd.AddCommand(pauseCmd, resumeCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
