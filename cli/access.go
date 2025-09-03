package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	accessCmd := &cobra.Command{Use: "access", Short: "Role based access control"}

	grantCmd := &cobra.Command{
		Use:   "grant [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Grant a role to an address",
		Run: func(cmd *cobra.Command, args []string) {
			if err := core.GrantRole(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("granted")
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Revoke a role from an address",
		Run: func(cmd *cobra.Command, args []string) {
			if err := core.RevokeRole(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("revoked")
		},
	}

	hasCmd := &cobra.Command{
		Use:   "has [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Check if address has role",
		Run: func(cmd *cobra.Command, args []string) {
			ok, err := core.HasRole(args[0], args[1])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(ok)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "List roles for an address",
		Run: func(cmd *cobra.Command, args []string) {
			roles, err := core.ListRoles(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			for _, r := range roles {
				fmt.Println(r)
			}
		},
	}

	accessCmd.AddCommand(grantCmd, revokeCmd, hasCmd, listCmd)
	rootCmd.AddCommand(accessCmd)
}
