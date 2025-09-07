package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var legalRegistry = core.NewLegalTokenRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn4700",
		Short: "Manage SYN4700 legal tokens",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a legal token",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			doctype, _ := cmd.Flags().GetString("doctype")
			hash, _ := cmd.Flags().GetString("hash")
			owner, _ := cmd.Flags().GetString("owner")
			expiry, _ := cmd.Flags().GetInt64("expiry")
			supply, _ := cmd.Flags().GetUint64("supply")
			parties, _ := cmd.Flags().GetStringArray("party")
			if id == "" || name == "" || symbol == "" || doctype == "" || hash == "" || owner == "" {
				return fmt.Errorf("all fields required")
			}
			if expiry <= 0 {
				return fmt.Errorf("invalid expiry")
			}
			if supply == 0 {
				return fmt.Errorf("supply must be positive")
			}
			if len(parties) == 0 {
				return fmt.Errorf("at least one party required")
			}
			t := core.NewLegalToken(id, name, symbol, doctype, hash, owner, time.Unix(expiry, 0), supply, parties)
			legalRegistry.Add(t)
			cmd.Println("token created")
			return nil
		},
	}
	createCmd.Flags().String("id", "", "token id")
	createCmd.Flags().String("name", "", "name")
	createCmd.Flags().String("symbol", "", "symbol")
	createCmd.Flags().String("doctype", "", "document type")
	createCmd.Flags().String("hash", "", "document hash")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Int64("expiry", 0, "expiry unix timestamp")
	createCmd.Flags().Uint64("supply", 0, "token supply")
	createCmd.Flags().StringArray("party", []string{}, "parties")
	createCmd.MarkFlagRequired("id")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("doctype")
	createCmd.MarkFlagRequired("hash")
	createCmd.MarkFlagRequired("owner")
	createCmd.MarkFlagRequired("expiry")
	createCmd.MarkFlagRequired("supply")
	createCmd.MarkFlagRequired("party")

	signCmd := &cobra.Command{
		Use:   "sign <id> <party> <sig>",
		Args:  cobra.ExactArgs(3),
		Short: "Add a party signature",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, ok := legalRegistry.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			if err := t.Sign(args[1], args[2]); err != nil {
				return err
			}
			cmd.Println("signed")
			return nil
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke <id> <party>",
		Args:  cobra.ExactArgs(2),
		Short: "Revoke a signature",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, ok := legalRegistry.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			t.RevokeSignature(args[1])
			cmd.Println("revoked")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <id> <status>",
		Args:  cobra.ExactArgs(2),
		Short: "Update token status",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, ok := legalRegistry.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			s := core.LegalTokenStatus(args[1])
			switch s {
			case core.LegalTokenStatusPending, core.LegalTokenStatusActive, core.LegalTokenStatusCompleted, core.LegalTokenStatusDisputed:
				t.UpdateStatus(s)
				return nil
			default:
				return fmt.Errorf("invalid status")
			}
		},
	}

	disputeCmd := &cobra.Command{
		Use:   "dispute <id> <action> [result]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Record a dispute",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, ok := legalRegistry.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			if args[1] == "" {
				return fmt.Errorf("action required")
			}
			res := ""
			if len(args) == 3 {
				res = args[2]
			}
			t.Dispute(args[1], res)
			cmd.Println("disputed")
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, ok := legalRegistry.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.MarshalIndent(t, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	cmd.AddCommand(createCmd, signCmd, revokeCmd, statusCmd, disputeCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
