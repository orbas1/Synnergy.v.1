package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy"
	nodes "synnergy/internal/nodes"
)

var holoNode = nodes.NewHolographicNode(nodes.Address("holo-1"))

func init() {
	cmd := &cobra.Command{
		Use:   "holographic",
		Short: "Holographic node operations",
	}

	storeCmd := &cobra.Command{
		Use:   "store [id] [data] [shards]",
		Args:  cobra.ExactArgs(3),
		Short: "Store data as holographic frame",
		Run: func(cmd *cobra.Command, args []string) {
			n, err := strconv.Atoi(args[2])
			if err != nil {
				printOutput("invalid shard count: " + err.Error())
				return
			}
			frame := synnergy.SplitHolographic(args[0], []byte(args[1]), n)
			holoNode.Store(frame)
			gasPrint("HolographicStore")
			printOutput(map[string]any{"status": "stored", "id": args[0], "shards": n})
		},
	}

	retrieveCmd := &cobra.Command{
		Use:   "retrieve [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve stored data",
		Run: func(cmd *cobra.Command, args []string) {
			frame, ok := holoNode.Retrieve(args[0])
			if !ok {
				printOutput("not found")
				return
			}
			data := synnergy.ReconstructHolographic(frame)
			gasPrint("HolographicRetrieve")
			printOutput(map[string]any{"id": args[0], "data": string(data)})
		},
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List known peers",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("HolographicPeers")
			printOutput(map[string]any{"peers": holoNode.Peers()})
		},
	}

	dialCmd := &cobra.Command{
		Use:   "dial [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Dial a seed peer",
		Run: func(cmd *cobra.Command, args []string) {
			if err := holoNode.DialSeed(nodes.Address(args[0])); err != nil {
				printOutput("dial error: " + err.Error())
				return
			}
			gasPrint("HolographicDial")
			printOutput(map[string]any{"status": "connected", "addr": args[0]})
		},
	}

	cmd.AddCommand(storeCmd, retrieveCmd, peersCmd, dialCmd)
	rootCmd.AddCommand(cmd)
}
