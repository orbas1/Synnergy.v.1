package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	replicator = core.NewReplicator(ledger)
	initSvc    = core.NewInitService(replicator)
)

func init() {
	initrepCmd := &cobra.Command{
		Use:   "initrep",
		Short: "Initialization replication control",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Bootstrap ledger and start replication",
		Run: func(cmd *cobra.Command, args []string) {
			initSvc.Start()
			gasPrint("InitrepStart")
			printOutput(map[string]any{"status": "started"})
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop initialization service",
		Run: func(cmd *cobra.Command, args []string) {
			initSvc.Stop()
			gasPrint("InitrepStop")
			printOutput(map[string]any{"status": "stopped"})
		},
	}

	initrepCmd.AddCommand(startCmd, stopCmd)
	rootCmd.AddCommand(initrepCmd)
}
