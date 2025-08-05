package cli

import (
	"encoding/hex"
	"fmt"
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
			val, err := hex.DecodeString(args[1])
			if err != nil {
				val = []byte(args[1])
			}
			kademlia.Store(args[0], val)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get [key]",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a value",
		Run: func(cmd *cobra.Command, args []string) {
			if v, ok := kademlia.FindValue(args[0]); ok {
				fmt.Println(string(v))
			} else {
				fmt.Println("not found")
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "closest [target] [n]",
		Args:  cobra.ExactArgs(2),
		Short: "List n closest keys to target",
		Run: func(cmd *cobra.Command, args []string) {
			n, _ := strconv.Atoi(args[1])
			keys := kademlia.Closest(args[0], n)
			for _, k := range keys {
				fmt.Println(k)
			}
		},
	})

	rootCmd.AddCommand(cmd)
}
