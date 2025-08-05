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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%d: OpInitContracts\n", core.OpInitContracts)
			fmt.Printf("%d: OpPauseContract\n", core.OpPauseContract)
			fmt.Printf("%d: OpResumeContract\n", core.OpResumeContract)
			fmt.Printf("%d: OpUpgradeContract\n", core.OpUpgradeContract)
			fmt.Printf("%d: OpContractInfo\n", core.OpContractInfo)
			fmt.Printf("%d: OpDeployAIContract\n", core.OpDeployAIContract)
			fmt.Printf("%d: OpInvokeAIContract\n", core.OpInvokeAIContract)
		},
	}
	rootCmd.AddCommand(cmd)
}
