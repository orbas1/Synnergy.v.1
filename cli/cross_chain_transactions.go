package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var crossTxManager = core.NewCrossChainTxManager(core.NewLedger())

func init() {
	cmd := &cobra.Command{
		Use:   "cross_tx",
		Short: "Execute cross-chain asset transfers",
	}

	var listJSON bool
	var getJSON bool

	lockMintCmd := &cobra.Command{
		Use:   "lockmint <bridge_id> <asset_id> <amount> <proof> --from <addr> --to <addr>",
		Args:  cobra.ExactArgs(4),
		Short: "Lock native assets and mint wrapped tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			id, err := crossTxManager.LockMint(parseInt(args[0]), from, to, args[1], amount, args[3])
			if err != nil {
				return err
			}
			fmt.Printf("%d gas:%d\n", id, synnergy.GasCost("LockMint"))
			return nil
		},
	}
	lockMintCmd.Flags().String("from", "", "Sender address")
	lockMintCmd.Flags().String("to", "", "Recipient address")
	lockMintCmd.MarkFlagRequired("from")
	lockMintCmd.MarkFlagRequired("to")

	burnReleaseCmd := &cobra.Command{
		Use:   "burnrelease <bridge_id> <to> <asset_id> <amount> --from <addr>",
		Args:  cobra.ExactArgs(4),
		Short: "Burn wrapped tokens and release native assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			amount, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			id, err := crossTxManager.BurnRelease(parseInt(args[0]), from, args[1], args[2], amount)
			if err != nil {
				return err
			}
			fmt.Printf("%d gas:%d\n", id, synnergy.GasCost("BurnRelease"))
			return nil
		},
	}
	burnReleaseCmd.Flags().String("from", "", "Sender address")
	burnReleaseCmd.MarkFlagRequired("from")

	getCmd := &cobra.Command{
		Use:   "get <tx_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a cross-chain transfer by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			t, err := crossTxManager.GetTransfer(id)
			if err != nil {
				return err
			}
			if getJSON {
				enc, _ := json.Marshal(t)
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("%d: bridge=%d from=%s to=%s asset=%s amount=%d type=%s completed=%v\n", t.ID, t.BridgeID, t.From, t.To, t.AssetID, t.Amount, t.Type, t.Completed)
			return nil
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List cross-chain transfer records",
		Run: func(cmd *cobra.Command, args []string) {
			ts := crossTxManager.ListTransfers()
			if listJSON {
				enc, _ := json.Marshal(ts)
				fmt.Println(string(enc))
				return
			}
			for _, t := range ts {
				fmt.Printf("%d: bridge=%d from=%s to=%s asset=%s amount=%d type=%s completed=%v\n", t.ID, t.BridgeID, t.From, t.To, t.AssetID, t.Amount, t.Type, t.Completed)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	cmd.AddCommand(lockMintCmd, burnReleaseCmd, listCmd, getCmd)
	rootCmd.AddCommand(cmd)
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
