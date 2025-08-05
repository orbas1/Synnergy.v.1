package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var bridgeRegistry = core.NewBridgeRegistry()

func init() {
	ccCmd := &cobra.Command{
		Use:   "cross_chain",
		Short: "Manage cross-chain bridges",
	}

	registerCmd := &cobra.Command{
		Use:   "register <source_chain> <target_chain> <relayer_addr>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a bridge",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := bridgeRegistry.RegisterBridge(args[0], args[1], args[2])
			if err != nil {
				return err
			}
			fmt.Println(b.ID)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered bridges",
		Run: func(cmd *cobra.Command, args []string) {
			for _, b := range bridgeRegistry.ListBridges() {
				fmt.Printf("%s: %s -> %s\n", b.ID, b.SourceChain, b.TargetChain)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <bridge_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a bridge configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, ok := bridgeRegistry.GetBridge(args[0])
			if !ok {
				return fmt.Errorf("bridge not found")
			}
			fmt.Printf("%s: %s -> %s relayers=%d\n", b.ID, b.SourceChain, b.TargetChain, len(b.Relayers))
			return nil
		},
	}

	authorizeCmd := &cobra.Command{
		Use:   "authorize <bridge_id> <relayer_addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Whitelist a relayer address",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeRegistry.AuthorizeRelayer(args[0], args[1])
		},
	}

	revokeCmd := &cobra.Command{
		Use:   "revoke <bridge_id> <relayer_addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Remove a relayer from the whitelist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridgeRegistry.RevokeRelayer(args[0], args[1])
		},
	}

	ccCmd.AddCommand(registerCmd, listCmd, getCmd, authorizeCmd, revokeCmd)
	rootCmd.AddCommand(ccCmd)
}
