package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"synnergy/core"
	"synnergy/treasury"
)

var (
	coinJSON     bool
	treasuryOnce sync.Once
	treasuryInst *treasury.SynthronTreasury
	treasuryErr  error
)

func init() {
	coinCmd := &cobra.Command{Use: "coin", Short: "Synthron coin utilities"}
	coinCmd.PersistentFlags().BoolVar(&coinJSON, "json", false, "output as JSON")

	telemetryCmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Report treasury diagnostics and optionally orchestrate actions",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			operator, _ := cmd.Flags().GetString("operator")
			if operator != "" {
				ctx = treasury.WithOperator(ctx, operator)
			}

			treas, err := getSynthronTreasury(ctx)
			if err != nil {
				return err
			}

			authorize, _ := cmd.Flags().GetStringSlice("authorize-operator")
			for _, addr := range authorize {
				addr = strings.TrimSpace(addr)
				if addr == "" {
					continue
				}
				if err := treas.AuthorizeOperator(ctx, addr); err != nil {
					return err
				}
			}

			revoke, _ := cmd.Flags().GetStringSlice("revoke-operator")
			for _, addr := range revoke {
				addr = strings.TrimSpace(addr)
				if addr == "" {
					continue
				}
				if err := treas.RevokeOperator(ctx, addr); err != nil {
					return err
				}
			}

			issueArg, _ := cmd.Flags().GetString("issue")
			if issueArg != "" {
				addr, amount, err := parseAddressAmount(issueArg)
				if err != nil {
					return err
				}
				if _, err := treas.Issue(ctx, addr, amount); err != nil {
					return err
				}
			}

			burnArg, _ := cmd.Flags().GetString("burn")
			if burnArg != "" {
				addr, amount, err := parseAddressAmount(burnArg)
				if err != nil {
					return err
				}
				if err := treas.Burn(ctx, addr, amount); err != nil {
					return err
				}
			}

			transferArg, _ := cmd.Flags().GetString("transfer")
			if transferArg != "" {
				addr, amount, err := parseAddressAmount(transferArg)
				if err != nil {
					return err
				}
				if err := treas.Transfer(ctx, nil, addr, amount, 0); err != nil {
					return err
				}
			}

			bridgeArg, _ := cmd.Flags().GetString("bridge")
			if bridgeArg != "" {
				parts := strings.SplitN(bridgeArg, ":", 2)
				if len(parts) != 2 {
					return fmt.Errorf("invalid bridge format, expected source:target")
				}
				if _, err := treas.RegisterConsensusLink(ctx, parts[0], parts[1]); err != nil {
					return err
				}
			}

			diag := treas.Diagnostics(ctx)
			if coinJSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(diag)
			}
			fmt.Fprint(cmd.OutOrStdout(), treasury.SynthronTreasurySummary(diag))
			fmt.Fprintf(cmd.OutOrStdout(), "component health: vm=%s ledger=%s wallet=%s consensus=%s authorities=%s\n", diag.Health.VM, diag.Health.Ledger, diag.Health.Wallet, diag.Health.Consensus, diag.Health.Authorities)
			if len(diag.InsertedOpcodes) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "inserted opcodes: %v\n", diag.InsertedOpcodes)
			}
			return nil
		},
	}
	telemetryCmd.Flags().String("issue", "", "issue coins using addr:amount syntax")
	telemetryCmd.Flags().String("burn", "", "burn coins using addr:amount syntax")
	telemetryCmd.Flags().String("transfer", "", "transfer coins using addr:amount syntax")
	telemetryCmd.Flags().String("bridge", "", "register a consensus bridge using source:target syntax")
	telemetryCmd.Flags().String("operator", "", "execute privileged actions as the specified operator")
	telemetryCmd.Flags().StringSlice("authorize-operator", nil, "grant treasury operator access (repeatable)")
	telemetryCmd.Flags().StringSlice("revoke-operator", nil, "revoke treasury operator access (repeatable)")

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

	coinCmd.AddCommand(infoCmd, rewardCmd, supplyCmd, priceCmd, alphaCmd, minstakeCmd, telemetryCmd)
	rootCmd.AddCommand(coinCmd)
}

func getSynthronTreasury(ctx context.Context) (*treasury.SynthronTreasury, error) {
	treasuryOnce.Do(func() {
		treasuryInst, treasuryErr = treasury.DefaultSynthronTreasury(ctx)
	})
	return treasuryInst, treasuryErr
}

func parseAddressAmount(input string) (string, uint64, error) {
	parts := strings.SplitN(input, ":", 2)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid format, expected addr:amount")
	}
	amount, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid amount: %w", err)
	}
	return parts[0], amount, nil
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
