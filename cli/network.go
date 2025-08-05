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
		Run:   func(cmd *cobra.Command, args []string) { network.Start() },
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop network services",
		Run:   func(cmd *cobra.Command, args []string) { network.Stop() },
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List peers",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range network.Peers() {
				fmt.Println(p)
			}
		},
	}

	broadcastCmd := &cobra.Command{
		Use:   "broadcast [topic] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Publish data on the network",
		Run: func(cmd *cobra.Command, args []string) {
			network.Publish(args[0], []byte(args[1]))
		},
	}

	subscribeCmd := &cobra.Command{
		Use:   "subscribe [topic]",
		Args:  cobra.ExactArgs(1),
		Short: "Subscribe and print messages for a topic",
		Run: func(cmd *cobra.Command, args []string) {
			ch := network.Subscribe(args[0])
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			for {
				select {
				case msg := <-ch:
					fmt.Println(string(msg))
				case <-ctx.Done():
					return
				}
			}
		},
	}

	netCmd.AddCommand(startCmd, stopCmd, peersCmd, broadcastCmd, subscribeCmd)
	rootCmd.AddCommand(netCmd)
}
