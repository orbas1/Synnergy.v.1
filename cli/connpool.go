package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/core"
)

var pool = core.NewConnectionPool(8)

func init() {
	poolCmd := &cobra.Command{Use: "connpool", Short: "Manage connection pool"}

	statsCmd := &cobra.Command{Use: "stats", Short: "Show pool statistics", Run: func(cmd *cobra.Command, args []string) {
		s := pool.Stats()
		fmt.Printf("active: %d capacity: %d\n", s.Active, s.Capacity)
	}}

	dialCmd := &cobra.Command{Use: "dial [addr]", Args: cobra.ExactArgs(1), Short: "Dial an address using the pool", RunE: func(cmd *cobra.Command, args []string) error {
		_, err := pool.Dial(args[0])
		return err
	}}

	closeCmd := &cobra.Command{Use: "close", Short: "Close the pool", Run: func(cmd *cobra.Command, args []string) {
		pool.Close()
	}}

	poolCmd.AddCommand(statsCmd, dialCmd, closeCmd)
	rootCmd.AddCommand(poolCmd)
}
