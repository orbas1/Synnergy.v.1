package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	coinCmd := &cobra.Command{
		Use:   "coin",
		Short: "Display coin information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s total supply: %d\n", core.CoinName, core.TotalSupply)
		},
	}
	rootCmd.AddCommand(coinCmd)
}
