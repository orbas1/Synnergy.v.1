package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var auditManager = core.NewAuditManager()

func init() {
	auditCmd := &cobra.Command{Use: "audit", Short: "Audit log management"}

	logCmd := &cobra.Command{
		Use:   "log [address] [event] [key=value]...",
		Args:  cobra.MinimumNArgs(2),
		Short: "Record an audit event",
		RunE: func(cmd *cobra.Command, args []string) error {
			meta := make(map[string]string)
			for _, kv := range args[2:] {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) == 2 {
					meta[parts[0]] = parts[1]
				}
			}
			return auditManager.Log(args[0], args[1], meta)
		},
	}

	var jsonOut bool
	listCmd := &cobra.Command{
		Use:   "list [address]",
		Args:  cobra.ExactArgs(1),
		Short: "List audit events for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries := auditManager.List(args[0])
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

	auditCmd.AddCommand(logCmd, listCmd)
	rootCmd.AddCommand(auditCmd)
}
