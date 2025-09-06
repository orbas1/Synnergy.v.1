package cli

import (
	"context"
	"fmt"
	"strconv"

	synn "synnergy"
	"synnergy/core"

	"github.com/spf13/cobra"
)

var nftMarket = core.NewNFTMarketplace()

func init() {
	cmd := &cobra.Command{
		Use:   "nft",
		Short: "Mint and trade NFTs",
	}

	mintCmd := &cobra.Command{
		Use:   "mint [id] [owner] [metadata] [price]",
		Args:  cobra.ExactArgs(4),
		Short: "Mint a new NFT",
		RunE: func(cmd *cobra.Command, args []string) error {
			price, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil || price == 0 {
				return fmt.Errorf("invalid price: %s", args[3])
			}
			gasPrint("MintNFT")
			_, err = nftMarket.Mint(cmd.Context(), args[0], args[1], args[2], price, synn.GasCost("MintNFT"))
			if err != nil {
				return err
			}
			printOutput(map[string]any{
				"status":   "minted",
				"id":       args[0],
				"owner":    args[1],
				"metadata": args[2],
				"price":    price,
			})
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Show NFT details",
		RunE: func(cmd *cobra.Command, args []string) error {
			nft, err := nftMarket.List(args[0])
			if err != nil {
				return err
			}
			gasPrint("ListNFT")
			printOutput(nft)
			return nil
		},
	}

	buyCmd := &cobra.Command{
		Use:   "buy [id] [newOwner]",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer ownership of an NFT",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("BuyNFT")
			if err := nftMarket.Buy(context.Background(), args[0], args[1], synn.GasCost("BuyNFT")); err != nil {
				return err
			}
			printOutput(map[string]any{"status": "transferred", "id": args[0], "newOwner": args[1]})
			return nil
		},
	}

	cmd.AddCommand(mintCmd, listCmd, buyCmd)
	rootCmd.AddCommand(cmd)
}
