package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var idRegistry = core.NewIDRegistry()

func init() {
	idwCmd := &cobra.Command{
		Use:   "idwallet",
		Short: "ID wallet registration",
	}

	registerCmd := &cobra.Command{
		Use:   "register [address] [info]",
		Args:  cobra.ExactArgs(2),
		Short: "Register an ID wallet",
		Run: func(cmd *cobra.Command, args []string) {
			if err := idRegistry.Register(args[0], args[1]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("IDWalletRegister")
			printOutput(map[string]any{"status": "registered", "address": args[0]})
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Check registration info",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := idRegistry.Info(args[0])
			if !ok {
				printOutput(map[string]any{"error": "not registered"})
				return
			}
			gasPrint("IDWalletCheck")
			printOutput(map[string]any{"info": info})
		},
	}

	idwCmd.AddCommand(registerCmd, checkCmd)
	rootCmd.AddCommand(idwCmd)
}
