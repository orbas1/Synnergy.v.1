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
		Run:   func(cmd *cobra.Command, args []string) { mobileMiner.Start() },
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop mobile mining",
		Run:   func(cmd *cobra.Command, args []string) { mobileMiner.Stop() },
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show mining status",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println(mobileMiner.IsMining()) },
	}

	mineCmd := &cobra.Command{
		Use:   "mine [data]",
		Args:  cobra.ExactArgs(1),
		Short: "Mine once with provided data",
		Run: func(cmd *cobra.Command, args []string) {
			hash, err := mobileMiner.Mine([]byte(args[0]))
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(hash)
		},
	}

	setPower := &cobra.Command{
		Use:   "set-power [limit]",
		Args:  cobra.ExactArgs(1),
		Short: "Set power limit",
		Run: func(cmd *cobra.Command, args []string) {
			limit, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid limit")
				return
			}
			mobileMiner.SetPowerLimit(limit)
		},
	}

	powerCmd := &cobra.Command{
		Use:   "power",
		Short: "Show power limit",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println(mobileMiner.PowerLimit()) },
	}

	mmCmd.AddCommand(startCmd, stopCmd, statusCmd, mineCmd, setPower, powerCmd)
	rootCmd.AddCommand(mmCmd)
}
