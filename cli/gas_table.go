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
			gasPrint("GasSet")
			op, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				printOutput("invalid opcode")
				return
			}
			cost, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				printOutput("invalid cost")
				return
			}
			core.SetGasCost(core.Opcode(op), cost)
			printOutput("gas cost updated")
		},
	}

	var jsonOut bool
	snapCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Print current gas table snapshot",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("GasSnapshot")
			if jsonOut {
				data, err := core.GasTableSnapshotJSON()
				if err != nil {
					printOutput(err.Error())
					return
				}
				fmt.Println(string(data))
				return
			}
			snapshot := core.GasTableSnapshot()
			printOutput(snapshot)
		},
	}
	snapCmd.Flags().BoolVar(&jsonOut, "json", false, "output JSON")

	gasCmd.AddCommand(setCmd)
	gasCmd.AddCommand(snapCmd)
}
