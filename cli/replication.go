package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	repCmd := &cobra.Command{
		Use:   "replication",
		Short: "Control block replication",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the replicator",
		Run: func(cmd *cobra.Command, args []string) {
			replicator.Start()
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the replicator",
		Run: func(cmd *cobra.Command, args []string) {
			replicator.Stop()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show replicator status",
		Run: func(cmd *cobra.Command, args []string) {
			if replicator.Status() {
				fmt.Println("running")
			} else {
				fmt.Println("stopped")
			}
		},
	}

	replicateCmd := &cobra.Command{
		Use:   "replicate [hash]",
		Args:  cobra.ExactArgs(1),
		Short: "Mark a block as replicated",
		Run: func(cmd *cobra.Command, args []string) {
			if replicator.ReplicateBlock(args[0]) {
				fmt.Println("replicated")
			} else {
				fmt.Println("replication inactive")
			}
		},
	}

	repCmd.AddCommand(startCmd, stopCmd, statusCmd, replicateCmd)
	rootCmd.AddCommand(repCmd)
}
