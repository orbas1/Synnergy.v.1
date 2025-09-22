package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var grantRegistry = core.NewGrantRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3800",
		Short: "Manage SYN3800 grant records",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureGrantRegistryLoaded()
		},
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			authorizerSpec, _ := cmd.Flags().GetString("authorizer")
			authorizer := ""
			if authorizerSpec != "" {
				addr, err := walletAddressFromSpec(authorizerSpec)
				if err != nil {
					return err
				}
				authorizer = addr
			}
			id, err := grantRegistry.CreateGrant(args[0], args[1], amt, core.Address(authorizer))
			if err != nil {
				return err
			}
			if err := persistGrantRegistry(); err != nil {
				return err
			}
			cmd.Println(id)
			return nil
		},
	}
	createCmd.Flags().String("authorizer", "", "wallet:path:password of initial authorizer")
	cmd.AddCommand(createCmd)

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
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			if err := grantRegistry.Disburse(id, amt, note, core.Address(wallet.Address)); err != nil {
				return err
			}
			if err := persistGrantRegistry(); err != nil {
				return err
			}
			printOutput("released")
			return nil
		},
	}
	addWalletFlags(releaseCmd)
	cmd.AddCommand(releaseCmd)

	authorizeCmd := &cobra.Command{
		Use:   "authorize <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Authorize a wallet to release grant funds",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := grantRegistry.Authorize(id, core.Address(wallet.Address)); err != nil {
				return err
			}
			if err := persistGrantRegistry(); err != nil {
				return err
			}
			printOutput("authorized")
			return nil
		},
	}
	addWalletFlags(authorizeCmd)
	cmd.AddCommand(authorizeCmd)

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
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		Run: func(cmd *cobra.Command, args []string) {
			gs := grantRegistry.ListGrants()
			b, _ := json.MarshalIndent(gs, "", "  ")
			cmd.Println(string(b))
		},
	}
	cmd.AddCommand(listCmd)

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show audit log for a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events := grantRegistry.AuditLog(id)
			b, _ := json.MarshalIndent(events, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(auditCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show grant registry status",
		Run: func(cmd *cobra.Command, args []string) {
			summary := grantRegistry.StatusSummary()
			b, _ := json.MarshalIndent(summary, "", "  ")
			cmd.Println(string(b))
		},
	}
	cmd.AddCommand(statusCmd)

	rootCmd.AddCommand(cmd)
}

func walletAddressFromSpec(spec string) (string, error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("wallet spec must be path:password")
	}
	w, err := loadWallet(parts[0], parts[1])
	if err != nil {
		return "", err
	}
	return w.Address, nil
}

func addWalletFlags(cmd *cobra.Command) {
	cmd.Flags().String("wallet", "", "path to wallet file")
	cmd.Flags().String("password", "", "wallet password")
}

func requireWalletFlags(cmd *cobra.Command) {
	prev := cmd.PreRunE
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if prev != nil {
			if err := prev(cmd, args); err != nil {
				return err
			}
		}
		path, _ := cmd.Flags().GetString("wallet")
		pass, _ := cmd.Flags().GetString("password")
		if path == "" || pass == "" {
			return fmt.Errorf("wallet and password required")
		}
		return nil
	}
}

func walletFromFlags(cmd *cobra.Command) (*core.Wallet, error) {
	path, _ := cmd.Flags().GetString("wallet")
	pass, _ := cmd.Flags().GetString("password")
	if path == "" || pass == "" {
		return nil, fmt.Errorf("wallet and password required")
	}
	return loadWallet(path, pass)
}
