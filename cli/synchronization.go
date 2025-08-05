package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/core"
)

var syncMgr = core.NewSyncManager(ledger)

func init() {
	syncCmd := &cobra.Command{Use: "synchronization", Short: "Blockchain synchronization"}

	startCmd := &cobra.Command{Use: "start", Short: "Start the sync manager", Run: func(cmd *cobra.Command, args []string) { syncMgr.Start() }}

	stopCmd := &cobra.Command{Use: "stop", Short: "Stop the sync manager", Run: func(cmd *cobra.Command, args []string) { syncMgr.Stop() }}

	statusCmd := &cobra.Command{Use: "status", Short: "Show sync status", Run: func(cmd *cobra.Command, args []string) {
		running, h := syncMgr.Status()
		fmt.Printf("running: %v height: %d\n", running, h)
	}}

	onceCmd := &cobra.Command{Use: "once", Short: "Run one sync round", RunE: func(cmd *cobra.Command, args []string) error {
		return syncMgr.Once()
	}}

	syncCmd.AddCommand(startCmd, stopCmd, statusCmd, onceCmd)
	rootCmd.AddCommand(syncCmd)
}
