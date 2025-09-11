package cli

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var miningNode = core.NewMiningNode(1000)

func init() {
	mineCmd := &cobra.Command{
		Use:   "mining",
		Short: "Control a mining node",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start mining",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MiningStart")
			miningNode.Start()
			printOutput("started")
			return nil
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop mining",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MiningStop")
			miningNode.Stop()
			printOutput("stopped")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show mining status",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MiningStatus")
			printOutput(map[string]bool{"mining": miningNode.IsMining()})
			return nil
		},
	}

	mineBlockCmd := &cobra.Command{
		Use:   "mine [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Perform one mining attempt",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MiningMine")
			hash, err := miningNode.Mine([]byte(args[0]))
			if err != nil {
				return err
			}
			printOutput(map[string]string{"hash": hash})
			return nil
		},
	}

	rateCmd := &cobra.Command{
		Use:   "hashrate",
		Short: "Display hash rate hint",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MiningHashrate")
			printOutput(map[string]uint64{"hashrate": miningNode.HashRateHint()})
			return nil
		},
	}

	var timeout int
	mineUntilCmd := &cobra.Command{
		Use:   "mine-until [data] [prefix]",
		Args:  cobra.ExactArgs(2),
		Short: "Mine until hash has prefix or timeout elapses",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MineUntil")
			ctx := context.Background()
			if timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
				defer cancel()
			}
			hash, nonce, err := miningNode.MineUntil(ctx, []byte(args[0]), args[1])
			if err != nil {
				return err
			}
			printOutput(map[string]interface{}{"hash": hash, "nonce": nonce})
			return nil
		},
	}
	mineUntilCmd.Flags().IntVar(&timeout, "timeout", 0, "timeout in seconds")

	// helper to show hex encoded data
	hexCmd := &cobra.Command{
		Use:   "hex [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Mine pre-hex-encoded input",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MiningHex")
			b, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("invalid hex")
			}
			hash, err := miningNode.Mine(b)
			if err != nil {
				return err
			}
			printOutput(map[string]string{"hash": hash})
			return nil
		},
	}

	mineCmd.AddCommand(startCmd, stopCmd, statusCmd, mineBlockCmd, mineUntilCmd, rateCmd, hexCmd)
	rootCmd.AddCommand(mineCmd)
}
