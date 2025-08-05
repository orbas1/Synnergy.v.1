package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "sidechainops",
		Short: "Side-chain escrow helpers",
	}

	depositCmd := &cobra.Command{
		Use:   "deposit [chainID] [from] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Deposit tokens to escrow",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := sideOps.Deposit(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	withdrawCmd := &cobra.Command{
		Use:   "withdraw [chainID] [from] [amount] [proof]",
		Args:  cobra.ExactArgs(4),
		Short: "Withdraw from escrow with proof",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := sideOps.Withdraw(args[0], args[1], amt, args[3]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance [chainID] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Show escrow balance for address",
		Run: func(cmd *cobra.Command, args []string) {
			bal, err := sideOps.EscrowBalance(args[0], args[1])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(bal)
		},
	}

	cmd.AddCommand(depositCmd, withdrawCmd, balanceCmd)
	rootCmd.AddCommand(cmd)
}
