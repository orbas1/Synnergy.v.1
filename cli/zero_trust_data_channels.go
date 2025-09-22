package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

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
			owner, _ := cmd.Flags().GetString("owner")
			meta, _ := cmd.Flags().GetStringArray("meta")
			retention, _ := cmd.Flags().GetInt("retention")
			opts := []core.ChannelOption{core.WithOwner(owner), core.WithChannelMetadata(parseMetadata(meta))}
			if retention > 0 {
				opts = append(opts, core.WithRetention(retention))
			}
			if err := ztEngine.OpenChannel(args[0], key, opts...); err != nil {
				return err
			}
			info, err := ztEngine.ChannelInfo(args[0])
			if err != nil {
				return err
			}
			printOutput(info)
			return nil
		},
	}
	openCmd.Flags().String("owner", "", "channel owner")
	openCmd.Flags().StringArray("meta", nil, "metadata key=value pairs")
	openCmd.Flags().Int("retention", 0, "message retention count")

	sendCmd := &cobra.Command{
		Use:   "send [id] [msg]",
		Args:  cobra.ExactArgs(2),
		Short: "Send an encrypted message",
		RunE: func(cmd *cobra.Command, args []string) error {
			sender, _ := cmd.Flags().GetString("sender")
			cipher, err := ztEngine.SendAs(args[0], sender, []byte(args[1]))
			if err != nil {
				return err
			}
			printOutput(map[string]string{"cipher": fmt.Sprintf("%x", cipher)})
			return nil
		},
	}
	sendCmd.Flags().String("sender", "", "sender identity")

	listCmd := &cobra.Command{
		Use:   "messages [id]",
		Args:  cobra.ExactArgs(1),
		Short: "List encrypted messages",
		RunE: func(cmd *cobra.Command, args []string) error {
			msgs := ztEngine.Messages(args[0])
			rows := make([]map[string]string, len(msgs))
			for i, m := range msgs {
				rows[i] = map[string]string{
					"index":     strconv.Itoa(i),
					"cipher":    fmt.Sprintf("%x", m.Cipher),
					"sender":    m.Sender,
					"timestamp": m.Timestamp.Format(time.RFC3339Nano),
				}
			}
			printOutput(rows)
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
			printOutput(map[string]string{"plaintext": string(pt)})
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
			printOutput(map[string]string{"status": "closed"})
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Show channel information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := ztEngine.ChannelInfo(args[0])
			if err != nil {
				return err
			}
			printOutput(info)
			return nil
		},
	}

	authorizeCmd := &cobra.Command{
		Use:   "authorize [id] [participant] [pubkey]",
		Args:  cobra.ExactArgs(3),
		Short: "Authorize a participant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ztEngine.AuthorizePeer(args[0], args[1], args[2]); err != nil {
				return err
			}
			printOutput(map[string]string{"status": "authorized", "participant": args[1]})
			return nil
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke [id] [participant]",
		Args:  cobra.ExactArgs(2),
		Short: "Revoke a participant",
		RunE: func(cmd *cobra.Command, args []string) error {
			ztEngine.RevokePeer(args[0], args[1])
			printOutput(map[string]string{"status": "revoked", "participant": args[1]})
			return nil
		},
	}

	rotateCmd := &cobra.Command{
		Use:   "rotate [id] [hexkey]",
		Args:  cobra.ExactArgs(2),
		Short: "Rotate channel key",
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			if err := ztEngine.RotateKey(args[0], key); err != nil {
				return err
			}
			printOutput(map[string]string{"status": "rotated"})
			return nil
		},
	}

	eventsCmd := &cobra.Command{
		Use:   "events [id]",
		Args:  cobra.ExactArgs(1),
		Short: "List channel events",
		RunE: func(cmd *cobra.Command, args []string) error {
			since, _ := cmd.Flags().GetUint64("since")
			var events []core.ChannelEvent
			if since > 0 {
				events = ztEngine.EventsSince(since)
			} else {
				events = ztEngine.Events()
			}
			filtered := make([]core.ChannelEvent, 0, len(events))
			for _, ev := range events {
				if ev.ChannelID == args[0] {
					filtered = append(filtered, ev)
				}
			}
			printOutput(filtered)
			return nil
		},
	}
	eventsCmd.Flags().Uint64("since", 0, "sequence filter")

	cmd.AddCommand(openCmd, sendCmd, listCmd, recvCmd, authorizeCmd, revokeCmd, rotateCmd, closeCmd, infoCmd, eventsCmd)
	rootCmd.AddCommand(cmd)
}
