package cli

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var simpleVM *core.SimpleVM

func init() {
	cmd := &cobra.Command{
		Use:   "simplevm",
		Short: "Manage the simple virtual machine",
	}

	createCmd := &cobra.Command{
		Use:   "create [mode]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Create a new VM instance",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("VMCreate")
			mode := core.VMLight
			if len(args) == 1 {
				switch args[0] {
				case "heavy":
					mode = core.VMHeavy
				case "superlight":
					mode = core.VMSuperLight
				}
			}
			simpleVM = core.NewSimpleVM(mode)
			printOutput(map[string]string{"status": "created"})
		},
	}
	cmd.AddCommand(createCmd)

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the VM",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("VMStart")
			if simpleVM == nil {
				printOutput(map[string]string{"error": "vm not created"})
				return
			}
			if err := simpleVM.Start(); err != nil {
				printOutput(map[string]string{"error": err.Error()})
				return
			}
			printOutput(map[string]string{"status": "started"})
		},
	}
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the VM",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("VMStop")
			if simpleVM == nil {
				printOutput(map[string]string{"error": "vm not created"})
				return
			}
			if err := simpleVM.Stop(); err != nil {
				printOutput(map[string]string{"error": err.Error()})
				return
			}
			printOutput(map[string]string{"status": "stopped"})
		},
	}
	cmd.AddCommand(stopCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show running status",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("VMStatus")
			if simpleVM == nil {
				printOutput(map[string]string{"error": "vm not created"})
				return
			}
			printOutput(map[string]bool{"running": simpleVM.Status()})
		},
	}
	cmd.AddCommand(statusCmd)

	var timeoutMS int
	execCmd := &cobra.Command{
		Use:   "exec <wasmHex> [argsHex] [gas]",
		Args:  cobra.RangeArgs(1, 3),
		Short: "Execute bytecode on the VM",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("VMExec")
			if simpleVM == nil {
				printOutput(map[string]string{"error": "vm not created"})
				return
			}
			if !simpleVM.Status() {
				printOutput(map[string]string{"error": "vm not running"})
				return
			}
			wasm, err := hex.DecodeString(args[0])
			if err != nil {
				printOutput(map[string]string{"error": "invalid wasm"})
				return
			}
			var in []byte
			if len(args) > 1 {
				in, err = hex.DecodeString(args[1])
				if err != nil {
					printOutput(map[string]string{"error": "invalid args"})
					return
				}
			}
			gas := uint64(100)
			if len(args) > 2 {
				gas, err = strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					printOutput(map[string]string{"error": "invalid gas"})
					return
				}
			}
			ctx := context.Background()
			if timeoutMS > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeoutMS)*time.Millisecond)
				defer cancel()
			}
			out, used, err := simpleVM.ExecuteContext(ctx, wasm, "", in, gas)
			if err != nil {
				printOutput(map[string]string{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"out": hex.EncodeToString(out), "gasUsed": used})
		},
	}
	execCmd.Flags().IntVar(&timeoutMS, "timeout", 0, "execution timeout in ms (0 for none)")
	cmd.AddCommand(execCmd)

	rootCmd.AddCommand(cmd)
}
