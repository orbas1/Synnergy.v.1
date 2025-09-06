package cli

import (
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
			gasPrint("LightAddHeader")
			h, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid height"})
				return
			}
			header := nodes.BlockHeader{Hash: args[0], Height: h, ParentHash: args[2]}
			lightNode.AddHeader(header)
			printOutput(header)
		},
	}

	latestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Show latest header",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LightLatest")
			h, ok := lightNode.LatestHeader()
			if !ok {
				printOutput(map[string]any{"error": "no headers"})
				return
			}
			printOutput(h)
		},
	}

	listCmd := &cobra.Command{
		Use:   "headers",
		Short: "List all headers",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LightHeaders")
			printOutput(lightNode.Headers())
		},
	}

	cmd.AddCommand(addCmd, latestCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
