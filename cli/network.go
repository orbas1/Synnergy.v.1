package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var network = core.NewNetwork(biometricSvc)

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
	relayCmd := &cobra.Command{
		Use:   "relay [id] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Add relay node to network",
		Run: func(cmd *cobra.Command, args []string) {
			n := core.NewNode(args[0], args[1], core.NewLedger())
			network.AddRelay(n)
			fmt.Println("relay node added")
		},
	}
	broadcastCmd := &cobra.Command{
		Use:   "broadcast [userID] [biometric] [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(7),
		Short: "Broadcast transaction with biometric verification",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[4], 10, 64)
			fee, _ := strconv.ParseUint(args[5], 10, 64)
			nonce, _ := strconv.ParseUint(args[6], 10, 64)
			tx := core.NewTransaction(args[2], args[3], amt, fee, nonce)
			if err := network.Broadcast(tx, args[0], []byte(args[1])); err != nil {
				fmt.Println("broadcast failed:", err)
				return
			}
			fmt.Println("transaction queued for broadcast")
		},
	}
	networkCmd.AddCommand(addCmd, relayCmd, broadcastCmd)
	rootCmd.AddCommand(networkCmd)
}
