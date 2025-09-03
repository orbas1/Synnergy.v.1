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

	var jsonOut bool
	listCmd := &cobra.Command{
		Use:   "list [address]",
		Args:  cobra.ExactArgs(1),
		Short: "List audit events for an address",
		Run: func(cmd *cobra.Command, args []string) {
			entries := auditManager.List(args[0])
			if jsonOut {
				b, _ := json.MarshalIndent(entries, "", "  ")
				fmt.Println(string(b))
				return
			}
			for _, e := range entries {
				fmt.Printf("%s %s %v\n", e.Timestamp.Format("2006-01-02T15:04:05"), e.Event, e.Metadata)
			}
		},
	}
	listCmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")

	auditCmd.AddCommand(logCmd, listCmd)
	rootCmd.AddCommand(auditCmd)
}
