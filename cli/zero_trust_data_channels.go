package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			msgs := ztEngine.Messages(args[0])
			for i, m := range msgs {
				fmt.Printf("%d:%x\n", i, m.Cipher)
			}
			return nil
		},
	}

	recvCmd := &cobra.Command{
		Use:   "receive [id] [index]",
		Args:  cobra.ExactArgs(2),
		Short: "Decrypt and verify a message",
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			pt, err := ztEngine.Receive(args[0], idx)
			if err != nil {
				return err
			}
			fmt.Println(string(pt))
			return nil
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

	cmd.AddCommand(openCmd, sendCmd, listCmd, recvCmd, closeCmd)
	rootCmd.AddCommand(cmd)
}
