package cli

import (
	"fmt"

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
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(adapter.ID())
		},
	}

	adCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(adCmd)
}
