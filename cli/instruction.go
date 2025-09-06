package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func parseOpcode(arg string) (core.Opcode, error) {
	if n, err := strconv.ParseUint(arg, 10, 32); err == nil {
		return core.Opcode(n), nil
	}
	for op, name := range core.Opcodes() {
		if strings.EqualFold(name, arg) {
			return op, nil
		}
	}
	return 0, fmt.Errorf("unknown opcode %s", arg)
}

func init() {
	cmd := &cobra.Command{
		Use:   "instruction",
		Short: "Work with VM instructions",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "new [opcode] [value]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Create an instruction",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("InstructionNew")
			op, err := parseOpcode(args[0])
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			var val int64
			if len(args) == 2 {
				v, err := strconv.ParseInt(args[1], 10, 64)
				if err != nil {
					printOutput("invalid value")
					return
				}
				val = v
			}
			inst := core.Instruction{Op: op, Value: val}
			printOutput(inst)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List registered opcodes",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("InstructionList")
			cat := core.Catalogue()
			printOutput(cat)
		},
	})

	rootCmd.AddCommand(cmd)
}
