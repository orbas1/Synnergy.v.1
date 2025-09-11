package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var enforcer *core.ImmutabilityEnforcer

func init() {
	immCmd := &cobra.Command{
		Use:   "immutability",
		Short: "Immutability enforcement utilities",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise enforcer with ledger genesis",
		Run: func(cmd *cobra.Command, args []string) {
			gen := core.NewBlock(nil, "")
			gen.Hash = gen.HeaderHash(0)
			if err := ledger.AddBlock(gen); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			enforcer = core.NewImmutabilityEnforcer(gen)
			gasPrint("ImmutabilityInit")
			printOutput(map[string]any{"genesis": gen.Hash})
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Verify ledger against genesis hash",
		Run: func(cmd *cobra.Command, args []string) {
			if enforcer == nil {
				printOutput("enforcer not initialised")
				return
			}
			if err := enforcer.CheckLedger(ledger); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("ImmutabilityCheck")
			printOutput(map[string]any{"status": "ok"})
		},
	}

	immCmd.AddCommand(initCmd, checkCmd)
	rootCmd.AddCommand(immCmd)
}
