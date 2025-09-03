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
			if err != nil {
				return err
			}
			_, err = nftMarket.Mint(cmd.Context(), args[0], args[1], args[2], price, synn.GasCost("MintNFT"))
			return err
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
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s %d\n", nft.ID, nft.Owner, nft.Metadata, nft.Price)
			return nil
		},
	}

	buyCmd := &cobra.Command{
		Use:   "buy [id] [newOwner]",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer ownership of an NFT",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nftMarket.Buy(context.Background(), args[0], args[1], synn.GasCost("BuyNFT"))
		},
	}

	cmd.AddCommand(mintCmd, listCmd, buyCmd)
	rootCmd.AddCommand(cmd)
}
