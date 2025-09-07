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
		Use:     "watchtower-node",
		Aliases: []string{"watchtowernode"},
		Short:   "Manage dedicated watchtower node",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create watchtower node",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("id required")
			}
			watchNode = core.NewWatchtowerNode(id, log.New(os.Stdout, "", 0))
			fmt.Println("watchtower node created")
			return nil
		},
	}
	createCmd.Flags().String("id", "", "node id")
	cmd.AddCommand(createCmd)

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start monitoring",
		RunE: func(cmd *cobra.Command, args []string) error {
			if watchNode == nil {
				return fmt.Errorf("node not initialised")
			}
			if err := watchNode.Start(context.Background()); err != nil {
				return fmt.Errorf("start error: %w", err)
			}
			fmt.Println("watchtower started")
			return nil
		},
	}
	cmd.AddCommand(startCmd)

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop monitoring",
		RunE: func(cmd *cobra.Command, args []string) error {
			if watchNode == nil {
				return fmt.Errorf("node not initialised")
			}
			if err := watchNode.Stop(); err != nil {
				return fmt.Errorf("stop error: %w", err)
			}
			fmt.Println("watchtower stopped")
			return nil
		},
	}
	cmd.AddCommand(stopCmd)

	forkCmd := &cobra.Command{
		Use:   "fork <height> <hash>",
		Args:  cobra.ExactArgs(2),
		Short: "Report fork event",
		RunE: func(cmd *cobra.Command, args []string) error {
			if watchNode == nil {
				return fmt.Errorf("node not initialised")
			}
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height")
			}
			watchNode.ReportFork(h, args[1])
			fmt.Println("fork reported")
			return nil
		},
	}
	cmd.AddCommand(forkCmd)

	metricsCmd := &cobra.Command{
		Use:   "metrics",
		Short: "Show latest system metrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			if watchNode == nil {
				return fmt.Errorf("node not initialised")
			}
			fmt.Printf("%+v\n", watchNode.Metrics())
			return nil
		},
	}
	cmd.AddCommand(metricsCmd)

	rootCmd.AddCommand(cmd)
}
