package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn1000 *tokens.SYN1000Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn1000",
		Short: "SYN1000 stablecoin operations",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN1000 token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint32("decimals")
			id := tokenRegistry.NextID()
			syn1000 = tokens.NewSYN1000Token(id, name, symbol, uint8(dec))
			tokenRegistry.Register(syn1000)
			fmt.Println("syn1000 initialised")
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint32("decimals", 18, "decimal places")
	cmd.AddCommand(initCmd)

	addResCmd := &cobra.Command{
		Use:   "add-reserve <asset> <amount>",
		Short: "Add reserve asset",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn1000 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt float64
			fmt.Sscanf(args[1], "%f", &amt)
			syn1000.AddReserve(args[0], amt)
			fmt.Println("reserve added")
		},
	}
	cmd.AddCommand(addResCmd)

	setPriceCmd := &cobra.Command{
		Use:   "set-price <asset> <price>",
		Short: "Set reserve price",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn1000 == nil {
				fmt.Println("token not initialised")
				return
			}
			var price float64
			fmt.Sscanf(args[1], "%f", &price)
			syn1000.SetReservePrice(args[0], price)
			fmt.Println("price updated")
		},
	}
	cmd.AddCommand(setPriceCmd)

	valueCmd := &cobra.Command{
		Use:   "value",
		Short: "Show total reserve value",
		Run: func(cmd *cobra.Command, args []string) {
			if syn1000 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn1000.TotalReserveValue())
		},
	}
	cmd.AddCommand(valueCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <to> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn1000 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn1000.Mint(args[0], amt); err != nil {
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
			if syn1000 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := syn1000.Transfer(args[0], args[1], amt); err != nil {
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
			if syn1000 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn1000.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balanceCmd)

	rootCmd.AddCommand(cmd)
}
