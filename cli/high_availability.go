package cli

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var failover *core.FailoverManager

func init() {
	haCmd := &cobra.Command{
		Use:   "highavailability",
		Short: "High availability failover management",
	}

	initCmd := &cobra.Command{
		Use:   "init [primary] [timeoutSec]",
		Args:  cobra.ExactArgs(2),
		Short: "Initialise failover manager",
		Run: func(cmd *cobra.Command, args []string) {
			t, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				printOutput("invalid timeout: " + err.Error())
				return
			}
			failover = core.NewFailoverManager(args[0], time.Duration(t)*time.Second)
			gasPrint("FailoverInit")
			printOutput(map[string]any{"status": "initialised", "primary": args[0], "timeout": t})
		},
	}

	addCmd := &cobra.Command{
		Use:   "add [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Register a backup node",
		Run: func(cmd *cobra.Command, args []string) {
			if failover == nil {
				printOutput("manager not initialised")
				return
			}
			failover.RegisterBackup(args[0])
			gasPrint("FailoverRegister")
			printOutput(map[string]any{"status": "registered", "id": args[0]})
		},
	}

	heartbeatCmd := &cobra.Command{
		Use:   "heartbeat [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Record a heartbeat",
		Run: func(cmd *cobra.Command, args []string) {
			if failover == nil {
				printOutput("manager not initialised")
				return
			}
			failover.Heartbeat(args[0])
			gasPrint("FailoverHeartbeat")
			printOutput(map[string]any{"status": "heartbeat", "id": args[0]})
		},
	}

	activeCmd := &cobra.Command{
		Use:   "active",
		Short: "Show active node",
		Run: func(cmd *cobra.Command, args []string) {
			if failover == nil {
				printOutput("manager not initialised")
				return
			}
			gasPrint("FailoverActive")
			printOutput(map[string]any{"active": failover.Active()})
		},
	}

	haCmd.AddCommand(initCmd, addCmd, heartbeatCmd, activeCmd)
	rootCmd.AddCommand(haCmd)
}
