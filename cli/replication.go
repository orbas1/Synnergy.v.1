package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "replication",
		Short: "Control block replication",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Launch replication goroutines",
		Run:   func(cmd *cobra.Command, args []string) { replicator.Start() },
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the replication subsystem",
		Run:   func(cmd *cobra.Command, args []string) { replicator.Stop() },
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show replication status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(replicator.Status())
		},
	}

	replicateCmd := &cobra.Command{
		Use:   "replicate [hash]",
		Args:  cobra.ExactArgs(1),
		Short: "Gossip a known block",
		Run: func(cmd *cobra.Command, args []string) {
			if replicator.ReplicateBlock(args[0]) {
				fmt.Println("replicated")
			} else {
				fmt.Println("replication not running")
			}
		},
	}

	cmd.AddCommand(startCmd, stopCmd, statusCmd, replicateCmd)
	rootCmd.AddCommand(cmd)
}
