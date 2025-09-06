package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var mobileMiner = core.NewMobileMiningNode(500, 100)

func init() {
	mmCmd := &cobra.Command{
		Use:   "mobile_mining",
		Short: "Operate a mobile mining node",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start mobile mining",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MobileMiningStart")
			mobileMiner.Start()
			printOutput("started")
			return nil
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop mobile mining",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MobileMiningStop")
			mobileMiner.Stop()
			printOutput("stopped")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show mining status",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MobileMiningStatus")
			printOutput(map[string]bool{"mining": mobileMiner.IsMining()})
			return nil
		},
	}

	mineCmd := &cobra.Command{
		Use:   "mine [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Mine once with provided data",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MobileMiningMine")
			hash, err := mobileMiner.Mine([]byte(args[0]))
			if err != nil {
				return err
			}
			printOutput(map[string]string{"hash": hash})
			return nil
		},
	}

	setPower := &cobra.Command{
		Use:   "set-power [limit]",
		Args:  cobra.ExactArgs(1),
		Short: "Set power limit",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MobileMiningSetPower")
			limit, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid limit")
			}
			mobileMiner.SetPowerLimit(limit)
			printOutput(map[string]uint64{"limit": limit})
			return nil
		},
	}

	powerCmd := &cobra.Command{
		Use:   "power",
		Short: "Show power limit",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MobileMiningPower")
			printOutput(map[string]uint64{"limit": mobileMiner.PowerLimit()})
			return nil
		},
	}

	mmCmd.AddCommand(startCmd, stopCmd, statusCmd, mineCmd, setPower, powerCmd)
	rootCmd.AddCommand(mmCmd)
}
