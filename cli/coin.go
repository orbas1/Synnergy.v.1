package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var coinJSON bool

func init() {
	coinCmd := &cobra.Command{Use: "coin", Short: "Synthron coin utilities"}
	coinCmd.PersistentFlags().BoolVar(&coinJSON, "json", false, "output as JSON")

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Display coin parameters",
		RunE: func(cmd *cobra.Command, args []string) error {
			resp := map[string]interface{}{
				"name":               core.CoinName,
				"max_supply":         core.MaxSupply,
				"genesis_allocation": core.GenesisAllocation,
			}
			return coinOutput(resp, fmt.Sprintf("name: %s\nmax supply: %d\ngenesis allocation: %d", core.CoinName, core.MaxSupply, core.GenesisAllocation))
		},
	}

	rewardCmd := &cobra.Command{
		Use:   "reward [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Show block reward at a given height",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height: %w", err)
			}
			reward := core.BlockReward(h)
			return coinOutput(map[string]uint64{"reward": reward}, fmt.Sprintf("reward: %d", reward))
		},
	}

	supplyCmd := &cobra.Command{
		Use:   "supply [height]",
		Args:  cobra.ExactArgs(1),
		Short: "Show circulating and remaining supply",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid height: %w", err)
			}
			circ := core.CirculatingSupply(h)
			rem := core.RemainingSupply(h)
			resp := map[string]uint64{"circulating": circ, "remaining": rem}
			return coinOutput(resp, fmt.Sprintf("circulating: %d remaining: %d", circ, rem))
		},
	}

	priceCmd := &cobra.Command{
		Use:   "price [C] [R] [M] [V] [T] [E]",
		Args:  cobra.ExactArgs(6),
		Short: "Calculate initial price from economic factors",
		RunE: func(cmd *cobra.Command, args []string) error {
			C, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid C: %w", err)
			}
			R, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("invalid R: %w", err)
			}
			M, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return fmt.Errorf("invalid M: %w", err)
			}
			V, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				return fmt.Errorf("invalid V: %w", err)
			}
			T, err := strconv.ParseFloat(args[4], 64)
			if err != nil {
				return fmt.Errorf("invalid T: %w", err)
			}
			E, err := strconv.ParseFloat(args[5], 64)
			if err != nil {
				return fmt.Errorf("invalid E: %w", err)
			}
			price := core.InitialPrice(C, R, M, V, T, E)
			return coinOutput(map[string]float64{"price": price}, fmt.Sprintf("price: %f", price))
		},
	}

	alphaCmd := &cobra.Command{
		Use:   "alpha [volatility] [participation] [economic] [norm]",
		Args:  cobra.ExactArgs(4),
		Short: "Compute alpha factor",
		RunE: func(cmd *cobra.Command, args []string) error {
			V, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid volatility: %w", err)
			}
			P, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("invalid participation: %w", err)
			}
			E, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return fmt.Errorf("invalid economic factor: %w", err)
			}
			N, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				return fmt.Errorf("invalid norm: %w", err)
			}
			alpha := core.AlphaFactor(V, P, E, N)
			return coinOutput(map[string]float64{"alpha": alpha}, fmt.Sprintf("alpha: %f", alpha))
		},
	}

	minstakeCmd := &cobra.Command{
		Use:   "minstake [totalTx] [currentReward] [circSupply] [alpha]",
		Args:  cobra.ExactArgs(4),
		Short: "Calculate minimum stake requirement",
		RunE: func(cmd *cobra.Command, args []string) error {
			tx, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("invalid totalTx: %w", err)
			}
			reward, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("invalid currentReward: %w", err)
			}
			supply, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return fmt.Errorf("invalid circSupply: %w", err)
			}
			alpha, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				return fmt.Errorf("invalid alpha: %w", err)
			}
			ms := core.MinimumStake(tx, reward, supply, alpha)
			return coinOutput(map[string]float64{"minimum_stake": ms}, fmt.Sprintf("minimum stake: %f", ms))
		},
	}

	coinCmd.AddCommand(infoCmd, rewardCmd, supplyCmd, priceCmd, alphaCmd, minstakeCmd)
	rootCmd.AddCommand(coinCmd)
}

func coinOutput(v interface{}, plain string) error {
	if coinJSON {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	} else {
		fmt.Println(plain)
	}
	return nil
}
