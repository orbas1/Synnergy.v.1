package cli

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ReplicationStart")
			replicator.Start()
			printOutput(map[string]string{"status": "started"})
			return nil
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the replication subsystem",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ReplicationStop")
			replicator.Stop()
			printOutput(map[string]string{"status": "stopped"})
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show replication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ReplicationStatus")
			printOutput(map[string]bool{"running": replicator.Status()})
			return nil
		},
	}

	replicateCmd := &cobra.Command{
		Use:   "replicate [hash]",
		Args:  cobra.ExactArgs(1),
		Short: "Gossip a known block",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ReplicateBlock")
			if replicator.ReplicateBlock(args[0]) {
				printOutput(map[string]string{"status": "replicated", "hash": args[0]})
			} else {
				printOutput(map[string]string{"error": "replication not running"})
			}
			return nil
		},
	}

	cmd.AddCommand(startCmd, stopCmd, statusCmd, replicateCmd)
	rootCmd.AddCommand(cmd)
}
