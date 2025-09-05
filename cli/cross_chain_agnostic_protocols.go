package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var protocolRegistry = core.NewProtocolRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "cross_chain_agnostic_protocols",
		Short: "Manage cross-chain protocols",
	}

	registerCmd := &cobra.Command{
		Use:   "register <name>",
		Args:  cobra.ExactArgs(1),
		Short: "Register a new protocol definition",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := protocolRegistry.Register(args[0])
			fmt.Fprintf(cmd.OutOrStdout(), "%d gas:%d\n", id, synnergy.GasCost("RegisterProtocol"))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered protocols",
		Run: func(cmd *cobra.Command, args []string) {
			ps := protocolRegistry.List()
			listJSON, _ := cmd.Flags().GetBool("json")
			if listJSON {
				enc, err := json.Marshal(ps)
				if err != nil {
					fmt.Fprintln(cmd.ErrOrStderr(), err)
					return
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(enc))
				return
			}
			for _, p := range ps {
				fmt.Fprintf(cmd.OutOrStdout(), "%d: %s\n", p.ID, p.Name)
			}
		},
	}
	listCmd.Flags().Bool("json", false, "output as JSON")

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
			getJSON, _ := cmd.Flags().GetBool("json")
			if getJSON {
				enc, err := json.Marshal(p)
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(enc))
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%d: %s\n", p.ID, p.Name)
			return nil
		},
	}
	getCmd.Flags().Bool("json", false, "output as JSON")

	cmd.AddCommand(registerCmd, listCmd, getCmd)
	rootCmd.AddCommand(cmd)
}
