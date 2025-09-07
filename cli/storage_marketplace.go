package cli

import (
	"context"
	"strconv"

	synn "synnergy"
	"synnergy/core"

	"github.com/spf13/cobra"
)

var storageMarket = core.NewStorageMarketplace(core.NewSimpleVM())

func init() {
	cmd := &cobra.Command{
		Use:   "storage_marketplace",
		Short: "List and lease decentralised storage",
	}

	var listGas uint64
	createCmd := &cobra.Command{
		Use:   "list [hash] [price] [owner]",
		Args:  cobra.ExactArgs(3),
		Short: "Create a storage listing",
		RunE: func(cmd *cobra.Command, args []string) error {
			price, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			id, err := storageMarket.CreateListing(cmd.Context(), args[0], price, args[2], listGas)
			if err != nil {
				return err
			}
			printOutput(map[string]string{"id": id})
			return nil
		},
	}
	createCmd.Flags().Uint64Var(&listGas, "gas", synn.GasCost("CreateListing"), "gas limit")

	listCmd := &cobra.Command{
		Use:   "listings",
		Short: "List storage offers as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := storageMarket.ListListings(context.Background())
			if err != nil {
				return err
			}
			printOutput(l)
			return nil
		},
	}

	var dealGas uint64
	dealCmd := &cobra.Command{
		Use:   "deal [listingID] [buyer]",
		Args:  cobra.ExactArgs(2),
		Short: "Open a storage deal",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := storageMarket.OpenDeal(cmd.Context(), args[0], args[1], dealGas)
			if err != nil {
				return err
			}
			printOutput(map[string]string{"id": id})
			return nil
		},
	}
	dealCmd.Flags().Uint64Var(&dealGas, "gas", synn.GasCost("OpenDeal"), "gas limit")

	cmd.AddCommand(createCmd, listCmd, dealCmd)
	rootCmd.AddCommand(cmd)
}
