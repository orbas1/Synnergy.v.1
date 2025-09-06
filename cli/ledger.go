package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var ledger = core.NewLedger()

func init() {
	cmd := &cobra.Command{
		Use:   "ledger",
		Short: "Interact with the ledger",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "head",
		Short: "Show chain height and latest block hash",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerHead")
			h, hash := ledger.Head()
			printOutput(map[string]any{"height": h, "hash": hash})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "block [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch a block by height",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerBlock")
			ht, err := strconv.Atoi(args[0])
			if err != nil {
				printOutput(map[string]any{"error": "invalid height"})
				return
			}
			if b, ok := ledger.GetBlock(ht); ok {
				printOutput(b)
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "balance [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Display token balance of an address",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerBalance")
			printOutput(ledger.GetBalance(args[0]))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "utxo [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "List UTXOs for an address",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerUTXO")
			outs := ledger.GetUTXOs(args[0])
			printOutput(outs)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "pool",
		Short: "List mem-pool transactions",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerPool")
			printOutput(ledger.Pool())
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "mint [addr] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Mint tokens to an address",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerMint")
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid amount"})
				return
			}
			ledger.Mint(args[0], amt)
			printOutput(map[string]any{"status": "minted", "address": args[0], "amount": amt})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "transfer [from] [to] [amount] [fee]",
		Args:  cobra.RangeArgs(3, 4),
		Short: "Transfer tokens between addresses",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LedgerTransfer")
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid amount"})
				return
			}
			var fee uint64
			if len(args) == 4 {
				fee, err = strconv.ParseUint(args[3], 10, 64)
				if err != nil {
					printOutput(map[string]any{"error": "invalid fee"})
					return
				}
			}
			if err := ledger.Transfer(args[0], args[1], amt, fee); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "transferred", "from": args[0], "to": args[1], "amount": amt, "fee": fee})
		},
	})

	rootCmd.AddCommand(cmd)
}
