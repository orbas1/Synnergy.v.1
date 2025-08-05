package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var consensusService = core.NewConsensusService(currentNode)

func init() {
	cmd := &cobra.Command{
		Use:   "consensus-service",
		Short: "Run consensus mining service",
	}

	startCmd := &cobra.Command{
		Use:   "start [ms]",
		Args:  cobra.ExactArgs(1),
		Short: "Start mining loop with interval milliseconds",
		Run: func(cmd *cobra.Command, args []string) {
			ms, _ := time.ParseDuration(args[0] + "ms")
			consensusService.Start(ms)
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the mining loop",
		Run: func(cmd *cobra.Command, args []string) {
			consensusService.Stop()
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show service info",
		Run: func(cmd *cobra.Command, args []string) {
			h, running := consensusService.Info()
			fmt.Printf("height: %d running: %v\n", h, running)
		},
	}

	cmd.AddCommand(startCmd, stopCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
