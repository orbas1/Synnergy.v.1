package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy"
)

var (
	watchtowerCtx    = context.Background()
	watchtowerLogger = log.New(os.Stdout, "", log.LstdFlags)
	watchtowerNode   = synnergy.NewWatchtowerNode("wt-1", watchtowerLogger)
)

func init() {
	cmd := &cobra.Command{
		Use:   "watchtower",
		Short: "Watchtower node operations",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start monitoring",
		Run: func(cmd *cobra.Command, args []string) {
			if err := watchtowerNode.Start(watchtowerCtx); err != nil {
				fmt.Println("start error:", err)
			}
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop monitoring",
		Run: func(cmd *cobra.Command, args []string) {
			if err := watchtowerNode.Stop(); err != nil {
				fmt.Println("stop error:", err)
			}
		},
	}

	forkCmd := &cobra.Command{
		Use:   "fork [height] [hash]",
		Short: "Report a fork",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid height:", err)
				return
			}
			watchtowerNode.ReportFork(h, args[1])
		},
	}

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show latest metrics",
		Run: func(cmd *cobra.Command, args []string) {
			m := watchtowerNode.Metrics()
			fmt.Printf("cpu=%.2f mem=%d peers=%d height=%d time=%s\n", m.CPUUsage, m.MemoryUsage, m.PeerCount, m.LastBlockHeight, m.Timestamp.Format(time.RFC3339))
		},
	}

	cmd.AddCommand(startCmd, stopCmd, forkCmd, metricsCmd)
	rootCmd.AddCommand(cmd)
}
