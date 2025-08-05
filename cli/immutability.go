package cli

import (
	"fmt"

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
			ledger.AddBlock(gen)
			enforcer = core.NewImmutabilityEnforcer(gen)
			fmt.Println("genesis hash", gen.Hash)
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Verify ledger against genesis hash",
		Run: func(cmd *cobra.Command, args []string) {
			if enforcer == nil {
				fmt.Println("enforcer not initialised")
				return
			}
			if err := enforcer.CheckLedger(ledger); err != nil {
				fmt.Println("check failed:", err)
				return
			}
			fmt.Println("ledger matches genesis hash")
		},
	}

	immCmd.AddCommand(initCmd, checkCmd)
	rootCmd.AddCommand(immCmd)
}
