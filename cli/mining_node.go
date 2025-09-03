package cli

import (
	"encoding/hex"
	"fmt"

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
			miningNode.Start()
			printOutput("started")
			return nil
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop mining",
		RunE: func(cmd *cobra.Command, args []string) error {
			miningNode.Stop()
			printOutput("stopped")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show mining status",
		RunE: func(cmd *cobra.Command, args []string) error {
			printOutput(miningNode.IsMining())
			return nil
		},
	}

	mineBlockCmd := &cobra.Command{
		Use:   "mine [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Perform one mining attempt",
		RunE: func(cmd *cobra.Command, args []string) error {
			hash, err := miningNode.Mine([]byte(args[0]))
			if err != nil {
				return err
			}
			printOutput(hash)
			return nil
		},
	}

	rateCmd := &cobra.Command{
		Use:   "hashrate",
		Short: "Display hash rate hint",
		RunE: func(cmd *cobra.Command, args []string) error {
			printOutput(miningNode.HashRateHint())
			return nil
		},
	}

	// helper to show hex encoded data
	hexCmd := &cobra.Command{
		Use:   "hex [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Mine pre-hex-encoded input",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("invalid hex")
			}
			hash, err := miningNode.Mine(b)
			if err != nil {
				return err
			}
			printOutput(hash)
			return nil
		},
	}

	mineCmd.AddCommand(startCmd, stopCmd, statusCmd, mineBlockCmd, rateCmd, hexCmd)
	rootCmd.AddCommand(mineCmd)
}
