package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var sandboxMgr = core.NewSandboxManager()

func init() {
	cmd := &cobra.Command{
		Use:   "sandbox",
		Short: "Manage VM sandboxes",
	}

	startCmd := &cobra.Command{
		Use:   "start <id> <contract> <gas> <memory>",
		Args:  cobra.ExactArgs(4),
		Short: "Start a sandbox",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxStart")
			gas, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid gas"})
				return
			}
			mem, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid memory"})
				return
			}
			if _, err := sandboxMgr.StartSandbox(args[0], args[1], gas, mem); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "started", "id": args[0]})
		},
	}
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Stop a sandbox",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxStop")
			if err := sandboxMgr.StopSandbox(args[0]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "stopped", "id": args[0]})
		},
	}
	cmd.AddCommand(stopCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a sandbox",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxDelete")
			if err := sandboxMgr.DeleteSandbox(args[0]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "deleted", "id": args[0]})
		},
	}
	cmd.AddCommand(deleteCmd)

	resetCmd := &cobra.Command{
		Use:   "reset <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Reset sandbox timer",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxReset")
			if err := sandboxMgr.ResetSandbox(args[0]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "reset", "id": args[0]})
		},
	}
	cmd.AddCommand(resetCmd)

	statusCmd := &cobra.Command{
		Use:   "status <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show sandbox status",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxStatus")
			sb, ok := sandboxMgr.SandboxStatus(args[0])
			if !ok {
				printOutput(map[string]any{"error": "not found"})
				return
			}
			printOutput(sb)
		},
	}
	cmd.AddCommand(statusCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List sandboxes",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxList")
			printOutput(sandboxMgr.ListSandboxes())
		},
	}
	cmd.AddCommand(listCmd)

	purgeCmd := &cobra.Command{
		Use:   "purge",
		Short: "Remove stopped sandboxes past TTL",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SandboxPurge")
			sandboxMgr.PurgeInactive()
			printOutput(map[string]any{"status": "purged"})
		},
	}
	cmd.AddCommand(purgeCmd)

	rootCmd.AddCommand(cmd)
}
