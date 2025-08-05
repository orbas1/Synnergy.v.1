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
		Run: func(cmd *cobra.Command, args []string) {
			title, _ := cmd.Flags().GetString("title")
			artist, _ := cmd.Flags().GetString("artist")
			album, _ := cmd.Flags().GetString("album")
			musicToken = core.NewMusicToken(title, artist, album)
			fmt.Println("token initialised")
		},
	}
	initCmd.Flags().String("title", "", "song title")
	initCmd.Flags().String("artist", "", "artist")
	initCmd.Flags().String("album", "", "album")
	cmd.AddCommand(initCmd)

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show token info",
		Run: func(cmd *cobra.Command, args []string) {
			if musicToken == nil {
				fmt.Println("token not initialised")
				return
			}
			t, a, al := musicToken.Info()
			fmt.Printf("%s by %s on %s\n", t, a, al)
		},
	}
	cmd.AddCommand(infoCmd)

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Update token metadata",
		Run: func(cmd *cobra.Command, args []string) {
			if musicToken == nil {
				fmt.Println("token not initialised")
				return
			}
			title, _ := cmd.Flags().GetString("title")
			artist, _ := cmd.Flags().GetString("artist")
			album, _ := cmd.Flags().GetString("album")
			musicToken.Update(title, artist, album)
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
		Run: func(cmd *cobra.Command, args []string) {
			if musicToken == nil {
				fmt.Println("token not initialised")
				return
			}
			s, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid share")
				return
			}
			musicToken.SetRoyaltyShare(args[0], s)
		},
	}
	cmd.AddCommand(shareCmd)

	payoutCmd := &cobra.Command{
		Use:   "payout <amount>",
		Args:  cobra.ExactArgs(1),
		Short: "Distribute payout",
		Run: func(cmd *cobra.Command, args []string) {
			if musicToken == nil {
				fmt.Println("token not initialised")
				return
			}
			amt, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			payouts, err := musicToken.Distribute(amt)
			if err != nil {
				fmt.Println(err)
				return
			}
			for addr, p := range payouts {
				fmt.Printf("%s:%d\n", addr, p)
			}
		},
	}
	cmd.AddCommand(payoutCmd)

	rootCmd.AddCommand(cmd)
}
