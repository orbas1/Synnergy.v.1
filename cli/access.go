package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var accessCtrl = core.NewAccessController()

func init() {
	accessCmd := &cobra.Command{Use: "access", Short: "Role based access control"}

	grantCmd := &cobra.Command{
		Use:   "grant [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Grant a role to an address",
		Run: func(cmd *cobra.Command, args []string) {
			accessCtrl.Grant(args[0], args[1])
			fmt.Println("granted")
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Revoke a role from an address",
		Run: func(cmd *cobra.Command, args []string) {
			accessCtrl.Revoke(args[0], args[1])
			fmt.Println("revoked")
		},
	}

	hasCmd := &cobra.Command{
		Use:   "has [role] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Check if address has role",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(accessCtrl.HasRole(args[0], args[1]))
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "List roles for an address",
		Run: func(cmd *cobra.Command, args []string) {
			roles := accessCtrl.List(args[0])
			for _, r := range roles {
				fmt.Println(r)
			}
		},
	}

	accessCmd.AddCommand(grantCmd, revokeCmd, hasCmd, listCmd)
	rootCmd.AddCommand(accessCmd)
}
