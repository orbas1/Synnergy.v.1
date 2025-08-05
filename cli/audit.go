package cli

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			meta := make(map[string]string)
			for _, kv := range args[2:] {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) == 2 {
					meta[parts[0]] = parts[1]
				}
			}
			if err := auditManager.Log(args[0], args[1], meta); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [address]",
		Args:  cobra.ExactArgs(1),
		Short: "List audit events for an address",
		Run: func(cmd *cobra.Command, args []string) {
			entries := auditManager.List(args[0])
			for _, e := range entries {
				fmt.Printf("%s %s %v\n", e.Timestamp.Format("2006-01-02T15:04:05"), e.Event, e.Metadata)
			}
		},
	}

	auditCmd.AddCommand(logCmd, listCmd)
	rootCmd.AddCommand(auditCmd)
}
