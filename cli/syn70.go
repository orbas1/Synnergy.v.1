package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn70 *tokens.SYN70Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn70",
		Short: "SYN70 game asset token",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN70 token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint32("decimals")
			id := tokenRegistry.NextID()
			syn70 = tokens.NewSYN70Token(id, name, symbol, uint8(dec))
			tokenRegistry.Register(syn70)
			fmt.Println("syn70 initialised")
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint32("decimals", 0, "decimal places")
	cmd.AddCommand(initCmd)

	registerCmd := &cobra.Command{
		Use:   "register <id> <owner> <name> <game>",
		Short: "Register an in-game asset",
		Args:  cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70.RegisterAsset(args[0], args[1], args[2], args[3]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(registerCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <newOwner>",
		Short: "Transfer asset ownership",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70.TransferAsset(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(transferCmd)

	attrCmd := &cobra.Command{
		Use:   "setattr <id> <key> <value>",
		Short: "Set asset attribute",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70.SetAttribute(args[0], args[1], args[2]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(attrCmd)

	achCmd := &cobra.Command{
		Use:   "achievement <id> <name>",
		Short: "Record achievement",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70.AddAchievement(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(achCmd)

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Short: "Show asset info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			a, err := syn70.AssetInfo(args[0])
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("ID:%s Owner:%s Name:%s Game:%s\n", a.ID, a.Owner, a.Name, a.Game)
			if len(a.Attributes) > 0 {
				for k, v := range a.Attributes {
					fmt.Printf("%s=%s ", k, v)
				}
				fmt.Println()
			}
			if len(a.Achievements) > 0 {
				fmt.Println("Achievements:", strings.Join(a.Achievements, ","))
			}
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List assets",
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			assets := syn70.ListAssets()
			for _, a := range assets {
				fmt.Printf("%s %s %s %s\n", a.ID, a.Owner, a.Name, a.Game)
			}
		},
	}
	cmd.AddCommand(listCmd)

	balCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show token balance",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn70.BalanceOf(args[0]))
		},
	}
	cmd.AddCommand(balCmd)

	rootCmd.AddCommand(cmd)
}
