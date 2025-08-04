package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	coinCmd := &cobra.Command{Use: "coin", Short: "Synthron coin utilities"}

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display coin parameters",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("name: %s\nmax supply: %d\ngenesis allocation: %d\n", core.CoinName, core.MaxSupply, core.GenesisAllocation)
		},
	}

	rewardCmd := &cobra.Command{
		Use:   "reward [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Show block reward at a given height",
		Run: func(cmd *cobra.Command, args []string) {
			h, _ := strconv.ParseUint(args[0], 10, 64)
			fmt.Println("reward:", core.BlockReward(h))
		},
	}

	coinCmd.AddCommand(infoCmd, rewardCmd)
	rootCmd.AddCommand(coinCmd)
}
