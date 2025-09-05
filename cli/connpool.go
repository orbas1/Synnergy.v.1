package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/core"
	ilog "synnergy/internal/log"
)

var pool = core.NewConnectionPool(8)

func init() {
	poolCmd := &cobra.Command{Use: "connpool", Short: "Manage connection pool"}

	statsCmd := &cobra.Command{Use: "stats", Short: "Show pool statistics", Run: func(cmd *cobra.Command, args []string) {
		s := pool.Stats()
		ilog.Info("cli_pool_stats", "active", s.Active, "capacity", s.Capacity)
		fmt.Printf("active: %d capacity: %d\n", s.Active, s.Capacity)
	}}

	dialCmd := &cobra.Command{Use: "dial [addr]", Args: cobra.ExactArgs(1), Short: "Dial an address using the pool", RunE: func(cmd *cobra.Command, args []string) error {
		_, err := pool.Dial(args[0])
		if err == nil {
			ilog.Info("cli_pool_dial", "id", args[0])
		} else {
			ilog.Error("cli_pool_dial", "error", err)
		}
		return err
	}}

	releaseCmd := &cobra.Command{Use: "release [addr]", Args: cobra.ExactArgs(1), Short: "Release a connection from the pool", Run: func(cmd *cobra.Command, args []string) {
		pool.Release(args[0])
		ilog.Info("cli_pool_release", "id", args[0])
	}}

	closeCmd := &cobra.Command{Use: "close", Short: "Close the pool", Run: func(cmd *cobra.Command, args []string) {
		pool.Close()
		ilog.Info("cli_pool_close")
	}}

	poolCmd.AddCommand(statsCmd, dialCmd, releaseCmd, closeCmd)
	rootCmd.AddCommand(poolCmd)
}
