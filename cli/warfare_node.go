package cli

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
	militarynodes "synnergy/internal/nodes/military_nodes"
)

var warfareNode *core.WarfareNode

func init() {
	cmd := &cobra.Command{
		Use:   "warfare",
		Short: "Interact with a warfare node",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create warfare node",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			if id == "" || addr == "" {
				return fmt.Errorf("id and addr required")
			}
			base := core.NewNode(id, addr, core.NewLedger())
			warfareNode = core.NewWarfareNode(base)
			printOutput(map[string]any{"status": "created", "id": id})
			return nil
		},
	}
	createCmd.Flags().String("id", "", "node id")
	createCmd.Flags().String("addr", "", "node address")
	cmd.AddCommand(createCmd)

	cmdCmd := &cobra.Command{
		Use:   "command <cmd>",
		Args:  cobra.ExactArgs(1),
		Short: "Execute secure command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			commander, _ := cmd.Flags().GetString("commander")
			privHex, _ := cmd.Flags().GetString("private")
			nonce, _ := cmd.Flags().GetUint64("nonce")
			meta, _ := cmd.Flags().GetStringArray("meta")
			metadata := parseMetadata(meta)
			command := strings.TrimSpace(args[0])
			if commander == "" && privHex == "" {
				if err := warfareNode.SecureCommand(command); err != nil {
					return err
				}
				printOutput(map[string]any{"status": "executed", "commander": "root"})
				return nil
			}
			if commander == "" || privHex == "" {
				return errors.New("commander and private key required for signed commands")
			}
			privBytes, err := hex.DecodeString(privHex)
			if err != nil {
				return fmt.Errorf("decode private key: %w", err)
			}
			req := core.CommandRequest{
				Commander: commander,
				Command:   command,
				Timestamp: time.Now().UTC(),
				Metadata:  metadata,
				Nonce:     nonce,
			}
			req.Signature = ed25519.Sign(ed25519.PrivateKey(privBytes), req.CanonicalPayload())
			record, err := warfareNode.ExecuteSecureCommand(context.Background(), req)
			if err != nil {
				return err
			}
			printOutput(map[string]any{
				"status":    "executed",
				"commander": record.Commander,
				"nonce":     record.Nonce,
				"accepted":  record.Accepted,
				"latency":   record.Latency.String(),
			})
			return nil
		},
	}
	cmdCmd.Flags().String("commander", "", "commander id for manual signing")
	cmdCmd.Flags().String("private", "", "hex encoded ed25519 private key")
	cmdCmd.Flags().Uint64("nonce", 0, "optional nonce override")
	cmdCmd.Flags().StringArray("meta", nil, "metadata key=value pairs")
	cmd.AddCommand(cmdCmd)

	trackCmd := &cobra.Command{
		Use:   "track <assetID> <location> <status>",
		Args:  cobra.ExactArgs(3),
		Short: "Record logistics information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			reporter, _ := cmd.Flags().GetString("reporter")
			meta, _ := cmd.Flags().GetStringArray("meta")
			update := core.LogisticsUpdate{
				AssetID:  args[0],
				Location: args[1],
				Status:   args[2],
				Reporter: reporter,
				Metadata: parseMetadata(meta),
			}
			rec, err := warfareNode.RecordLogistics(update)
			if err != nil {
				return err
			}
			printOutput(map[string]any{"status": "recorded", "asset": rec.AssetID, "location": rec.Location, "timestamp": rec.Timestamp})
			return nil
		},
	}
	trackCmd.Flags().String("reporter", "", "reporting unit or operator")
	trackCmd.Flags().StringArray("meta", nil, "metadata key=value pairs")
	cmd.AddCommand(trackCmd)

	listCmd := &cobra.Command{
		Use:   "logistics [assetID]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List logistics records",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			var recs []militarynodes.LogisticsRecord
			if len(args) == 1 {
				recs = warfareNode.LogisticsByAsset(args[0])
			} else {
				recs = warfareNode.Logistics()
			}
			printOutput(recs)
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	shareCmd := &cobra.Command{
		Use:   "share <info>",
		Args:  cobra.ExactArgs(1),
		Short: "Share tactical information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			meta, _ := cmd.Flags().GetStringArray("meta")
			if err := warfareNode.BroadcastTactical(args[0], parseMetadata(meta)); err != nil {
				return err
			}
			printOutput(map[string]any{"status": "broadcast"})
			return nil
		},
	}
	shareCmd.Flags().StringArray("meta", nil, "metadata key=value pairs")
	cmd.AddCommand(shareCmd)

	eventsCmd := &cobra.Command{
		Use:   "events",
		Args:  cobra.NoArgs,
		Short: "Stream recent events",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			since, _ := cmd.Flags().GetUint64("since")
			var events []core.WarfareEvent
			if since > 0 {
				events = warfareNode.EventsSince(since)
			} else {
				events = warfareNode.Events()
			}
			printOutput(events)
			return nil
		},
	}
	eventsCmd.Flags().Uint64("since", 0, "sequence filter")
	cmd.AddCommand(eventsCmd)

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Args:  cobra.NoArgs,
		Short: "Display warfare node metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			printOutput(warfareNode.MetricsSnapshot())
			return nil
		},
	}
	cmd.AddCommand(metricsCmd)

	commanderCmd := &cobra.Command{
		Use:   "commander",
		Short: "Manage commander keys",
	}

	issueCmd := &cobra.Command{
		Use:   "issue <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Issue a new commander credential",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			cred, err := warfareNode.IssueCommander(args[0])
			if err != nil {
				return err
			}
			printOutput(cred)
			return nil
		},
	}

	authCmd := &cobra.Command{
		Use:   "authorize <id> <pubkey>",
		Args:  cobra.ExactArgs(2),
		Short: "Authorize an external commander public key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			if err := warfareNode.AuthorizeCommander(args[0], args[1]); err != nil {
				return err
			}
			printOutput(map[string]any{"status": "authorized", "id": args[0]})
			return nil
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Revoke a commander",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return errors.New("node not initialised")
			}
			warfareNode.RevokeCommander(args[0])
			printOutput(map[string]any{"status": "revoked", "id": args[0]})
			return nil
		},
	}

	commanderCmd.AddCommand(issueCmd, authCmd, revokeCmd)
	cmd.AddCommand(commanderCmd)

	rootCmd.AddCommand(cmd)
}
