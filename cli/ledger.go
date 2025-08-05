package cli

import (
	"encoding/json"
	"fmt"
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
			h, hash := ledger.Head()
			fmt.Printf("%d %s\n", h, hash)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "block [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch a block by height",
		Run: func(cmd *cobra.Command, args []string) {
			ht, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("invalid height:", err)
				return
			}
			if b, ok := ledger.GetBlock(ht); ok {
				out, _ := json.MarshalIndent(b, "", "  ")
				fmt.Println(string(out))
			} else {
				fmt.Println("not found")
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "balance [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Display token balance of an address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ledger.GetBalance(args[0]))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "utxo [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "List UTXOs for an address",
		Run: func(cmd *cobra.Command, args []string) {
			outs := ledger.GetUTXOs(args[0])
			out, _ := json.MarshalIndent(outs, "", "  ")
			fmt.Println(string(out))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "pool",
		Short: "List mem-pool transactions",
		Run: func(cmd *cobra.Command, args []string) {
			out, _ := json.MarshalIndent(ledger.Pool(), "", "  ")
			fmt.Println(string(out))
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "mint [addr] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Mint tokens to an address",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			ledger.Mint(args[0], amt)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "transfer [from] [to] [amount] [fee]",
		Args:  cobra.RangeArgs(3, 4),
		Short: "Transfer tokens between addresses",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			var fee uint64
			if len(args) == 4 {
				fee, err = strconv.ParseUint(args[3], 10, 64)
				if err != nil {
					fmt.Println("invalid fee")
					return
				}
			}
			if err := ledger.Transfer(args[0], args[1], amt, fee); err != nil {
				fmt.Println("error:", err)
			}
		},
	})

	rootCmd.AddCommand(cmd)
}
