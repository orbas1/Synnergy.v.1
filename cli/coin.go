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

	supplyCmd := &cobra.Command{
		Use:   "supply [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Show circulating and remaining supply",
		Run: func(cmd *cobra.Command, args []string) {
			h, _ := strconv.ParseUint(args[0], 10, 64)
			circ := core.CirculatingSupply(h)
			rem := core.RemainingSupply(h)
			fmt.Printf("circulating: %d remaining: %d\n", circ, rem)
		},
	}

	priceCmd := &cobra.Command{
		Use:   "price [C] [R] [M] [V] [T] [E]",
		Args:  cobra.ExactArgs(6),
		Short: "Calculate initial price from economic factors",
		Run: func(cmd *cobra.Command, args []string) {
			C, _ := strconv.ParseFloat(args[0], 64)
			R, _ := strconv.ParseFloat(args[1], 64)
			M, _ := strconv.ParseFloat(args[2], 64)
			V, _ := strconv.ParseFloat(args[3], 64)
			T, _ := strconv.ParseFloat(args[4], 64)
			E, _ := strconv.ParseFloat(args[5], 64)
			fmt.Println("price:", core.InitialPrice(C, R, M, V, T, E))
		},
	}

	alphaCmd := &cobra.Command{
		Use:   "alpha [volatility] [participation] [economic] [norm]",
		Args:  cobra.ExactArgs(4),
		Short: "Compute alpha factor",
		Run: func(cmd *cobra.Command, args []string) {
			V, _ := strconv.ParseFloat(args[0], 64)
			P, _ := strconv.ParseFloat(args[1], 64)
			E, _ := strconv.ParseFloat(args[2], 64)
			N, _ := strconv.ParseFloat(args[3], 64)
			fmt.Println("alpha:", core.AlphaFactor(V, P, E, N))
		},
	}

	minstakeCmd := &cobra.Command{
		Use:   "minstake [totalTx] [currentReward] [circSupply] [alpha]",
		Args:  cobra.ExactArgs(4),
		Short: "Calculate minimum stake requirement",
		Run: func(cmd *cobra.Command, args []string) {
			tx, _ := strconv.ParseFloat(args[0], 64)
			reward, _ := strconv.ParseFloat(args[1], 64)
			supply, _ := strconv.ParseFloat(args[2], 64)
			alpha, _ := strconv.ParseFloat(args[3], 64)
			fmt.Println("minimum stake:", core.MinimumStake(tx, reward, supply, alpha))
		},
	}

	coinCmd.AddCommand(infoCmd, rewardCmd, supplyCmd, priceCmd, alphaCmd, minstakeCmd)
	rootCmd.AddCommand(coinCmd)
}
