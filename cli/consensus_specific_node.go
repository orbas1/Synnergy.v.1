package cli

import (
	"context"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := parseMode(args[0])
			if err != nil {
				return err
			}
			csNode = core.NewConsensusSpecificNode(mode, args[1], args[2], ledger)
			return nil
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
			b, err := csNode.MineBlock(context.Background())
			if err != nil {
				fmt.Printf("mining error: %v\n", err)
				return
			}
			if b != nil {
				fmt.Printf("mined block %s with nonce %d\n", b.Hash, b.Nonce)
			}
		},
	}

	stakeCmd := &cobra.Command{
		Use:   "stake [address] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Assign stake on node ledger",
		RunE: func(cmd *cobra.Command, args []string) error {
			if csNode == nil {
				return fmt.Errorf("no node")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			return csNode.SetStake(args[0], amt)
		},
	}

	cmd.AddCommand(createCmd, infoCmd, mineCmd, stakeCmd)
	rootCmd.AddCommand(cmd)
}
