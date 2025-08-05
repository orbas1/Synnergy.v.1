package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/nodes"
)

func init() {
	cmd := &cobra.Command{
		Use:   "nodeaddr",
		Short: "Node address utilities",
	}

	parseCmd := &cobra.Command{
		Use:   "parse <addr>",
		Short: "Parse a node address",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var a nodes.Address = nodes.Address(args[0])
			fmt.Println(a)
		},
	}
	cmd.AddCommand(parseCmd)

	validateCmd := &cobra.Command{
		Use:   "validate <addr>",
		Short: "Check if address is non-empty",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] == "" {
				fmt.Println("invalid")
			} else {
				fmt.Println("valid")
			}
		},
	}
	cmd.AddCommand(validateCmd)

	rootCmd.AddCommand(cmd)
}
