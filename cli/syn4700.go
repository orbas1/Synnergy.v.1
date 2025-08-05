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
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			doctype, _ := cmd.Flags().GetString("doctype")
			hash, _ := cmd.Flags().GetString("hash")
			owner, _ := cmd.Flags().GetString("owner")
			expiry, _ := cmd.Flags().GetInt64("expiry")
			supply, _ := cmd.Flags().GetUint64("supply")
			parties, _ := cmd.Flags().GetStringArray("party")
			t := core.NewLegalToken(id, name, symbol, doctype, hash, owner, time.Unix(expiry, 0), supply, parties)
			legalRegistry.Add(t)
			fmt.Println("token created")
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

	signCmd := &cobra.Command{
		Use:   "sign <id> <party> <sig>",
		Args:  cobra.ExactArgs(3),
		Short: "Add a party signature",
		Run: func(cmd *cobra.Command, args []string) {
			if t, ok := legalRegistry.Get(args[0]); ok {
				if err := t.Sign(args[1], args[2]); err != nil {
					fmt.Println("error:", err)
				}
			}
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke <id> <party>",
		Args:  cobra.ExactArgs(2),
		Short: "Revoke a signature",
		Run: func(cmd *cobra.Command, args []string) {
			if t, ok := legalRegistry.Get(args[0]); ok {
				t.RevokeSignature(args[1])
			}
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <id> <status>",
		Args:  cobra.ExactArgs(2),
		Short: "Update token status",
		Run: func(cmd *cobra.Command, args []string) {
			if t, ok := legalRegistry.Get(args[0]); ok {
				t.UpdateStatus(core.LegalTokenStatus(args[1]))
			}
		},
	}

	disputeCmd := &cobra.Command{
		Use:   "dispute <id> <action> [result]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Record a dispute",
		Run: func(cmd *cobra.Command, args []string) {
			if t, ok := legalRegistry.Get(args[0]); ok {
				res := ""
				if len(args) == 3 {
					res = args[2]
				}
				t.Dispute(args[1], res)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show token info",
		Run: func(cmd *cobra.Command, args []string) {
			if t, ok := legalRegistry.Get(args[0]); ok {
				b, _ := json.MarshalIndent(t, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	}

	cmd.AddCommand(createCmd, signCmd, revokeCmd, statusCmd, disputeCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
