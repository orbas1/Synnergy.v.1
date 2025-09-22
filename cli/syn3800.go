package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "syn3800",
		Short: "Manage SYN3800 grant records",
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800Create")
			store, err := stage73State()
			if err != nil {
				return err
			}
			if args[0] == "" || args[1] == "" {
				return fmt.Errorf("beneficiary and name required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			registry := store.Grants()
			id := registry.CreateGrant(args[0], args[1], amt)
			authorizer, _ := cmd.Flags().GetString("authorizer")
			if authorizer != "" {
				path, password, err := parseCredentialPair(authorizer)
				if err != nil {
					return err
				}
				wallet, err := loadWallet(path, password)
				if err != nil {
					return err
				}
				if err := registry.AddAuthorizer(id, wallet.Address); err != nil {
					return err
				}
			}
			markStage73Dirty()
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}
	createCmd.Flags().String("authorizer", "", "optional controller credential path:password")
	cmd.AddCommand(createCmd)

	releaseCmd := &cobra.Command{
		Use:   "release <id> <amount> [note]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Release funds for a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800Release")
			store, err := stage73State()
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			registry := store.Grants()
			if err := registry.DisburseWithActor(id, amt, note, wallet.Address); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("released")
			return nil
		},
	}
	releaseCmd.Flags().String("wallet", "", "controller wallet path")
	releaseCmd.Flags().String("password", "", "wallet password")
	cmd.AddCommand(releaseCmd)

	authorizeCmd := &cobra.Command{
		Use:   "authorize <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Authorize an additional wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800Authorize")
			store, err := stage73State()
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			registry := store.Grants()
			if err := registry.AddAuthorizer(id, wallet.Address); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("authorized")
			return nil
		},
	}
	authorizeCmd.Flags().String("wallet", "", "wallet path")
	authorizeCmd.Flags().String("password", "", "wallet password")
	cmd.AddCommand(authorizeCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant details",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800Get")
			store, err := stage73State()
			if err != nil {
				return err
			}
			registry := store.Grants()
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			g, ok := registry.GetGrant(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.MarshalIndent(g, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800List")
			store, err := stage73State()
			if err != nil {
				return err
			}
			gs := store.Grants().ListGrants()
			b, _ := json.MarshalIndent(gs, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800Audit")
			store, err := stage73State()
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events, ok := store.Grants().AuditTrail(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.MarshalIndent(events, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(auditCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show grant telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3800Status")
			store, err := stage73State()
			if err != nil {
				return err
			}
			summary := store.Grants().Summary()
			b, _ := json.MarshalIndent(summary, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	rootCmd.AddCommand(cmd)
}
