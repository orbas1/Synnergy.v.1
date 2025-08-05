package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var watchNode *core.Watchtower

func init() {
	cmd := &cobra.Command{
		Use:   "watchtower",
		Short: "Manage watchtower node",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create watchtower node",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			watchNode = core.NewWatchtowerNode(id, log.New(os.Stdout, "", 0))
			fmt.Println("watchtower node created")
		},
	}
	createCmd.Flags().String("id", "", "node id")
	cmd.AddCommand(createCmd)

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start monitoring",
		Run: func(cmd *cobra.Command, args []string) {
			if watchNode == nil {
				fmt.Println("node not initialised")
				return
			}
			if err := watchNode.Start(context.Background()); err != nil {
				fmt.Println("start error:", err)
				return
			}
			fmt.Println("watchtower started")
		},
	}
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop monitoring",
		Run: func(cmd *cobra.Command, args []string) {
			if watchNode == nil {
				fmt.Println("node not initialised")
				return
			}
			if err := watchNode.Stop(); err != nil {
				fmt.Println("stop error:", err)
				return
			}
			fmt.Println("watchtower stopped")
		},
	}
	cmd.AddCommand(stopCmd)

	forkCmd := &cobra.Command{
		Use:   "fork <height> <hash>",
		Args:  cobra.ExactArgs(2),
		Short: "Report fork event",
		Run: func(cmd *cobra.Command, args []string) {
			if watchNode == nil {
				fmt.Println("node not initialised")
				return
			}
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid height")
				return
			}
			watchNode.ReportFork(h, args[1])
			fmt.Println("fork reported")
		},
	}
	cmd.AddCommand(forkCmd)

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show latest system metrics",
		Run: func(cmd *cobra.Command, args []string) {
			if watchNode == nil {
				fmt.Println("node not initialised")
				return
			}
			fmt.Printf("%+v\n", watchNode.Metrics())
		},
	}
	cmd.AddCommand(metricsCmd)

	rootCmd.AddCommand(cmd)
}
