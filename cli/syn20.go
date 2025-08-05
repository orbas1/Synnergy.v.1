package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn20 *tokens.SYN20Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn20",
		Short: "SYN20 token with pause and freeze",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN20 token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint32("decimals")
			id := tokenRegistry.NextID()
			syn20 = tokens.NewSYN20Token(id, name, symbol, uint8(dec))
			tokenRegistry.Register(syn20)
			fmt.Println("syn20 initialised")
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint32("decimals", 18, "decimal places")
	cmd.AddCommand(initCmd)

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause token operations",
		Run: func(cmd *cobra.Command, args []string) {
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			syn20.Pause()
			fmt.Println("paused")
		},
	}
	cmd.AddCommand(pauseCmd)

	unpauseCmd := &cobra.Command{
		Use:   "unpause",
		Short: "Resume operations",
		Run: func(cmd *cobra.Command, args []string) {
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			syn20.Unpause()
			fmt.Println("unpaused")
		},
	}
	cmd.AddCommand(unpauseCmd)

	freezeCmd := &cobra.Command{
		Use:   "freeze <addr>",
		Short: "Freeze an address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			syn20.Freeze(args[0])
			fmt.Println("frozen")
		},
	}
	cmd.AddCommand(freezeCmd)

	unfreezeCmd := &cobra.Command{
		Use:   "unfreeze <addr>",
		Short: "Unfreeze an address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			syn20.Unfreeze(args[0])
			fmt.Println("unfrozen")
		},
	}
	cmd.AddCommand(unfreezeCmd)

	mintCmd := &cobra.Command{
		Use:   "mint <to> <amt>",
		Short: "Mint tokens",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn20.Mint(args[0], amt); err != nil {
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
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn20.Burn(args[0], amt); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("burned")
		},
	}
	cmd.AddCommand(burnCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amt>",
		Short: "Transfer tokens",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := syn20.Transfer(args[0], args[1], amt); err != nil {
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
			if syn20 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn20.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balanceCmd)

	rootCmd.AddCommand(cmd)
}
