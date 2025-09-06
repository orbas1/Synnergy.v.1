package cli

import (
	"encoding/hex"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var kademlia = core.NewKademlia()

func init() {
	cmd := &cobra.Command{
		Use:   "kademlia",
		Short: "Interact with the Kademlia DHT",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "store [key] [value]",
		Args:  cobra.ExactArgs(2),
		Short: "Store a key/value pair",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("KademliaStore")
			val, err := hex.DecodeString(args[1])
			if err != nil {
				val = []byte(args[1])
			}
			kademlia.Store(args[0], val)
			printOutput(map[string]any{"key": args[0], "stored": true})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a value",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("KademliaGet")
			if v, ok := kademlia.FindValue(args[0]); ok {
				printOutput(map[string]any{"value": string(v)})
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "closest [target] [n]",
		Args:  cobra.ExactArgs(2),
		Short: "List n closest keys to target",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("KademliaClosest")
			n, _ := strconv.Atoi(args[1])
			keys := kademlia.Closest(args[0], n)
			printOutput(keys)
		},
	})

	rootCmd.AddCommand(cmd)
}
