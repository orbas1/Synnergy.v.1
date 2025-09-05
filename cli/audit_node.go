package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

type dummyBootstrap struct{}

func (dummyBootstrap) Start() error {
	fmt.Println("bootstrap started")
	return nil
}

var auditNode = core.NewAuditNode(dummyBootstrap{}, auditManager)

func init() {
	nodeCmd := &cobra.Command{Use: "audit_node", Short: "Audit node operations"}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start bootstrap node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return auditNode.Start()
		},
	}

	logCmd := &cobra.Command{
		Use:   "log [address] [event] [key=value]...",
		Args:  cobra.MinimumNArgs(2),
		Short: "Log event through audit node",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				return err
			}
			meta := make(map[string]string)
			for _, kv := range args[2:] {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) == 2 {
					meta[parts[0]] = parts[1]
				}
			}
			return auditNode.LogEvent(addr.Hex(), args[1], meta)
		},
	}

	var jsonOut bool
	listCmd := &cobra.Command{
		Use:   "list [address]",
		Args:  cobra.ExactArgs(1),
		Short: "List events via audit node",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				return err
			}
			entries := auditNode.ListEvents(addr.Hex())
			if jsonOut {
				b, err := json.MarshalIndent(entries, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				return nil
			}
			for _, e := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s %v\n", e.Timestamp.Format("2006-01-02T15:04:05"), e.Event, e.Metadata)
			}
			return nil
		},
	}
	listCmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")

	nodeCmd.AddCommand(startCmd, logCmd, listCmd)
	rootCmd.AddCommand(nodeCmd)
}
