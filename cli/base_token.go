package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var baseToken *tokens.BaseToken

func init() {
	cmd := &cobra.Command{
		Use:   "basetoken",
		Short: "Base token operations",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise a base token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint32("decimals")
			id := tokenRegistry.NextID()
			baseToken = tokens.NewBaseToken(tokens.TokenID(id), name, symbol, uint8(dec))
			tokenRegistry.Register(baseToken)
			fmt.Println("base token initialised")
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint32("decimals", 18, "decimal places")
	cmd.AddCommand(initCmd)

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show balance of an address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(baseToken.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balanceCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amt>",
		Short: "Transfer tokens",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := baseToken.Transfer(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("transferred")
		},
	}
	cmd.AddCommand(transferCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <to> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := baseToken.Mint(args[0], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("minted")
		},
	}
	cmd.AddCommand(mintCmd)

	burnCmd := &cobra.Command{
		Use:   "burn <from> <amt>",
		Short: "Burn tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := baseToken.Burn(args[0], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("burned")
		},
	}
	cmd.AddCommand(burnCmd)

	approveCmd := &cobra.Command{
		Use:   "approve <owner> <spender> <amt>",
		Short: "Approve allowance",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := baseToken.Approve(args[0], args[1], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("approved")
		},
	}
	cmd.AddCommand(approveCmd)

	allowanceCmd := &cobra.Command{
		Use:   "allowance <owner> <spender>",
		Short: "Check allowance",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(baseToken.Allowance(args[0], args[1]))
		},
	}
	cmd.AddCommand(allowanceCmd)

	rootCmd.AddCommand(cmd)
}
