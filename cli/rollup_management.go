package cli

import "github.com/spf13/cobra"

func init() {
	cmd := &cobra.Command{
		Use:   "rollupmgr",
		Short: "Control the rollup aggregator",
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause aggregator",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupPause")
			rollupMgr.Pause()
			printOutput(map[string]string{"status": "paused"})
			return nil
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume aggregator",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupResume")
			rollupMgr.Resume()
			printOutput(map[string]string{"status": "resumed"})
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show pause status",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupStatus")
			printOutput(map[string]bool{"paused": rollupMgr.Status()})
			return nil
		},
	}

	cmd.AddCommand(pauseCmd, resumeCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
