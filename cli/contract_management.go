package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
	ierr "synnergy/internal/errors"
)

var (
	contractMgr = core.NewContractManager(contractRegistry)
)

func init() {
	contractVM.Start()
	cmd := &cobra.Command{
		Use:   "contract-mgr",
		Short: "Administrative contract management",
	}

	transferCmd := &cobra.Command{
		Use:   "transfer [addr] [newOwner]",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer contract ownership",
		RunE: func(cmd *cobra.Command, args []string) error {
			return contractMgr.Transfer(context.Background(), args[0], args[1])
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Pause a contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			return contractMgr.Pause(context.Background(), args[0])
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Resume a paused contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			return contractMgr.Resume(context.Background(), args[0])
		},
	}

	upgradeCmd := &cobra.Command{
		Use:   "upgrade [addr] [wasmHex] [gasLimit]",
		Args:  cobra.ExactArgs(3),
		Short: "Upgrade contract bytecode",
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			gas, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			return contractMgr.Upgrade(context.Background(), args[0], bytes, gas)
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
			fmt.Printf("owner:%s paused:%v gas:%d\n", c.Owner, c.Paused, c.GasLimit)
		},
	}

	cmd.AddCommand(transferCmd, pauseCmd, resumeCmd, upgradeCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}

func printErr(err error) {
	if e, ok := err.(*ierr.Error); ok {
		fmt.Printf("error (%s): %s\n", e.Code, e.Message)
	} else {
		fmt.Println("error:", err)
	}
}
