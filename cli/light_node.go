package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
	nodes "synnergy/internal/nodes"
)

var lightNode = core.NewLightNode(nodes.Address("light-1"))

func init() {
	cmd := &cobra.Command{
		Use:   "light",
		Short: "Light node operations",
	}

	addCmd := &cobra.Command{
		Use:   "add-header [hash] [height] [parent]",
		Args:  cobra.ExactArgs(3),
		Short: "Add a block header",
		Run: func(cmd *cobra.Command, args []string) {
			h, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid height:", err)
				return
			}
			header := nodes.BlockHeader{Hash: args[0], Height: h, ParentHash: args[2]}
			lightNode.AddHeader(header)
		},
	}

	latestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Show latest header",
		Run: func(cmd *cobra.Command, args []string) {
			h, ok := lightNode.LatestHeader()
			if !ok {
				fmt.Println("no headers")
				return
			}
			fmt.Printf("%d %s %s\n", h.Height, h.Hash, h.ParentHash)
		},
	}

	listCmd := &cobra.Command{
		Use:   "headers",
		Short: "List all headers",
		Run: func(cmd *cobra.Command, args []string) {
			for _, h := range lightNode.Headers() {
				fmt.Printf("%d %s %s\n", h.Height, h.Hash, h.ParentHash)
			}
		},
	}

	cmd.AddCommand(addCmd, latestCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
