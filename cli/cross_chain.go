package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	synnergy "synnergy"
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
			fmt.Fprintf(cmd.OutOrStdout(), "%s gas:%d\n", b.ID, synnergy.GasCost("RegisterBridge"))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered bridges",
		Run: func(cmd *cobra.Command, args []string) {
			bridges := bridgeRegistry.ListBridges()
			listJSON, _ := cmd.Flags().GetBool("json")
			if listJSON {
				enc, err := json.Marshal(bridges)
				if err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err)
					return
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(enc))
				return
			}
			for _, b := range bridges {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s -> %s\n", b.ID, b.SourceChain, b.TargetChain)
			}
		},
	}
	listCmd.Flags().Bool("json", false, "output as JSON")

	getCmd := &cobra.Command{
		Use:   "get <bridge_id>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a bridge configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, ok := bridgeRegistry.GetBridge(args[0])
			if !ok {
				return fmt.Errorf("bridge not found")
			}
			getJSON, _ := cmd.Flags().GetBool("json")
			if getJSON {
				enc, err := json.Marshal(b)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(enc))
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s -> %s relayers=%d\n", b.ID, b.SourceChain, b.TargetChain, len(b.Relayers))
			return nil
		},
	}
	getCmd.Flags().Bool("json", false, "output as JSON")

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
