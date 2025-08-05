package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	opsCmd := &cobra.Command{
		Use:   "sidechain_ops",
		Short: "Operate on side-chain escrows",
	}

	depositCmd := &cobra.Command{
		Use:   "deposit <chainID> <addr> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens to a side-chain escrow",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			return sidechainOps.Deposit(args[0], args[1], amt)
		},
	}

	withdrawCmd := &cobra.Command{
		Use:   "withdraw <chainID> <addr> <amount> <proof>",
		Args:  cobra.ExactArgs(4),
		Short: "Withdraw from a side-chain escrow",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			return sidechainOps.Withdraw(args[0], args[1], amt, args[3])
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance <chainID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Show escrow balance for address",
		RunE: func(cmd *cobra.Command, args []string) error {
			bal, err := sidechainOps.EscrowBalance(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(bal)
			return nil
		},
	}

	opsCmd.AddCommand(depositCmd, withdrawCmd, balanceCmd)
	rootCmd.AddCommand(opsCmd)
}
