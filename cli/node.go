package cli

import (
	"fmt"
	"strconv"

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
			fmt.Printf("ID: %s\nAddr: %s\nBlockchain Height: %d\n", currentNode.ID, currentNode.Addr, len(currentNode.Blockchain))
		},
	}
	stakeCmd := &cobra.Command{
		Use:   "stake [address] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Assign stake to an address",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			if err := currentNode.SetStake(args[0], amt); err != nil {
				fmt.Println("error:", err)
			}
		},
	}
	slashCmd := &cobra.Command{
		Use:   "slash [address] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Slash a validator (reason: double|downtime)",
		Run: func(cmd *cobra.Command, args []string) {
			switch args[1] {
			case "double":
				currentNode.ReportDoubleSign(args[0])
			case "downtime":
				currentNode.ReportDowntime(args[0])
			}
		},
	}
	rehabCmd := &cobra.Command{
		Use:   "rehab [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Rehabilitate a slashed validator",
		Run: func(cmd *cobra.Command, args []string) {
			currentNode.Rehabilitate(args[0])
		},
	}
	addTxCmd := &cobra.Command{
		Use:   "addtx [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Add a transaction to the mempool",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			if err := currentNode.AddTransaction(tx); err != nil {
				fmt.Println("error:", err)
			}
		},
	}
	mempoolCmd := &cobra.Command{
		Use:   "mempool",
		Short: "Show mempool size",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(len(currentNode.Mempool))
		},
	}
	mineCmd := &cobra.Command{
		Use:   "mine",
		Short: "Mine a block from the current mempool",
		Run: func(cmd *cobra.Command, args []string) {
			block := currentNode.MineBlock()
			if block == nil {
				fmt.Println("no transactions to mine")
				return
			}
			fmt.Printf("mined block %s with nonce %d\n", block.Hash, block.Nonce)
		},
	}
	nodeCmd.AddCommand(infoCmd, stakeCmd, slashCmd, rehabCmd, addTxCmd, mempoolCmd, mineCmd)
	rootCmd.AddCommand(nodeCmd)
}
