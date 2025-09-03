package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var peerMgr = core.NewPeerManager()

func init() {
	cmd := &cobra.Command{
		Use:   "peer",
		Short: "Peer discovery and connections",
	}

	discoverCmd := &cobra.Command{
		Use:   "discover [topic]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List known peers or those advertising a topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				for _, id := range peerMgr.Discover(args[0]) {
					fmt.Fprintln(cmd.OutOrStdout(), id)
				}
				return nil
			}
			for _, id := range peerMgr.ListPeers() {
				fmt.Fprintln(cmd.OutOrStdout(), id)
			}
			return nil
		},
	}

	connectCmd := &cobra.Command{
		Use:   "connect [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to a peer by address",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := peerMgr.Connect(args[0])
			fmt.Fprintln(cmd.OutOrStdout(), "connected", id)
			return nil
		},
	}

	advertiseCmd := &cobra.Command{
		Use:   "advertise [topic]",
		Args:  cobra.ExactArgs(1),
		Short: "Advertise current node on a topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			peerMgr.Advertise(currentNode.ID, args[0])
			return nil
		},
	}

	cmd.AddCommand(discoverCmd, connectCmd, advertiseCmd)
	rootCmd.AddCommand(cmd)
}
