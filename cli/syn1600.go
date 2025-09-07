package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var musicToken *core.MusicToken

func init() {
	cmd := &cobra.Command{
		Use:   "syn1600",
		Short: "Music token utilities",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise a music token",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			title, _ := cmd.Flags().GetString("title")
			artist, _ := cmd.Flags().GetString("artist")
			album, _ := cmd.Flags().GetString("album")
			if title == "" || artist == "" || album == "" {
				return fmt.Errorf("title, artist and album must be provided")
			}
			musicToken = core.NewMusicToken(title, artist, album)
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("title", "", "song title")
	initCmd.Flags().String("artist", "", "artist")
	initCmd.Flags().String("album", "", "album")
	initCmd.MarkFlagRequired("title")
	initCmd.MarkFlagRequired("artist")
	initCmd.MarkFlagRequired("album")
	cmd.AddCommand(initCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show token info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if musicToken == nil {
				return fmt.Errorf("token not initialised")
			}
			t, a, al := musicToken.Info()
			cmd.Printf("%s by %s on %s\n", t, a, al)
			return nil
		},
	}
	cmd.AddCommand(infoCmd)

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update token metadata",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if musicToken == nil {
				return fmt.Errorf("token not initialised")
			}
			changed := cmd.Flags().Changed("title") || cmd.Flags().Changed("artist") || cmd.Flags().Changed("album")
			if !changed {
				return fmt.Errorf("at least one metadata field must be provided")
			}
			title, _ := cmd.Flags().GetString("title")
			artist, _ := cmd.Flags().GetString("artist")
			album, _ := cmd.Flags().GetString("album")
			musicToken.Update(title, artist, album)
			cmd.Println("metadata updated")
			return nil
		},
	}
	updateCmd.Flags().String("title", "", "song title")
	updateCmd.Flags().String("artist", "", "artist")
	updateCmd.Flags().String("album", "", "album")
	cmd.AddCommand(updateCmd)

	shareCmd := &cobra.Command{
		Use:   "share <addr> <share>",
		Args:  cobra.ExactArgs(2),
		Short: "Set royalty share",
		RunE: func(cmd *cobra.Command, args []string) error {
			if musicToken == nil {
				return fmt.Errorf("token not initialised")
			}
			s, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid share")
			}
			musicToken.SetRoyaltyShare(args[0], s)
			return nil
		},
	}
	cmd.AddCommand(shareCmd)

	payoutCmd := &cobra.Command{
		Use:   "payout <amount>",
		Args:  cobra.ExactArgs(1),
		Short: "Distribute payout",
		RunE: func(cmd *cobra.Command, args []string) error {
			if musicToken == nil {
				return fmt.Errorf("token not initialised")
			}
			amt, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			payouts, err := musicToken.Distribute(amt)
			if err != nil {
				return err
			}
			for addr, p := range payouts {
				fmt.Printf("%s:%d\n", addr, p)
			}
			return nil
		},
	}
	cmd.AddCommand(payoutCmd)

	rootCmd.AddCommand(cmd)
}
