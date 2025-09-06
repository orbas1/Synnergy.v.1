package cli

import (
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "liquidity_views",
		Short: "Inspect liquidity pool views",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "info [poolID]",
		Args:  cobra.ExactArgs(1),
		Short: "Show pool state",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityViewsInfo")
			if view, ok := poolRegistry.PoolInfo(args[0]); ok {
				printOutput(view)
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all pools",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityViewsList")
			views := poolRegistry.PoolViews()
			printOutput(views)
		},
	})

	rootCmd.AddCommand(cmd)
}
