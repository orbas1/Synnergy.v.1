package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var currentNode = core.NewNode("node1", "localhost", ledger)

func init() {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Node operations",
	}
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show node information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ID: %s\nAddr: %s\n", currentNode.ID, currentNode.Addr)
		},
	}
	nodeCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(nodeCmd)
}
