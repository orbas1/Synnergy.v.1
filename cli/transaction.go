package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction utilities",
	}
	createCmd := &cobra.Command{
		Use:   "create [from] [to] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Create a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, 0)
			fmt.Printf("tx: %+v\n", tx)
		},
	}
	txCmd.AddCommand(createCmd)
	rootCmd.AddCommand(txCmd)
}
