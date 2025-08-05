package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var ztEngine = core.NewZeroTrustEngine()

func init() {
	cmd := &cobra.Command{
		Use:   "zero-trust",
		Short: "Manage zero trust data channels",
	}

	openCmd := &cobra.Command{
		Use:   "open [id] [hexkey]",
		Args:  cobra.ExactArgs(2),
		Short: "Open a secure channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			if err := ztEngine.OpenChannel(args[0], key); err != nil {
				return err
			}
			fmt.Println("channel opened")
			return nil
		},
	}

	sendCmd := &cobra.Command{
		Use:   "send [id] [msg]",
		Args:  cobra.ExactArgs(2),
		Short: "Send an encrypted message",
		RunE: func(cmd *cobra.Command, args []string) error {
			cipher, err := ztEngine.Send(args[0], []byte(args[1]))
			if err != nil {
				return err
			}
			fmt.Printf("%x\n", cipher)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "messages [id]",
		Args:  cobra.ExactArgs(1),
		Short: "List encrypted messages",
		Run: func(cmd *cobra.Command, args []string) {
			for _, m := range ztEngine.Messages(args[0]) {
				fmt.Printf("%x\n", m)
			}
		},
	}

	closeCmd := &cobra.Command{
		Use:   "close [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Close a channel",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ztEngine.CloseChannel(args[0]); err != nil {
				return err
			}
			fmt.Println("channel closed")
			return nil
		},
	}

	cmd.AddCommand(openCmd, sendCmd, listCmd, closeCmd)
	rootCmd.AddCommand(cmd)
}
