package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn10 *tokens.SYN10Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn10",
		Short: "SYN10 CBDC token operations",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN10 token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			issuer, _ := cmd.Flags().GetString("issuer")
			rate, _ := cmd.Flags().GetFloat64("rate")
			dec, _ := cmd.Flags().GetUint32("decimals")
			id := tokenRegistry.NextID()
			syn10 = tokens.NewSYN10Token(id, name, symbol, issuer, rate, uint8(dec))
			tokenRegistry.Register(syn10)
			fmt.Println("syn10 initialised")
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("issuer", "", "issuer")
	initCmd.Flags().Float64("rate", 1.0, "fiat exchange rate")
	initCmd.Flags().Uint32("decimals", 2, "decimal places")
	cmd.AddCommand(initCmd)

	setRateCmd := &cobra.Command{
		Use:   "set-rate <rate>",
		Short: "Update exchange rate",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn10 == nil {
				fmt.Println("token not initialised")
				return
			}
			var rate float64
			fmt.Sscanf(args[0], "%f", &rate)
			syn10.SetExchangeRate(rate)
			fmt.Println("rate updated")
		},
	}
	cmd.AddCommand(setRateCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show token info",
		Run: func(cmd *cobra.Command, args []string) {
			if syn10 == nil {
				fmt.Println("token not initialised")
				return
			}
			info := syn10.Info()
			fmt.Printf("Name:%s Symbol:%s Issuer:%s Rate:%f Supply:%d\n", info.Name, info.Symbol, info.Issuer, info.ExchangeRate, info.TotalSupply)
		},
	}
	cmd.AddCommand(infoCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <to> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn10 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn10.Mint(args[0], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("minted")
		},
	}
	cmd.AddCommand(mintCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amt>",
		Short: "Transfer tokens",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if syn10 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := syn10.Transfer(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("transferred")
		},
	}
	cmd.AddCommand(transferCmd)

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show balance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn10 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn10.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balanceCmd)

	rootCmd.AddCommand(cmd)
}
