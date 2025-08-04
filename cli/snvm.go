package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var vm = core.NewSNVM()

func init() {
	snvmCmd := &cobra.Command{
		Use:   "snvm",
		Short: "Interact with the Synnergy VM",
	}

	execCmd := &cobra.Command{
		Use:   "exec [add|sub|mul|div] <a> <b>",
		Short: "Execute a simple arithmetic program",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			opStr, aStr, bStr := args[0], args[1], args[2]
			a, err := strconv.ParseInt(aStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid operand a: %w", err)
			}
			b, err := strconv.ParseInt(bStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid operand b: %w", err)
			}
			var op core.Opcode
			switch opStr {
			case "add":
				op = core.OpAdd
			case "sub":
				op = core.OpSub
			case "mul":
				op = core.OpMul
			case "div":
				op = core.OpDiv
			default:
				return fmt.Errorf("unknown opcode %s", opStr)
			}

			tx := core.NewTransaction("", "", 0, 0, 0)
			tx.Program = []core.Instruction{{Op: core.OpPush, Value: a}, {Op: core.OpPush, Value: b}, {Op: op}}
			res, err := vm.Execute(tx)
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		},
	}

	snvmCmd.AddCommand(execCmd)
	rootCmd.AddCommand(snvmCmd)
}
