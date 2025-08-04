package cli

import (
	"fmt"
	"strconv"

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
	broadcastCmd := &cobra.Command{
		Use:   "broadcast [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Broadcast transaction to all nodes",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			network.Broadcast(tx)
		},
	}
	networkCmd.AddCommand(addCmd, broadcastCmd)
	rootCmd.AddCommand(networkCmd)
}
