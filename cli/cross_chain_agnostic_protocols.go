package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var protocolRegistry = core.NewProtocolRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "cross_chain_agnostic_protocols",
		Short: "Register cross-chain protocols",
	}

	registerCmd := &cobra.Command{
		Use:   "register <name>",
		Args:  cobra.ExactArgs(1),
		Short: "Register a new protocol definition",
		Run: func(cmd *cobra.Command, args []string) {
			id := protocolRegistry.Register(args[0])
			fmt.Println(id)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered protocols",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range protocolRegistry.List() {
				fmt.Printf("%d: %s\n", p.ID, p.Name)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a protocol configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			p, ok := protocolRegistry.Get(id)
			if !ok {
				return fmt.Errorf("protocol not found")
			}
			fmt.Printf("%d: %s\n", p.ID, p.Name)
			return nil
		},
	}

	cmd.AddCommand(registerCmd, listCmd, getCmd)
	rootCmd.AddCommand(cmd)
}
