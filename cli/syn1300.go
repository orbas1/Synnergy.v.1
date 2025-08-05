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
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			desc, _ := cmd.Flags().GetString("desc")
			owner, _ := cmd.Flags().GetString("owner")
			loc, _ := cmd.Flags().GetString("loc")
			if _, err := syn1300.Register(id, desc, owner, loc); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("asset registered")
		},
	}
	regCmd.Flags().String("id", "", "asset id")
	regCmd.Flags().String("desc", "", "description")
	regCmd.Flags().String("owner", "", "owner")
	regCmd.Flags().String("loc", "", "initial location")
	cmd.AddCommand(regCmd)

	updCmd := &cobra.Command{
		Use:   "update <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Update asset status",
		Run: func(cmd *cobra.Command, args []string) {
			loc, _ := cmd.Flags().GetString("loc")
			status, _ := cmd.Flags().GetString("status")
			note, _ := cmd.Flags().GetString("note")
			if err := syn1300.Update(args[0], loc, status, note); err != nil {
				fmt.Println(err)
			}
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
		Run: func(cmd *cobra.Command, args []string) {
			asset, ok := syn1300.Get(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%s owned by %s at %s status %s events %d\n", asset.ID, asset.Owner, asset.Location, asset.Status, len(asset.History))
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
