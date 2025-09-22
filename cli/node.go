package cli

import (
	"context"
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
			gasPrint("NodeInfo")
			printOutput(map[string]any{
				"id":     currentNode.ID,
				"addr":   currentNode.Addr,
				"height": len(currentNode.Blockchain),
			})
			return nil
		},
	}
	stakeCmd := &cobra.Command{
		Use:   "stake [address] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Assign stake to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeStake")
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			if err := currentNode.SetStake(args[0], amt); err != nil {
				return err
			}
			printOutput(map[string]any{"address": args[0], "stake": amt})
			return nil
		},
	}
	slashCmd := &cobra.Command{
		Use:   "slash [address] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Slash a validator (reason: double|downtime)",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeSlash")
			switch args[1] {
			case "double":
				currentNode.ReportDoubleSign(args[0])
			case "downtime":
				currentNode.ReportDowntime(args[0])
			default:
				return fmt.Errorf("invalid reason")
			}
			printOutput("slashed")
			return nil
		},
	}
	rehabCmd := &cobra.Command{
		Use:   "rehab [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Rehabilitate a slashed validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeRehab")
			currentNode.Rehabilitate(args[0])
			printOutput("rehabilitated")
			return nil
		},
	}
	addTxCmd := &cobra.Command{
		Use:   "addtx [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Add a transaction to the mempool",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeAddTx")
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			fee, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee")
			}
			nonce, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid nonce")
			}
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			if err := currentNode.AddTransaction(tx); err != nil {
				return err
			}
			printOutput("transaction added")
			return nil
		},
	}
	mempoolCmd := &cobra.Command{
		Use:   "mempool",
		Short: "Show mempool size",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeMempool")
			printOutput(map[string]int{"size": len(currentNode.Mempool)})
			return nil
		},
	}
	mineCmd := &cobra.Command{
		Use:   "mine",
		Short: "Mine a block from the current mempool",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeMine")
			block, err := currentNode.MineBlock(context.Background())
			if err != nil {
				return err
			}
			if block == nil {
				printOutput("no transactions to mine")
				return nil
			}
			printOutput(map[string]any{"hash": block.Hash, "nonce": block.Nonce})
			return nil
		},
	}
	nodeCmd.AddCommand(infoCmd, stakeCmd, slashCmd, rehabCmd, addTxCmd, mempoolCmd, mineCmd)
	rootCmd.AddCommand(nodeCmd)
}
