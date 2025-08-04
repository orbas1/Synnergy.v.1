package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var network = core.NewNetwork()

func init() {
	networkCmd := &cobra.Command{
		Use:   "network",
		Short: "Network operations",
	}
	addCmd := &cobra.Command{
		Use:   "add [id] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Add node to network",
		Run: func(cmd *cobra.Command, args []string) {
			n := core.NewNode(args[0], args[1], core.NewLedger())
			network.AddNode(n)
			fmt.Println("node added")
		},
	}
	networkCmd.AddCommand(addCmd)
	rootCmd.AddCommand(networkCmd)
}
