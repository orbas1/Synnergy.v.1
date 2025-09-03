package cli

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn1000Index = tokens.NewSYN1000Index()

func init() {
	cmd := &cobra.Command{
		Use:   "syn1000index",
		Short: "Manage multiple SYN1000 tokens",
	}

	createCmd := &cobra.Command{
		Use:   "create <name> <symbol>",
		Short: "Create a SYN1000 token",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dec, _ := cmd.Flags().GetUint32("decimals")
			id := syn1000Index.Create(args[0], args[1], uint8(dec))
			fmt.Println(id)
		},
	}
	createCmd.Flags().Uint32("decimals", 18, "decimal places")
	cmd.AddCommand(createCmd)

	addResCmd := &cobra.Command{
		Use:   "add-reserve <id> <asset> <amount>",
		Short: "Add reserve asset to token",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			amt, ok := new(big.Rat).SetString(args[2])
			if !ok {
				fmt.Println("invalid amount")
				return
			}
			if err := syn1000Index.AddReserve(tokens.TokenID(id), args[1], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("reserve added")
		},
	}
	cmd.AddCommand(addResCmd)

	setPriceCmd := &cobra.Command{
		Use:   "set-price <id> <asset> <price>",
		Short: "Set reserve price",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			price, ok := new(big.Rat).SetString(args[2])
			if !ok {
				fmt.Println("invalid price")
				return
			}
			if err := syn1000Index.SetReservePrice(tokens.TokenID(id), args[1], price); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("price updated")
		},
	}
	cmd.AddCommand(setPriceCmd)

	valueCmd := &cobra.Command{
		Use:   "value <id>",
		Short: "Show total reserve value of token",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			v, err := syn1000Index.TotalValue(tokens.TokenID(id))
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(v.FloatString(2))
		},
	}
	cmd.AddCommand(valueCmd)

	rootCmd.AddCommand(cmd)
}
