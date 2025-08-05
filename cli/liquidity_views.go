package cli

import (
	"encoding/json"
	"fmt"

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
			if view, ok := poolRegistry.PoolInfo(args[0]); ok {
				b, _ := json.MarshalIndent(view, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all pools",
		Run: func(cmd *cobra.Command, args []string) {
			views := poolRegistry.PoolViews()
			b, _ := json.MarshalIndent(views, "", "  ")
			fmt.Println(string(b))
		},
	})

	rootCmd.AddCommand(cmd)
}
