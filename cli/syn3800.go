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
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			beneficiary := args[0]
			if beneficiary == "" {
				return fmt.Errorf("beneficiary required")
			}
			name := args[1]
			if name == "" {
				return fmt.Errorf("name required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			authorizerFlag, _ := cmd.Flags().GetString("authorizer")
			var authorizerAddr string
			if authorizerFlag != "" {
				path, password, err := parseWalletCredential(authorizerFlag)
				if err != nil {
					return err
				}
				wallet, err := loadWallet(path, password)
				if err != nil {
					return err
				}
				authorizerAddr = wallet.Address
			}
			id := grantRegistry.CreateGrant(beneficiary, name, amt)
			if authorizerAddr != "" {
				_ = grantRegistry.Authorize(id, authorizerAddr)
			}
			cmd.Println(id)
			return nil
		},
	}
	createCmd.Flags().String("authorizer", "", "wallet credentials for initial authorizer path:password")

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
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			if err := grantRegistry.Disburse(id, amt, note, wallet.Address); err != nil {
				return err
			}
			cmd.Println("released")
			return nil
		},
	}
	releaseCmd.Flags().String("wallet", "", "path to wallet file")
	releaseCmd.Flags().String("password", "", "wallet password")

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
			payload, _ := json.MarshalIndent(g, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		Run: func(cmd *cobra.Command, args []string) {
			gs := grantRegistry.ListGrants()
			payload, _ := json.MarshalIndent(gs, "", "  ")
			cmd.Println(string(payload))
		},
	}

	authorizeCmd := &cobra.Command{
		Use:   "authorize <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Authorize a wallet to release funds",
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
			cmd.Println("authorized")
			return nil
		},
	}
	authorizeCmd.Flags().String("wallet", "", "path to wallet file")
	authorizeCmd.Flags().String("password", "", "wallet password")

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Audit log for a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events, ok := grantRegistry.Audit(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			payload, _ := json.MarshalIndent(events, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show grant telemetry",
		Run: func(cmd *cobra.Command, args []string) {
			tele := grantRegistry.Telemetry()
			payload, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(payload))
		},
	}

	cmd.AddCommand(createCmd, releaseCmd, getCmd, listCmd, authorizeCmd, auditCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
