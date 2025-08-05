package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
	"synnergy/internal/nodes"
)

var baseNode = core.NewBaseNode(nodes.Address("base1"))

func init() {
	bnCmd := &cobra.Command{
		Use:   "basenode",
		Short: "Manage base node lifecycle and peers",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the node",
		RunE:  func(cmd *cobra.Command, args []string) error { return baseNode.Start() },
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the node",
		RunE:  func(cmd *cobra.Command, args []string) error { return baseNode.Stop() },
	}

	runningCmd := &cobra.Command{
		Use:   "running",
		Short: "Check if the node is running",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println(baseNode.IsRunning()) },
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List known peers",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range baseNode.Peers() {
				fmt.Println(p)
			}
		},
	}

	dialCmd := &cobra.Command{
		Use:   "dial [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Dial a seed peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			return baseNode.DialSeed(nodes.Address(args[0]))
		},
	}

	bnCmd.AddCommand(startCmd, stopCmd, runningCmd, peersCmd, dialCmd)
	rootCmd.AddCommand(bnCmd)
}
