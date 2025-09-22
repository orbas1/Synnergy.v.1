package cli

import (
	"context"
	"encoding/hex"
	"strconv"

	"github.com/spf13/cobra"
	ierr "synnergy/internal/errors"
)

func init() {
	ensureContractComponents()
	_ = contractVM.Start()
	cmd := &cobra.Command{
		Use:   "contract-mgr",
		Short: "Administrative contract management",
	}

	transferCmd := &cobra.Command{
		Use:   "transfer [addr] [newOwner]",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer contract ownership",
		Run: func(cmd *cobra.Command, args []string) {
			if err := contractMgr.Transfer(context.Background(), args[0], args[1]); err != nil {
				printErr(err)
				return
			}
			gasPrint("TransferContract")
			printOutput(map[string]any{"status": "transferred", "address": args[0], "owner": args[1]})
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Pause a contract",
		Run: func(cmd *cobra.Command, args []string) {
			if err := contractMgr.Pause(context.Background(), args[0]); err != nil {
				printErr(err)
				return
			}
			gasPrint("PauseContract")
			printOutput(map[string]any{"status": "paused", "address": args[0]})
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Resume a paused contract",
		Run: func(cmd *cobra.Command, args []string) {
			if err := contractMgr.Resume(context.Background(), args[0]); err != nil {
				printErr(err)
				return
			}
			gasPrint("ResumeContract")
			printOutput(map[string]any{"status": "resumed", "address": args[0]})
		},
	}

	upgradeCmd := &cobra.Command{
		Use:   "upgrade [addr] [wasmHex] [gasLimit]",
		Args:  cobra.ExactArgs(3),
		Short: "Upgrade contract bytecode",
		Run: func(cmd *cobra.Command, args []string) {
			bytes, err := hex.DecodeString(args[1])
			if err != nil {
				printOutput(map[string]any{"error": "invalid wasm"})
				return
			}
			gas, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid gas"})
				return
			}
			if err := contractMgr.Upgrade(context.Background(), args[0], bytes, gas); err != nil {
				printErr(err)
				return
			}
			gasPrint("UpgradeContract")
			printOutput(map[string]any{"status": "upgraded", "address": args[0], "gas": gas})
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show contract metadata",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := contractMgr.Info(context.Background(), args[0])
			if err != nil {
				printErr(err)
				return
			}
			gasPrint("ContractInfo")
			printOutput(map[string]any{"owner": c.Owner, "paused": c.Paused, "gas": c.GasLimit})
		},
	}

	cmd.AddCommand(transferCmd, pauseCmd, resumeCmd, upgradeCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}

func printErr(err error) {
	if e, ok := err.(*ierr.Error); ok {
		printOutput(map[string]any{"error": e.Message, "code": e.Code})
	} else {
		printOutput(map[string]any{"error": err.Error()})
	}
}
