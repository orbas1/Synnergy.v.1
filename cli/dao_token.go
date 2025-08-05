package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var daoTokenLedger = core.NewDAOTokenLedger()

func init() {
	tokenCmd := &cobra.Command{
		Use:   "dao-token",
		Short: "DAO token ledger operations",
	}

	mintCmd := &cobra.Command{
		Use:   "mint <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Mint tokens to an address",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			daoTokenLedger.Mint(args[0], amt)
		},
	}

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Transfer tokens",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := daoTokenLedger.Transfer(args[0], args[1], amt); err != nil {
				fmt.Println(err)
			}
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Get token balance",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(daoTokenLedger.Balance(args[0]))
		},
	}

	tokenCmd.AddCommand(mintCmd, transferCmd, balanceCmd)
	rootCmd.AddCommand(tokenCmd)
}
