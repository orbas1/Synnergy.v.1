package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	setCmd := &cobra.Command{
		Use:   "set <opcode> <cost>",
		Short: "Set gas cost for an opcode",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			op, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				fmt.Printf("invalid opcode: %v\n", err)
				return
			}
			cost, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Printf("invalid cost: %v\n", err)
				return
			}
			core.SetGasCost(core.Opcode(op), cost)
			fmt.Println("gas cost updated")
		},
	}

	snapCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Print current gas table snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			snapshot := core.GasTableSnapshot()
			for op, cost := range snapshot {
				fmt.Printf("%v: %d\n", op, cost)
			}
		},
	}

	gasCmd.AddCommand(setCmd)
	gasCmd.AddCommand(snapCmd)
}
