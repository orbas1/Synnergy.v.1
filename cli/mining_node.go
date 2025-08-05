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
		Run: func(cmd *cobra.Command, args []string) {
			miningNode.Start()
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop mining",
		Run: func(cmd *cobra.Command, args []string) {
			miningNode.Stop()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show mining status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(miningNode.IsMining())
		},
	}

	mineBlockCmd := &cobra.Command{
		Use:   "mine [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Perform one mining attempt",
		Run: func(cmd *cobra.Command, args []string) {
			hash, err := miningNode.Mine([]byte(args[0]))
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(hash)
		},
	}

	rateCmd := &cobra.Command{
		Use:   "hashrate",
		Short: "Display hash rate hint",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(miningNode.HashRateHint())
		},
	}

	// helper to show hex encoded data
	hexCmd := &cobra.Command{
		Use:   "hex [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Mine pre-hex-encoded input",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := hex.DecodeString(args[0])
			if err != nil {
				fmt.Println("invalid hex")
				return
			}
			hash, err := miningNode.Mine(b)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(hash)
		},
	}

	mineCmd.AddCommand(startCmd, stopCmd, statusCmd, mineBlockCmd, rateCmd, hexCmd)
	rootCmd.AddCommand(mineCmd)
}
