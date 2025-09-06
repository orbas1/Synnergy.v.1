package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var adapter = core.NewNodeAdapter(currentNode)

func init() {
	adCmd := &cobra.Command{
		Use:   "node_adapter",
		Short: "Inspect node adapter",
	}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show adapted node ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("NodeAdapterInfo")
			printOutput(map[string]any{"id": adapter.ID()})
			return nil
		},
	}

	adCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(adCmd)
}
