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
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "ID: %s\nAddr: %s\nBlockchain Height: %d\n", currentNode.ID, currentNode.Addr, len(currentNode.Blockchain))
			return nil
		},
	}
	stakeCmd := &cobra.Command{
		Use:   "stake [address] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Assign stake to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			return currentNode.SetStake(args[0], amt)
		},
	}
	slashCmd := &cobra.Command{
		Use:   "slash [address] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Slash a validator (reason: double|downtime)",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[1] {
			case "double":
				currentNode.ReportDoubleSign(args[0])
			case "downtime":
				currentNode.ReportDowntime(args[0])
			}
			return nil
		},
	}
	rehabCmd := &cobra.Command{
		Use:   "rehab [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Rehabilitate a slashed validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentNode.Rehabilitate(args[0])
			return nil
		},
	}
	addTxCmd := &cobra.Command{
		Use:   "addtx [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Add a transaction to the mempool",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			return currentNode.AddTransaction(tx)
		},
	}
	mempoolCmd := &cobra.Command{
		Use:   "mempool",
		Short: "Show mempool size",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), len(currentNode.Mempool))
			return nil
		},
	}
	mineCmd := &cobra.Command{
		Use:   "mine",
		Short: "Mine a block from the current mempool",
		RunE: func(cmd *cobra.Command, args []string) error {
			block := currentNode.MineBlock()
			if block == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "no transactions to mine")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "mined block %s with nonce %d\n", block.Hash, block.Nonce)
			return nil
		},
	}
	nodeCmd.AddCommand(infoCmd, stakeCmd, slashCmd, rehabCmd, addTxCmd, mempoolCmd, mineCmd)
	rootCmd.AddCommand(nodeCmd)
}
