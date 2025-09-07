package cli

import (
	"fmt"
	"strconv"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			underlying, _ := cmd.Flags().GetString("underlying")
			qty, _ := cmd.Flags().GetUint64("qty")
			price, _ := cmd.Flags().GetUint64("price")
			expStr, _ := cmd.Flags().GetString("expiration")
			if underlying == "" {
				return fmt.Errorf("underlying required")
			}
			if qty == 0 || price == 0 {
				return fmt.Errorf("quantity and price must be positive")
			}
			exp, err := time.Parse(time.RFC3339, expStr)
			if err != nil {
				return fmt.Errorf("invalid expiration: %w", err)
			}
			contract = core.NewFuturesContract(underlying, qty, price, exp)
			cmd.Println("contract created")
			return nil
		},
	}
	createCmd.Flags().String("underlying", "", "underlying asset")
	createCmd.Flags().Uint64("qty", 0, "quantity")
	createCmd.Flags().Uint64("price", 0, "price per unit")
	createCmd.Flags().String("expiration", "", "expiration time (RFC3339)")
	_ = createCmd.MarkFlagRequired("underlying")
	_ = createCmd.MarkFlagRequired("qty")
	_ = createCmd.MarkFlagRequired("price")
	_ = createCmd.MarkFlagRequired("expiration")
	cmd.AddCommand(createCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Check if contract expired",
		RunE: func(cmd *cobra.Command, args []string) error {
			if contract == nil {
				return fmt.Errorf("no contract")
			}
			cmd.Println(contract.IsExpired(time.Now()))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	settleCmd := &cobra.Command{
		Use:   "settle <marketPrice>",
		Short: "Settle contract and show PnL",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if contract == nil {
				return fmt.Errorf("no contract")
			}
			price, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid price")
			}
			pnl := contract.Settle(price)
			cmd.Printf("pnl %d\n", pnl)
			return nil
		},
	}
	cmd.AddCommand(settleCmd)

	rootCmd.AddCommand(cmd)
}
