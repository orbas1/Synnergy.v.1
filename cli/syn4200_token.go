package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var charityTok = tokens.NewSYN4200Token()

func init() {
	cmd := &cobra.Command{
		Use:   "syn4200",
		Short: "Charity token utilities",
	}

	donateCmd := &cobra.Command{
		Use:   "donate <symbol> <from> <amount>",
		Short: "Donate to campaign",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			purpose, _ := cmd.Flags().GetString("purpose")
			charityTok.Donate(args[0], args[1], amt, purpose)
			fmt.Println("donated")
		},
	}
	donateCmd.Flags().String("purpose", "", "campaign purpose")
	cmd.AddCommand(donateCmd)

	progressCmd := &cobra.Command{
		Use:   "progress <symbol>",
		Short: "Show campaign progress",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			amt, ok := charityTok.CampaignProgress(args[0])
			if !ok {
				fmt.Println("campaign not found")
				return
			}
			fmt.Println(amt)
		},
	}
	cmd.AddCommand(progressCmd)

	infoCmd := &cobra.Command{
		Use:   "info <symbol>",
		Short: "Show campaign info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c, ok := charityTok.Campaign(args[0])
			if !ok {
				fmt.Println("campaign not found")
				return
			}
			fmt.Printf("%+v\n", *c)
		},
	}
	cmd.AddCommand(infoCmd)

	rootCmd.AddCommand(cmd)
}
