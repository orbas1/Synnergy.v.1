package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn3500 *core.SYN3500Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn3500",
		Short: "SYN3500 currency token",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			issuer, _ := cmd.Flags().GetString("issuer")
			rate, _ := cmd.Flags().GetFloat64("rate")
			if name == "" || symbol == "" || issuer == "" {
				return fmt.Errorf("name, symbol and issuer required")
			}
			if rate <= 0 {
				return fmt.Errorf("rate must be positive")
			}
			syn3500 = core.NewSYN3500Token(name, symbol, issuer, rate)
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "name")
	initCmd.Flags().String("symbol", "", "symbol")
	initCmd.Flags().String("issuer", "", "issuer")
	initCmd.Flags().Float64("rate", 1.0, "exchange rate")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("symbol")
	_ = initCmd.MarkFlagRequired("issuer")
	_ = initCmd.MarkFlagRequired("rate")
	cmd.AddCommand(initCmd)

	rateCmd := &cobra.Command{
		Use:   "setrate <rate>",
		Short: "Update exchange rate",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3500 == nil {
				return fmt.Errorf("token not initialised")
			}
			r, err := strconv.ParseFloat(args[0], 64)
			if err != nil || r <= 0 {
				return fmt.Errorf("invalid rate")
			}
			syn3500.SetRate(r)
			cmd.Println("rate updated")
			return nil
		},
	}
	cmd.AddCommand(rateCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3500 == nil {
				return fmt.Errorf("token not initialised")
			}
			sym, issuer, rate := syn3500.Info()
			cmd.Printf("%s %s %.2f\n", sym, issuer, rate)
			return nil
		},
	}
	cmd.AddCommand(infoCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <addr> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3500 == nil {
				return fmt.Errorf("token not initialised")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			syn3500.Mint(args[0], amt)
			cmd.Println("minted")
			return nil
		},
	}
	cmd.AddCommand(mintCmd)

	redeemCmd := &cobra.Command{
		Use:   "redeem <addr> <amt>",
		Short: "Redeem tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3500 == nil {
				return fmt.Errorf("token not initialised")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			if err := syn3500.Redeem(args[0], amt); err != nil {
				return err
			}
			cmd.Println("redeemed")
			return nil
		},
	}
	cmd.AddCommand(redeemCmd)

	balCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3500 == nil {
				return fmt.Errorf("token not initialised")
			}
			cmd.Println(syn3500.BalanceOf(args[0]))
			return nil
		},
	}
	cmd.AddCommand(balCmd)

	rootCmd.AddCommand(cmd)
}
