package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var gasTable = core.DefaultGasTable()

func init() {
	gasCmd := &cobra.Command{
		Use:   "gas",
		Short: "Interact with gas table",
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List gas costs",
		Run: func(cmd *cobra.Command, args []string) {
			for op, cost := range gasTable {
				fmt.Printf("%v: %d\n", op, cost)
			}
		},
	}
	gasCmd.AddCommand(listCmd)
	rootCmd.AddCommand(gasCmd)
}
