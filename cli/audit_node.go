package cli

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			if err := auditNode.Start(); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	logCmd := &cobra.Command{
		Use:   "log [address] [event] [key=value]...",
		Args:  cobra.MinimumNArgs(2),
		Short: "Log event through audit node",
		Run: func(cmd *cobra.Command, args []string) {
			meta := make(map[string]string)
			for _, kv := range args[2:] {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) == 2 {
					meta[parts[0]] = parts[1]
				}
			}
			if err := auditNode.LogEvent(args[0], args[1], meta); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list [address]",
		Args:  cobra.ExactArgs(1),
		Short: "List events via audit node",
		Run: func(cmd *cobra.Command, args []string) {
			for _, e := range auditNode.ListEvents(args[0]) {
				fmt.Printf("%s %s %v\n", e.Timestamp.Format("2006-01-02T15:04:05"), e.Event, e.Metadata)
			}
		},
	}

	nodeCmd.AddCommand(startCmd, logCmd, listCmd)
	rootCmd.AddCommand(nodeCmd)
}
