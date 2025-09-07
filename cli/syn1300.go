package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn1300 = core.NewSupplyChainRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn1300",
		Short: "Supply chain asset registry",
	}

	regCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new asset",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			desc, _ := cmd.Flags().GetString("desc")
			owner, _ := cmd.Flags().GetString("owner")
			loc, _ := cmd.Flags().GetString("loc")
			if id == "" || desc == "" || owner == "" || loc == "" {
				return fmt.Errorf("id, desc, owner and loc must be provided")
			}
			if _, err := syn1300.Register(id, desc, owner, loc); err != nil {
				return err
			}
			cmd.Println("asset registered")
			return nil
		},
	}
	regCmd.Flags().String("id", "", "asset id")
	regCmd.Flags().String("desc", "", "description")
	regCmd.Flags().String("owner", "", "owner")
	regCmd.Flags().String("loc", "", "initial location")
	regCmd.MarkFlagRequired("id")
	regCmd.MarkFlagRequired("desc")
	regCmd.MarkFlagRequired("owner")
	regCmd.MarkFlagRequired("loc")
	cmd.AddCommand(regCmd)

	updCmd := &cobra.Command{
		Use:   "update <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Update asset status",
		RunE: func(cmd *cobra.Command, args []string) error {
			loc, _ := cmd.Flags().GetString("loc")
			status, _ := cmd.Flags().GetString("status")
			note, _ := cmd.Flags().GetString("note")
			if loc == "" && status == "" && note == "" {
				return fmt.Errorf("at least one of loc, status or note must be provided")
			}
			if err := syn1300.Update(args[0], loc, status, note); err != nil {
				return err
			}
			cmd.Println("updated")
			return nil
		},
	}
	updCmd.Flags().String("loc", "", "location")
	updCmd.Flags().String("status", "", "status")
	updCmd.Flags().String("note", "", "note")
	cmd.AddCommand(updCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get asset info",
		RunE: func(cmd *cobra.Command, args []string) error {
			asset, ok := syn1300.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			cmd.Printf("%s owned by %s at %s status %s events %d\n", asset.ID, asset.Owner, asset.Location, asset.Status, len(asset.History))
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
