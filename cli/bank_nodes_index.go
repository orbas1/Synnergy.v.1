package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	bankNodesCmd := &cobra.Command{
		Use:   "banknodes",
		Short: "Bank node type utilities",
	}

	typesCmd := &cobra.Command{
		Use:   "types",
		Short: "List supported bank node types",
		Run: func(cmd *cobra.Command, args []string) {
			for _, t := range core.BankNodeTypes {
				fmt.Println(t)
			}
		},
	}

	bankNodesCmd.AddCommand(typesCmd)
	rootCmd.AddCommand(bankNodesCmd)
}
