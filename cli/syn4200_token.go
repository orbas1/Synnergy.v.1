package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "syn4200_token",
		Short: "Charity token operations",
	}

	donateCmd := &cobra.Command{
		Use:   "donate <symbol>",
		Args:  cobra.ExactArgs(1),
		Short: "Donate to a charity campaign",
		RunE: func(cmd *cobra.Command, args []string) error {
			from, _ := cmd.Flags().GetString("from")
			amt, _ := cmd.Flags().GetUint64("amt")
			purpose, _ := cmd.Flags().GetString("purpose")
			if from == "" || amt == 0 {
				return fmt.Errorf("from and positive amt required")
			}
			store, err := stage73State()
			if err != nil {
				return err
			}
			registry := store.Charity()
			registry.Donate(args[0], from, amt, purpose)
			markStage73Dirty()
			cmd.Println("donation recorded")
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := stage73State()
			if err != nil {
				return err
			}
			amt, ok := store.Charity().CampaignProgress(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			fmt.Fprintln(cmd.OutOrStdout(), amt)
			return nil
		},
	}

	cmd.AddCommand(donateCmd, progressCmd)
	rootCmd.AddCommand(cmd)
}
