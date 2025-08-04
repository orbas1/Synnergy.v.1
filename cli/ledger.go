package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var ledger = core.NewLedger()

func init() {
	ledgerCmd := &cobra.Command{
		Use:   "ledger",
		Short: "Interact with the ledger",
	}
	balanceCmd := &cobra.Command{
		Use:   "balance [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Get balance for address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ledger.GetBalance(args[0]))
		},
	}
	creditCmd := &cobra.Command{
		Use:   "credit [address] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Credit an address",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			ledger.Credit(args[0], amt)
		},
	}
	ledgerCmd.AddCommand(balanceCmd, creditCmd)
	rootCmd.AddCommand(ledgerCmd)
}
