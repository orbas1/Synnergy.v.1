package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
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
		Run: func(cmd *cobra.Command, args []string) {
			if err := contractMgr.Transfer(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Pause a contract",
		Run: func(cmd *cobra.Command, args []string) {
			if err := contractMgr.Pause(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Resume a paused contract",
		Run: func(cmd *cobra.Command, args []string) {
			if err := contractMgr.Resume(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	upgradeCmd := &cobra.Command{
		Use:   "upgrade [addr] [wasmHex] [gasLimit]",
		Args:  cobra.ExactArgs(3),
		Short: "Upgrade contract bytecode",
		Run: func(cmd *cobra.Command, args []string) {
			bytes, err := hex.DecodeString(args[1])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			gas, _ := strconv.ParseUint(args[2], 10, 64)
			if err := contractMgr.Upgrade(args[0], bytes, gas); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show contract metadata",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := contractMgr.Info(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Printf("owner:%s paused:%v gas:%d\n", c.Owner, c.Paused, c.GasLimit)
		},
	}

	cmd.AddCommand(transferCmd, pauseCmd, resumeCmd, upgradeCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
