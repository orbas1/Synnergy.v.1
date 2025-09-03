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
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := watchtowerNode.Start(watchtowerCtx); err != nil {
				return err
			}
			printOutput("started")
			return nil
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop monitoring",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := watchtowerNode.Stop(); err != nil {
				return err
			}
			printOutput("stopped")
			return nil
		},
	}

	forkCmd := &cobra.Command{
		Use:   "fork [height] [hash]",
		Short: "Report a fork",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height: %w", err)
			}
			watchtowerNode.ReportFork(h, args[1])
			printOutput("reported")
			return nil
		},
	}

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show latest metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			m := watchtowerNode.Metrics()
			printOutput(map[string]any{
				"cpu":    fmt.Sprintf("%.2f", m.CPUUsage),
				"mem":    m.MemoryUsage,
				"peers":  m.PeerCount,
				"height": m.LastBlockHeight,
				"time":   m.Timestamp.Format(time.RFC3339),
			})
			return nil
		},
	}

	cmd.AddCommand(startCmd, stopCmd, forkCmd, metricsCmd)
	rootCmd.AddCommand(cmd)
}
