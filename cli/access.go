package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	// Stage 38 hardens role-based access control commands with explicit success messages.
	accessCmd := &cobra.Command{Use: "access", Short: "Role based access control"}

	grantCmd := &cobra.Command{
		Use:   "grant [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Grant a role to an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.GrantRole(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "granted")
			return nil
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Revoke a role from an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.RevokeRole(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "revoked")
			return nil
		},
	}

	hasCmd := &cobra.Command{
		Use:   "has [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Check if address has role",
		RunE: func(cmd *cobra.Command, args []string) error {
			ok, err := core.HasRole(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), ok)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "List roles for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			roles, err := core.ListRoles(args[0])
			if err != nil {
				return err
			}
			for _, r := range roles {
				fmt.Fprintln(cmd.OutOrStdout(), r)
			}
			return nil
		},
	}

	accessCmd.AddCommand(grantCmd, revokeCmd, hasCmd, listCmd)
	rootCmd.AddCommand(accessCmd)
}
