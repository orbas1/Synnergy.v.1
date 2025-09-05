package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	synn "synnergy"
	"synnergy/core"
)

func init() {
	cmd := &cobra.Command{
		Use:   "contractopcodes",
		Short: "List contract-related opcodes with gas costs",
		RunE: func(cmd *cobra.Command, args []string) error {
			entries := []struct {
				name string
				op   core.Opcode
			}{
				{"InitContracts", core.OpInitContracts},
				{"PauseContract", core.OpPauseContract},
				{"ResumeContract", core.OpResumeContract},
				{"UpgradeContract", core.OpUpgradeContract},
				{"ContractInfo", core.OpContractInfo},
				{"DeployAIContract", core.OpDeployAIContract},
				{"InvokeAIContract", core.OpInvokeAIContract},
			}
			for _, e := range entries {
				fmt.Fprintf(cmd.OutOrStdout(), "%d: %s (gas %d)\n", e.op, e.name, synn.GasCost(e.name))
			}
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
