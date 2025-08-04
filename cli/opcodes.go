package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	opcodeCmd := &cobra.Command{
		Use:   "opcodes",
		Short: "List supported opcodes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%d: OpNoop\n", core.OpNoop)
			fmt.Printf("%d: OpTransfer\n", core.OpTransfer)
		},
	}
	rootCmd.AddCommand(opcodeCmd)
}
