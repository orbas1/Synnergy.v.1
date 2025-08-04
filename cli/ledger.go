package cli

import (
	"fmt"

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
	ledgerCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(ledgerCmd)
}
