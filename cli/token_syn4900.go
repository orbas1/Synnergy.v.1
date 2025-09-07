package cli

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var agriRegistry = core.NewAgriculturalRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "token_syn4900",
		Short: "Manage SYN4900 agricultural assets",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <type> <owner> <origin> <qty> <harvest> <expiry> <cert>",
		Args:  cobra.ExactArgs(8),
		Short: "Register an agricultural asset",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Syn4900Register")
			qty, _ := strconv.ParseUint(args[4], 10, 64)
			harvest, _ := strconv.ParseInt(args[5], 10, 64)
			expiry, _ := strconv.ParseInt(args[6], 10, 64)
			if _, err := agriRegistry.Register(args[0], args[1], args[2], args[3], qty, time.Unix(harvest, 0), time.Unix(expiry, 0), args[7]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "registered", "id": args[0]})
		},
	}

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <owner>",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer ownership",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Syn4900Transfer")
			if err := agriRegistry.Transfer(args[0], args[1]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "transferred", "id": args[0], "owner": args[1]})
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <id> <status>",
		Args:  cobra.ExactArgs(2),
		Short: "Update asset status",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Syn4900Status")
			if err := agriRegistry.UpdateStatus(args[0], args[1]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "updated", "id": args[0], "state": args[1]})
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show asset info",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Syn4900Info")
			if a, ok := agriRegistry.Get(args[0]); ok {
				printOutput(a)
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	}

	cmd.AddCommand(registerCmd, transferCmd, statusCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
