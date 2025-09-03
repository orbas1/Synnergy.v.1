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
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height: %w", err)
			}
			header := nodes.BlockHeader{Hash: args[0], Height: h, ParentHash: args[2]}
			lightNode.AddHeader(header)
			printOutput(header)
			return nil
		},
	}

	latestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Show latest header",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, ok := lightNode.LatestHeader()
			if !ok {
				return fmt.Errorf("no headers")
			}
			printOutput(h)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "headers",
		Short: "List all headers",
		RunE: func(cmd *cobra.Command, args []string) error {
			printOutput(lightNode.Headers())
			return nil
		},
	}

	cmd.AddCommand(addCmd, latestCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
