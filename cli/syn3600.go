package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var contract *core.FuturesContract

func init() {
	cmd := &cobra.Command{
		Use:   "syn3600",
		Short: "Futures contract utilities",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a futures contract",
		Run: func(cmd *cobra.Command, args []string) {
			underlying, _ := cmd.Flags().GetString("underlying")
			qty, _ := cmd.Flags().GetUint64("qty")
			price, _ := cmd.Flags().GetUint64("price")
			expStr, _ := cmd.Flags().GetString("expiration")
			exp, _ := time.Parse(time.RFC3339, expStr)
			contract = core.NewFuturesContract(underlying, qty, price, exp)
			fmt.Println("contract created")
		},
	}
	createCmd.Flags().String("underlying", "", "underlying asset")
	createCmd.Flags().Uint64("qty", 0, "quantity")
	createCmd.Flags().Uint64("price", 0, "price per unit")
	createCmd.Flags().String("expiration", time.Now().Add(24*time.Hour).Format(time.RFC3339), "expiration time")
	cmd.AddCommand(createCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check if contract expired",
		Run: func(cmd *cobra.Command, args []string) {
			if contract == nil {
				fmt.Println("no contract")
				return
			}
			fmt.Println(contract.IsExpired(time.Now()))
		},
	}
	cmd.AddCommand(statusCmd)

	settleCmd := &cobra.Command{
		Use:   "settle <marketPrice>",
		Short: "Settle contract and show PnL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if contract == nil {
				fmt.Println("no contract")
				return
			}
			var price uint64
			fmt.Sscanf(args[0], "%d", &price)
			pnl := contract.Settle(price)
			fmt.Printf("pnl %d\n", pnl)
		},
	}
	cmd.AddCommand(settleCmd)

	rootCmd.AddCommand(cmd)
}
