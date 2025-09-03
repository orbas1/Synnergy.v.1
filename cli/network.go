package cli

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var network = core.NewNetwork(biometricSvc)

func init() {
	netCmd := &cobra.Command{
		Use:   "network",
		Short: "Control networking stack",
	}

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start network services",
		RunE: func(cmd *cobra.Command, args []string) error {
			network.Start()
			fmt.Fprintln(cmd.OutOrStdout(), "network started")
			return nil
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop network services",
		RunE: func(cmd *cobra.Command, args []string) error {
			network.Stop()
			fmt.Fprintln(cmd.OutOrStdout(), "network stopped")
			return nil
		},
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List peers",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, p := range network.Peers() {
				fmt.Fprintln(cmd.OutOrStdout(), p)
			}
			return nil
		},
	}

	broadcastCmd := &cobra.Command{
		Use:   "broadcast [topic] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Publish data on the network",
		RunE: func(cmd *cobra.Command, args []string) error {
			network.Publish(args[0], []byte(args[1]))
			return nil
		},
	}

	subscribeCmd := &cobra.Command{
		Use:   "subscribe [topic]",
		Args:  cobra.ExactArgs(1),
		Short: "Subscribe and print messages for a topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			ch := network.Subscribe(args[0])
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			for {
				select {
				case msg := <-ch:
					fmt.Fprintln(cmd.OutOrStdout(), string(msg))
				case <-ctx.Done():
					return nil
				}
			}
		},
	}

	netCmd.AddCommand(startCmd, stopCmd, peersCmd, broadcastCmd, subscribeCmd)
	rootCmd.AddCommand(netCmd)
}
