package cli

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeAddrParse")
			var a nodes.Address = nodes.Address(args[0])
			printOutput(map[string]string{"parsed": string(a)})
			return nil
		},
	}
	cmd.AddCommand(parseCmd)

	validateCmd := &cobra.Command{
		Use:   "validate <addr>",
		Short: "Check if address is non-empty",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeAddrValidate")
			valid := args[0] != ""
			printOutput(map[string]bool{"valid": valid})
			return nil
		},
	}
	cmd.AddCommand(validateCmd)

	rootCmd.AddCommand(cmd)
}
