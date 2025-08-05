package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn70Tok *tokens.SYN70Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn70",
		Short: "SYN70 in-game assets",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN70 token",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetUint64("id")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint64("decimals")
			syn70Tok = tokens.NewSYN70Token(tokens.TokenID(id), name, symbol, uint8(dec))
			fmt.Println("token initialised")
		},
	}
	initCmd.Flags().Uint64("id", 0, "token id")
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint64("decimals", 0, "decimal places")
	cmd.AddCommand(initCmd)

	regCmd := &cobra.Command{
		Use:   "register-asset",
		Short: "Register asset",
		Run: func(cmd *cobra.Command, args []string) {
			if syn70Tok == nil {
				fmt.Println("token not initialised")
				return
			}
			id, _ := cmd.Flags().GetString("id")
			owner, _ := cmd.Flags().GetString("owner")
			name, _ := cmd.Flags().GetString("name")
			game, _ := cmd.Flags().GetString("game")
			if err := syn70Tok.RegisterAsset(id, owner, name, game); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("asset registered")
			}
		},
	}
	regCmd.Flags().String("id", "", "asset id")
	regCmd.Flags().String("owner", "", "asset owner")
	regCmd.Flags().String("name", "", "asset name")
	regCmd.Flags().String("game", "", "game name")
	cmd.AddCommand(regCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer-asset <assetID> <newOwner>",
		Short: "Transfer asset",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70Tok == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70Tok.TransferAsset(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("transferred")
			}
		},
	}
	cmd.AddCommand(transferCmd)

	attrCmd := &cobra.Command{
		Use:   "set-attr <assetID> <key> <value>",
		Short: "Set asset attribute",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70Tok == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70Tok.SetAttribute(args[0], args[1], args[2]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("attribute set")
			}
		},
	}
	cmd.AddCommand(attrCmd)

	achCmd := &cobra.Command{
		Use:   "add-achievement <assetID> <name>",
		Short: "Add asset achievement",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70Tok == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn70Tok.AddAchievement(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("achievement added")
			}
		},
	}
	cmd.AddCommand(achCmd)

	infoCmd := &cobra.Command{
		Use:   "info <assetID>",
		Short: "Show asset info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn70Tok == nil {
				fmt.Println("token not initialised")
				return
			}
			asset, err := syn70Tok.AssetInfo(args[0])
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("%+v\n", asset)
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List assets",
		Run: func(cmd *cobra.Command, args []string) {
			if syn70Tok == nil {
				fmt.Println("token not initialised")
				return
			}
			for _, a := range syn70Tok.ListAssets() {
				fmt.Printf("%+v\n", a)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
