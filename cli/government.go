package cli

import (
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
			gasPrint("GovernmentNew")
			n := core.NewGovernmentAuthorityNode(args[0], args[1], args[2])
			printOutput(map[string]string{
				"address":    n.Address,
				"role":       n.Role,
				"department": n.Department,
			})
		},
	}

	govCmd.AddCommand(newCmd)
	rootCmd.AddCommand(govCmd)
}
