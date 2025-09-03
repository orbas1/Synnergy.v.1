package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	cmd := &cobra.Command{
		Use:   "contractopcodes",
		Short: "List contract-related opcodes",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpInitContracts\n", core.OpInitContracts)
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpPauseContract\n", core.OpPauseContract)
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpResumeContract\n", core.OpResumeContract)
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpUpgradeContract\n", core.OpUpgradeContract)
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpContractInfo\n", core.OpContractInfo)
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpDeployAIContract\n", core.OpDeployAIContract)
			fmt.Fprintf(cmd.OutOrStdout(), "%d: OpInvokeAIContract\n", core.OpInvokeAIContract)
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
