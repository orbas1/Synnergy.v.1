//go:build experimental

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	nodes "synnergy/internal/nodes"
)

var expNode *nodes.ExperimentalNode

func init() {
	cmd := &cobra.Command{
		Use:   "experimental",
		Short: "Manage experimental node (requires -tags=experimental)",
	}

	createCmd := &cobra.Command{
		Use:   "create [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a new experimental node",
		Run: func(cmd *cobra.Command, args []string) {
			expNode = nodes.NewExperimentalNode(nodes.Address(args[0]))
			fmt.Println("created")
		},
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if expNode == nil {
				return fmt.Errorf("node not created")
			}
			return expNode.Start()
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the node",
		RunE: func(cmd *cobra.Command, args []string) error {
			if expNode == nil {
				return fmt.Errorf("node not created")
			}
			return expNode.Stop()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show node status",
		Run: func(cmd *cobra.Command, args []string) {
			if expNode == nil {
				fmt.Println("no node")
				return
			}
			if expNode.IsRunning() {
				fmt.Println("running")
			} else {
				fmt.Println("stopped")
			}
		},
	}

	dialCmd := &cobra.Command{
		Use:   "dial [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Dial a seed peer",
		RunE: func(cmd *cobra.Command, args []string) error {
			if expNode == nil {
				return fmt.Errorf("node not created")
			}
			return expNode.DialSeed(nodes.Address(args[0]))
		},
	}

	cmd.AddCommand(createCmd, startCmd, stopCmd, statusCmd, dialCmd)
	rootCmd.AddCommand(cmd)
}
