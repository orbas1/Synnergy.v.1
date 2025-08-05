package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	wt "synnergy/internal/nodes/watchtower"
	"time"
)

type simpleWatchtower struct {
	id      string
	metrics wt.Metrics
	running bool
	forks   []string
}

func (w *simpleWatchtower) ID() string { return w.id }
func (w *simpleWatchtower) Start(ctx context.Context) error {
	w.running = true
	w.metrics.Timestamp = time.Now()
	return nil
}
func (w *simpleWatchtower) Stop() error { w.running = false; return nil }
func (w *simpleWatchtower) ReportFork(height uint64, hash string) {
	w.forks = append(w.forks, fmt.Sprintf("%d:%s", height, hash))
}
func (w *simpleWatchtower) Metrics() wt.Metrics { return w.metrics }

var watch *simpleWatchtower
=======
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

		Use:   "start <id>",
		Short: "Start watchtower",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			watch = &simpleWatchtower{id: args[0]}
			if err := watch.Start(context.Background()); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("started")
		},
	}
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop watchtower",
		Run: func(cmd *cobra.Command, args []string) {
			if watch == nil {
				fmt.Println("not running")
				return
			}
			watch.Stop()
			fmt.Println("stopped")
		},
	}
	cmd.AddCommand(stopCmd)

	reportCmd := &cobra.Command{
		Use:   "report <height> <hash>",
		Short: "Report a fork",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if watch == nil {
				fmt.Println("not running")
				return
			}
			var h uint64
			fmt.Sscanf(args[0], "%d", &h)
			watch.ReportFork(h, args[1])
			fmt.Println("fork recorded")
		},
	}
	cmd.AddCommand(reportCmd)

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show metrics",
		Run: func(cmd *cobra.Command, args []string) {
			if watch == nil {
				fmt.Println("not running")
				return
			}
			m := watch.Metrics()
			fmt.Printf("CPU:%.2f Mem:%d Peers:%d Height:%d Time:%s\n", m.CPUUsage, m.MemoryUsage, m.PeerCount, m.LastBlockHeight, m.Timestamp.Format(time.RFC3339))
		},
	}
	cmd.AddCommand(metricsCmd)

=======
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
