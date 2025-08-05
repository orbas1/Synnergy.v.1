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
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			issuer, _ := cmd.Flags().GetString("issuer")
			rate, _ := cmd.Flags().GetFloat64("rate")
			syn3500 = core.NewSYN3500Token(name, symbol, issuer, rate)
			fmt.Println("token initialised")
		},
	}
	initCmd.Flags().String("name", "", "name")
	initCmd.Flags().String("symbol", "", "symbol")
	initCmd.Flags().String("issuer", "", "issuer")
	initCmd.Flags().Float64("rate", 1.0, "exchange rate")
	cmd.AddCommand(initCmd)

	rateCmd := &cobra.Command{
		Use:   "setrate <rate>",
		Short: "Update exchange rate",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn3500 == nil {
				fmt.Println("token not initialised")
				return
			}
			r, _ := strconv.ParseFloat(args[0], 64)
			syn3500.SetRate(r)
			fmt.Println("rate updated")
		},
	}
	cmd.AddCommand(rateCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show token info",
		Run: func(cmd *cobra.Command, args []string) {
			if syn3500 == nil {
				fmt.Println("token not initialised")
				return
			}
			sym, issuer, rate := syn3500.Info()
			fmt.Printf("%s %s %.2f\n", sym, issuer, rate)
		},
	}
	cmd.AddCommand(infoCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <addr> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn3500 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			syn3500.Mint(args[0], amt)
			fmt.Println("minted")
		},
	}
	cmd.AddCommand(mintCmd)

	redeemCmd := &cobra.Command{
		Use:   "redeem <addr> <amt>",
		Short: "Redeem tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn3500 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn3500.Redeem(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("redeemed")
			}
		},
	}
	cmd.AddCommand(redeemCmd)

	balCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show balance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn3500 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn3500.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balCmd)

	rootCmd.AddCommand(cmd)
}
