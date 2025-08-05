package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var crossContractRegistry = core.NewCrossChainRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "xcontract",
		Short: "Register cross-chain contract mappings",
	}

	registerCmd := &cobra.Command{
		Use:   "register <local_addr> <remote_chain> <remote_addr>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a contract mapping",
		Run: func(cmd *cobra.Command, args []string) {
			crossContractRegistry.RegisterMapping(args[0], args[1], args[2])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered mappings",
		Run: func(cmd *cobra.Command, args []string) {
			for _, m := range crossContractRegistry.ListMappings() {
				fmt.Printf("%s -> %s:%s\n", m.LocalAddress, m.RemoteChain, m.RemoteAddress)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <local_addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve mapping info",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, ok := crossContractRegistry.GetMapping(args[0])
			if !ok {
				return fmt.Errorf("mapping not found")
			}
			fmt.Printf("%s -> %s:%s\n", m.LocalAddress, m.RemoteChain, m.RemoteAddress)
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <local_addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			return crossContractRegistry.RemoveMapping(args[0])
		},
	}

	cmd.AddCommand(registerCmd, listCmd, getCmd, removeCmd)
	rootCmd.AddCommand(cmd)
}
