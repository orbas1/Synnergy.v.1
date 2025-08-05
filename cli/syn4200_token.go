package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn4200 = core.NewSYN4200Token()

func init() {
	cmd := &cobra.Command{
		Use:   "syn4200_token",
		Short: "Charity token operations",
	}

	donateCmd := &cobra.Command{
		Use:   "donate <symbol>",
		Args:  cobra.ExactArgs(1),
		Short: "Donate to a charity campaign",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString("from")
			amt, _ := cmd.Flags().GetUint64("amt")
			purpose, _ := cmd.Flags().GetString("purpose")
			syn4200.Donate(args[0], from, amt, purpose)
			fmt.Println("donation recorded")
		},
	}
	donateCmd.Flags().String("from", "", "donor address")
	donateCmd.Flags().Uint64("amt", 0, "donation amount")
	donateCmd.Flags().String("purpose", "", "campaign purpose")
	donateCmd.MarkFlagRequired("from")
	donateCmd.MarkFlagRequired("amt")

	progressCmd := &cobra.Command{
		Use:   "progress <symbol>",
		Args:  cobra.ExactArgs(1),
		Short: "Show campaign progress",
		Run: func(cmd *cobra.Command, args []string) {
			if amt, ok := syn4200.CampaignProgress(args[0]); ok {
				fmt.Println(amt)
			} else {
				fmt.Println("not found")
			}
		},
	}

	cmd.AddCommand(donateCmd, progressCmd)
	rootCmd.AddCommand(cmd)
}
