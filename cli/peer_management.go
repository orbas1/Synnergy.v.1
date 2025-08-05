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
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				for _, id := range peerMgr.Discover(args[0]) {
					fmt.Println(id)
				}
				return
			}
			for _, id := range peerMgr.ListPeers() {
				fmt.Println(id)
			}
		},
	}

	connectCmd := &cobra.Command{
		Use:   "connect [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to a peer by address",
		Run: func(cmd *cobra.Command, args []string) {
			id := peerMgr.Connect(args[0])
			fmt.Println("connected", id)
		},
	}

	advertiseCmd := &cobra.Command{
		Use:   "advertise [topic]",
		Args:  cobra.ExactArgs(1),
		Short: "Advertise current node on a topic",
		Run: func(cmd *cobra.Command, args []string) {
			peerMgr.Advertise(currentNode.ID, args[0])
		},
	}

	cmd.AddCommand(discoverCmd, connectCmd, advertiseCmd)
	rootCmd.AddCommand(cmd)
}
