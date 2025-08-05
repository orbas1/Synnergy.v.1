package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var transferManager = core.NewBridgeTransferManager()

func init() {
	cmd := &cobra.Command{
		Use:   "cross_chain_bridge",
		Short: "Manage cross-chain token transfers",
	}

	depositCmd := &cobra.Command{
		Use:   "deposit <bridge_id> <from> <to> <amount> [tokenID]",
		Args:  cobra.RangeArgs(4, 5),
		Short: "Lock assets for bridging",
		RunE: func(cmd *cobra.Command, args []string) error {
			amount, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			tokenID := ""
			if len(args) == 5 {
				tokenID = args[4]
			}
			t, err := transferManager.Deposit(args[0], args[1], args[2], amount, tokenID)
			if err != nil {
				return err
			}
			fmt.Println(t.ID)
			return nil
		},
	}

	claimCmd := &cobra.Command{
		Use:   "claim <transfer_id> <proof>",
		Args:  cobra.ExactArgs(2),
		Short: "Release assets using a proof",
		RunE: func(cmd *cobra.Command, args []string) error {
			return transferManager.Claim(args[0], args[1])
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show a transfer record",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, ok := transferManager.GetTransfer(args[0])
			if !ok {
				return fmt.Errorf("transfer not found")
			}
			fmt.Printf("%s: bridge=%s from=%s to=%s amount=%d token=%s status=%s\n", t.ID, t.BridgeID, t.From, t.To, t.Amount, t.TokenID, t.Status)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all transfers",
		Run: func(cmd *cobra.Command, args []string) {
			for _, t := range transferManager.ListTransfers() {
				fmt.Printf("%s: bridge=%s from=%s to=%s amount=%d token=%s status=%s\n", t.ID, t.BridgeID, t.From, t.To, t.Amount, t.TokenID, t.Status)
			}
		},
	}

	cmd.AddCommand(depositCmd, claimCmd, getCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
