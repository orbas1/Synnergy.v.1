package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var legalReg = tokens.NewLegalTokenRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn4700",
		Short: "Legal token registry",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create legal token",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			doctype, _ := cmd.Flags().GetString("doctype")
			hash, _ := cmd.Flags().GetString("hash")
			owner, _ := cmd.Flags().GetString("owner")
			expiryStr, _ := cmd.Flags().GetString("expiry")
			supply, _ := cmd.Flags().GetUint64("supply")
			partiesStr, _ := cmd.Flags().GetString("parties")
			expiry, _ := time.Parse(time.RFC3339, expiryStr)
			parties := []string{}
			if partiesStr != "" {
				parties = strings.Split(partiesStr, ",")
			}
			t := tokens.NewLegalToken(id, name, symbol, doctype, hash, owner, expiry, supply, parties)
			legalReg.Add(t)
			fmt.Println("created")
		},
	}
	createCmd.Flags().String("id", "", "token id")
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("doctype", "", "document type")
	createCmd.Flags().String("hash", "", "document hash")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().String("expiry", time.Now().Add(24*time.Hour).Format(time.RFC3339), "expiry time")
	createCmd.Flags().Uint64("supply", 0, "token supply")
	createCmd.Flags().String("parties", "", "comma separated parties")
	cmd.AddCommand(createCmd)

	signCmd := &cobra.Command{
		Use:   "sign <id> <party> <sig>",
		Short: "Add party signature",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			t, ok := legalReg.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			if err := t.Sign(args[1], args[2]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("signed")
			}
		},
	}
	cmd.AddCommand(signCmd)

	revokeCmd := &cobra.Command{
		Use:   "revoke <id> <party>",
		Short: "Revoke party signature",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			t, ok := legalReg.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			t.RevokeSignature(args[1])
			fmt.Println("revoked")
		},
	}
	cmd.AddCommand(revokeCmd)

	statusCmd := &cobra.Command{
		Use:   "status <id> <status>",
		Short: "Update token status",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			t, ok := legalReg.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			t.UpdateStatus(tokens.LegalTokenStatus(args[1]))
			fmt.Println("updated")
		},
	}
	cmd.AddCommand(statusCmd)

	disputeCmd := &cobra.Command{
		Use:   "dispute <id> <action> [result]",
		Short: "Record dispute action",
		Args:  cobra.RangeArgs(2, 3),
		Run: func(cmd *cobra.Command, args []string) {
			t, ok := legalReg.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			res := ""
			if len(args) == 3 {
				res = args[2]
			}
			t.Dispute(args[1], res)
			fmt.Println("disputed")
		},
	}
	cmd.AddCommand(disputeCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get token info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			t, ok := legalReg.Get(args[0])
			if !ok {
				fmt.Println("token not found")
				return
			}
			fmt.Printf("%+v\n", *t)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List legal tokens",
		Run: func(cmd *cobra.Command, args []string) {
			for _, t := range legalReg.List() {
				fmt.Printf("%+v\n", *t)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
