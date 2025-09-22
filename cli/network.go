package cli

import (
	"context"
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
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("NetworkStart")
			network.Start()
			printOutput("network started")
		},
	}

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop network services",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("NetworkStop")
			network.Stop()
			printOutput("network stopped")
		},
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List peers",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("NetworkPeers")
			printOutput(map[string][]string{"peers": network.Peers()})
		},
	}

	broadcastCmd := &cobra.Command{
		Use:   "broadcast [topic] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Publish data on the network",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("NetworkBroadcast")
			network.Publish(args[0], []byte(args[1]))
			printOutput("message published")
		},
	}

	subscribeCmd := &cobra.Command{
		Use:   "subscribe [topic]",
		Args:  cobra.ExactArgs(1),
		Short: "Subscribe and print messages for a topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NetworkSubscribe")
			ch, cancel := network.Subscribe(args[0])
			defer cancel()
			ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
			defer stop()
			for {
				select {
				case msg := <-ch:
					printOutput(string(msg))
				case <-ctx.Done():
					return nil
				}
			}
		},
	}

	netCmd.AddCommand(startCmd, stopCmd, peersCmd, broadcastCmd, subscribeCmd)
	rootCmd.AddCommand(netCmd)
}
