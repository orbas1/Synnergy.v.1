package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var identitySvc = core.NewIdentityService()

func init() {
	idCmd := &cobra.Command{
		Use:   "identity",
		Short: "Identity verification service",
	}

	registerCmd := &cobra.Command{
		Use:   "register [addr] [name] [dob] [nationality]",
		Args:  cobra.ExactArgs(4),
		Short: "Register identity information",
		Run: func(cmd *cobra.Command, args []string) {
			if err := identitySvc.Register(args[0], args[1], args[2], args[3]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("IdentityRegister")
			printOutput(map[string]any{"status": "registered", "address": args[0]})
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [addr] [method]",
		Args:  cobra.ExactArgs(2),
		Short: "Record a verification method",
		Run: func(cmd *cobra.Command, args []string) {
			if err := identitySvc.Verify(args[0], args[1]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("IdentityVerify")
			printOutput(map[string]any{"status": "verified", "address": args[0], "method": args[1]})
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show identity information",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := identitySvc.Info(args[0])
			if !ok {
				printOutput(map[string]any{"error": "not found"})
				return
			}
			gasPrint("IdentityInfo")
			printOutput(map[string]any{"name": info.Name, "dob": info.DateOfBirth, "nationality": info.Nationality})
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show verification logs",
		Run: func(cmd *cobra.Command, args []string) {
			logs := identitySvc.Logs(args[0])
			gasPrint("IdentityLogs")
			printOutput(logs)
		},
	}

	idCmd.AddCommand(registerCmd, verifyCmd, infoCmd, logsCmd)
	rootCmd.AddCommand(idCmd)
}
