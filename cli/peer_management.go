package cli

import (
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
			gasPrint("PeerDiscover")
			var peers []string
			if len(args) == 1 {
				peers = peerMgr.Discover(args[0])
			} else {
				peers = peerMgr.ListPeers()
			}
			printOutput(peers)
			return nil
		},
	}

	connectCmd := &cobra.Command{
		Use:   "connect [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Connect to a peer by address",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("PeerConnect")
			id := peerMgr.Connect(args[0])
			printOutput(map[string]string{"status": "connected", "id": id})
			return nil
		},
	}

	advertiseCmd := &cobra.Command{
		Use:   "advertise [topic]",
		Args:  cobra.ExactArgs(1),
		Short: "Advertise current node on a topic",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("PeerAdvertise")
			peerMgr.Advertise(currentNode.ID, args[0])
			printOutput(map[string]string{"status": "advertised", "topic": args[0]})
			return nil
		},
	}

	countCmd := &cobra.Command{
		Use:   "count",
		Short: "Show number of known peers",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("PeerCount")
			printOutput(map[string]int{"count": peerMgr.Count()})
			return nil
		},
	}

	cmd.AddCommand(discoverCmd, connectCmd, advertiseCmd, countCmd)
	rootCmd.AddCommand(cmd)
}
