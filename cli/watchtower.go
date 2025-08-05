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
	wtCtx    = context.Background()
	wtLogger = log.New(os.Stdout, "", log.LstdFlags)
	wtNode   = synnergy.NewWatchtowerNode("wt-1", wtLogger)
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
			if err := wtNode.Start(wtCtx); err != nil {
				fmt.Println("start error:", err)
			}
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop monitoring",
		Run: func(cmd *cobra.Command, args []string) {
			if err := wtNode.Stop(); err != nil {
				fmt.Println("stop error:", err)
			}
		},
	}

	forkCmd := &cobra.Command{
		Use:   "fork [height] [hash]",
		Args:  cobra.ExactArgs(2),
		Short: "Report a fork",
		Run: func(cmd *cobra.Command, args []string) {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid height:", err)
				return
			}
			wtNode.ReportFork(h, args[1])
		},
	}

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show latest metrics",
		Run: func(cmd *cobra.Command, args []string) {
			m := wtNode.Metrics()
			fmt.Printf("cpu=%.2f mem=%d peers=%d height=%d time=%s\n", m.CPUUsage, m.MemoryUsage, m.PeerCount, m.LastBlockHeight, m.Timestamp.Format(time.RFC3339))
		},
	}

	cmd.AddCommand(startCmd, stopCmd, forkCmd, metricsCmd)
	rootCmd.AddCommand(cmd)
}
