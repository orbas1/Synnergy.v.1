package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var sysLogger = core.NewSystemHealthLogger()
var sysLogs []string

func init() {
	cmd := &cobra.Command{
		Use:   "system_health",
		Short: "System health utilities",
	}

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Display current system metrics",
		Run: func(cmd *cobra.Command, args []string) {
			m := sysLogger.Collect(0, 0)
			b, _ := json.MarshalIndent(m, "", "  ")
			fmt.Println(string(b))
		},
	}

	logCmd := &cobra.Command{
		Use:   "log <level> <msg>",
		Args:  cobra.ExactArgs(2),
		Short: "Append a log message",
		Run: func(cmd *cobra.Command, args []string) {
			sysLogs = append(sysLogs, fmt.Sprintf("%s: %s", args[0], args[1]))
		},
	}

	cmd.AddCommand(snapshotCmd, logCmd)
	rootCmd.AddCommand(cmd)
}
