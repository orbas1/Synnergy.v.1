package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var grantRegistry = core.NewGrantRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3800",
		Short: "Manage SYN3800 grants",
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			beneficiary := args[0]
			name := args[1]
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			if beneficiary == "" {
				return fmt.Errorf("beneficiary required")
			}
			if name == "" {
				return fmt.Errorf("name required")
			}
			authorizerSpec, _ := cmd.Flags().GetString("authorizer")
			if authorizerSpec == "" {
				return fmt.Errorf("authorizer descriptor required")
			}
			path, password, err := parseWalletDescriptor(authorizerSpec)
			if err != nil {
				return err
			}
			wallet, err := loadWallet(path, password)
			if err != nil {
				return err
			}
			id, err := grantRegistry.CreateGrantWithAuthorizer(beneficiary, name, amt, wallet.Address)
			if err != nil {
				return err
			}
			gasPrint("CreateGrant")
			cmd.Println(id)
			return nil
		},
	}
	createCmd.Flags().String("authorizer", "", "wallet descriptor path:password for the initial authorizer")
	cmd.AddCommand(createCmd)

	releaseCmd := &cobra.Command{
		Use:   "release <id> <amount> [note]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Release grant funds",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			amount, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || amount == 0 {
				return fmt.Errorf("invalid amount")
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := grantRegistry.DisburseWithAuthorizer(id, amount, note, wallet.Address); err != nil {
				return err
			}
			gasPrint("ReleaseGrant")
			cmd.Println("released")
			return nil
		},
	}
	releaseCmd.Flags().String("wallet", "", "authorizer wallet path")
	releaseCmd.Flags().String("password", "", "authorizer wallet password")
	cmd.AddCommand(releaseCmd)

	authorizeCmd := &cobra.Command{
		Use:   "authorize <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Authorize an additional wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := grantRegistry.AuthorizeSigner(id, wallet.Address); err != nil {
				return err
			}
			gasPrint("AuthorizeGrantWallet")
			cmd.Println("authorized")
			return nil
		},
	}
	authorizeCmd.Flags().String("wallet", "", "wallet path to authorise")
	authorizeCmd.Flags().String("password", "", "wallet password")
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
			grant, ok := grantRegistry.GetGrant(id)
			if !ok {
				return fmt.Errorf("grant not found")
			}
			return printJSON(cmd, grant)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		RunE: func(cmd *cobra.Command, args []string) error {
			grants := grantRegistry.ListGrants()
			return printJSON(cmd, grants)
		},
	}
	cmd.AddCommand(listCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Summarise grant lifecycle counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			counts := grantRegistry.StatusSummary()
			out := struct {
				Total     int `json:"total"`
				Pending   int `json:"pending"`
				Active    int `json:"active"`
				Completed int `json:"completed"`
			}{
				Total:     counts["total"],
				Pending:   counts[string(core.GrantStatusPending)],
				Active:    counts[string(core.GrantStatusActive)],
				Completed: counts[string(core.GrantStatusCompleted)],
			}
			return printJSON(cmd, out)
		},
	}
	cmd.AddCommand(statusCmd)

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant audit events",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events, err := grantRegistry.AuditTrail(id)
			if err != nil {
				return err
			}
			return printJSON(cmd, events)
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}
