package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var grantRegistry = core.NewGrantRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3800",
		Short: "Manage SYN3800 grant records",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureStage73Loaded()
		},
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" || args[1] == "" {
				return fmt.Errorf("beneficiary and name required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			var authorizer string
			authSpec, _ := cmd.Flags().GetString("authorizer")
			if authSpec != "" {
				wallet, err := loadWalletSpec(authSpec)
				if err != nil {
					return err
				}
				authorizer = wallet.Address
			}
			id := grantRegistry.CreateGrant(args[0], args[1], amt, authorizer)
			if err := persistStage73(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}
	createCmd.Flags().String("authorizer", "", "Wallet path:password granting release permissions")

	releaseCmd := &cobra.Command{
		Use:   "release <id> <amount> [note]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Release funds for a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := grantRegistry.Disburse(id, amt, note, wallet.Address); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("released")
			return nil
		},
	}
	releaseCmd.Flags().String("wallet", "", "Path to signing wallet")
	releaseCmd.Flags().String("password", "", "Wallet password")

	authorizeCmd := &cobra.Command{
		Use:   "authorize <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Authorize an additional wallet for releases",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := grantRegistry.Authorize(id, wallet.Address); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("authorized")
			return nil
		},
	}
	authorizeCmd.Flags().String("wallet", "", "Path to wallet to authorize")
	authorizeCmd.Flags().String("password", "", "Wallet password")

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant audit trail",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events, err := grantRegistry.Audit(id)
			if err != nil {
				return err
			}
			b, _ := json.MarshalIndent(events, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant details",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			g, ok := grantRegistry.GetGrant(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.MarshalIndent(g, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			gs := grantRegistry.ListGrants()
			b, _ := json.MarshalIndent(gs, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show grant telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			tele := grantRegistry.Telemetry()
			b, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	cmd.AddCommand(createCmd, releaseCmd, authorizeCmd, auditCmd, getCmd, listCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
