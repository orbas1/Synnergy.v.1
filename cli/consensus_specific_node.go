package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var csNode *core.ConsensusSpecificNode

func init() {
	cmd := &cobra.Command{
		Use:   "consensus-node",
		Short: "Consensus specific node operations",
	}

	createCmd := &cobra.Command{
		Use:   "create [mode] [id] [addr]",
		Args:  cobra.ExactArgs(3),
		Short: "Create a node locked to a consensus mode",
		Run: func(cmd *cobra.Command, args []string) {
			mode := parseMode(args[0])
			csNode = core.NewConsensusSpecificNode(mode, args[1], args[2], ledger)
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show node information",
		Run: func(cmd *cobra.Command, args []string) {
			if csNode == nil {
				fmt.Println("no node")
				return
			}
			fmt.Printf("ID: %s\nAddr: %s\nMode: %s\nHeight: %d\n", csNode.ID, csNode.Addr, csNode.Mode, len(csNode.Blockchain))
		},
	}

	mineCmd := &cobra.Command{
		Use:   "mine",
		Short: "Mine a block with current node",
		Run: func(cmd *cobra.Command, args []string) {
			if csNode == nil {
				fmt.Println("no node")
				return
			}
			b := csNode.MineBlock()
			if b != nil {
				fmt.Printf("mined block %s with nonce %d\n", b.Hash, b.Nonce)
			}
		},
	}

	stakeCmd := &cobra.Command{
		Use:   "stake [address] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Assign stake on node ledger",
		Run: func(cmd *cobra.Command, args []string) {
			if csNode == nil {
				fmt.Println("no node")
				return
			}
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			if err := csNode.SetStake(args[0], amt); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	cmd.AddCommand(createCmd, infoCmd, mineCmd, stakeCmd)
	rootCmd.AddCommand(cmd)
}
