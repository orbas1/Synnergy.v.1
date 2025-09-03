package cli

import (
	"fmt"
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
			gas, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid gas")
				return
			}
			mem, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				fmt.Println("invalid memory limit")
				return
			}
			if _, err := sandboxMgr.StartSandbox(args[0], args[1], gas, mem); err != nil {
				fmt.Println("start error:", err)
				return
			}
			fmt.Println("sandbox started")
		},
	}
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Stop a sandbox",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sandboxMgr.StopSandbox(args[0]); err != nil {
				fmt.Println("stop error:", err)
				return
			}
			fmt.Println("sandbox stopped")
		},
	}
	cmd.AddCommand(stopCmd)

	deleteCmd := &cobra.Command{
		Use:   "delete <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a sandbox",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sandboxMgr.DeleteSandbox(args[0]); err != nil {
				fmt.Println("delete error:", err)
				return
			}
			fmt.Println("sandbox deleted")
		},
	}
	cmd.AddCommand(deleteCmd)

	resetCmd := &cobra.Command{
		Use:   "reset <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Reset sandbox timer",
		Run: func(cmd *cobra.Command, args []string) {
			if err := sandboxMgr.ResetSandbox(args[0]); err != nil {
				fmt.Println("reset error:", err)
				return
			}
			fmt.Println("sandbox reset")
		},
	}
	cmd.AddCommand(resetCmd)

	statusCmd := &cobra.Command{
		Use:   "status <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show sandbox status",
		Run: func(cmd *cobra.Command, args []string) {
			sb, ok := sandboxMgr.SandboxStatus(args[0])
			if !ok {
				fmt.Println("sandbox not found")
				return
			}
			fmt.Printf("%+v\n", sb)
		},
	}
	cmd.AddCommand(statusCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List sandboxes",
		Run: func(cmd *cobra.Command, args []string) {
			for _, sb := range sandboxMgr.ListSandboxes() {
				fmt.Printf("%+v\n", sb)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
