package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	govCmd := &cobra.Command{
		Use:   "government",
		Short: "Government authority node operations",
	}

	newCmd := &cobra.Command{
		Use:   "new [address] [role] [department]",
		Args:  cobra.ExactArgs(3),
		Short: "Create a government authority node",
		Run: func(cmd *cobra.Command, args []string) {
			n := core.NewGovernmentAuthorityNode(args[0], args[1], args[2])
			fmt.Printf("node address=%s role=%s department=%s\n", n.Address, n.Role, n.Department)
		},
	}

	govCmd.AddCommand(newCmd)
	rootCmd.AddCommand(govCmd)
}
