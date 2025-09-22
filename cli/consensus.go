package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
	ilog "synnergy/internal/log"
)

var consensus = core.NewSynnergyConsensus()

func init() {
	consensusCmd := &cobra.Command{
		Use:   "consensus",
		Short: "Consensus operations",
	}
	mineCmd := &cobra.Command{
		Use:   "mine [difficulty]",
		Args:  cobra.ExactArgs(1),
		Short: "Mine a block",
		Run: func(cmd *cobra.Command, args []string) {
			diff, err := strconv.ParseUint(args[0], 10, 8)
			if err != nil {
				printOutput(map[string]any{"error": "invalid difficulty"})
				return
			}
			sb := core.NewSubBlock([]*core.Transaction{}, "validator")
			b := core.NewBlock([]*core.SubBlock{sb}, "")
			if err := consensus.MineBlock(context.Background(), b, uint8(diff)); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			ilog.Info("cli_mine", "nonce", b.Nonce)
			gasPrint("MineBlock")
			printOutput(map[string]any{"nonce": b.Nonce})
		},
	}

	weightsCmd := &cobra.Command{
		Use:   "weights",
		Short: "Show current consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Weights")
			ilog.Info("cli_weights", "pow", consensus.Weights.PoW, "pos", consensus.Weights.PoS, "poh", consensus.Weights.PoH)
			printOutput(map[string]float64{"pow": consensus.Weights.PoW, "pos": consensus.Weights.PoS, "poh": consensus.Weights.PoH})
		},
	}

	adjustCmd := &cobra.Command{
		Use:   "adjust [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Adjust consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid demand"})
				return
			}
			s, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid stake"})
				return
			}
			consensus.AdjustWeights(d, s)
			ilog.Info("cli_adjust", "pow", consensus.Weights.PoW, "pos", consensus.Weights.PoS, "poh", consensus.Weights.PoH)
			gasPrint("AdjustWeights")
			printOutput(map[string]float64{"pow": consensus.Weights.PoW, "pos": consensus.Weights.PoS, "poh": consensus.Weights.PoH})
		},
	}

	thresholdCmd := &cobra.Command{
		Use:   "threshold [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Calculate switching threshold",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid demand"})
				return
			}
			s, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid stake"})
				return
			}
			th := consensus.Threshold(d, s)
			ilog.Info("cli_threshold", "value", th)
			gasPrint("Threshold")
			printOutput(map[string]float64{"threshold": th})
		},
	}

	transitionCmd := &cobra.Command{
		Use:   "transition [demand] [threat] [stake]",
		Args:  cobra.ExactArgs(3),
		Short: "Compute full transition threshold",
		Run: func(cmd *cobra.Command, args []string) {
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid demand"})
				return
			}
			t, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid threat"})
				return
			}
			s, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid stake"})
				return
			}
			thr := consensus.TransitionThreshold(d, t, s)
			ilog.Info("cli_transition", "value", thr)
			gasPrint("TransitionThreshold")
			printOutput(map[string]float64{"threshold": thr})
		},
	}

	difficultyCmd := &cobra.Command{
		Use:   "difficulty [old] [actual] [expected]",
		Args:  cobra.ExactArgs(3),
		Short: "Adjust mining difficulty",
		Run: func(cmd *cobra.Command, args []string) {
			old, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid old"})
				return
			}
			actual, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid actual"})
				return
			}
			expected, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid expected"})
				return
			}
			nd := consensus.DifficultyAdjust(old, actual, expected)
			ilog.Info("cli_difficulty", "value", nd)
			gasPrint("DifficultyAdjust")
			printOutput(map[string]float64{"difficulty": nd})
		},
	}

	availabilityCmd := &cobra.Command{
		Use:   "availability [pow] [pos] [poh]",
		Args:  cobra.ExactArgs(3),
		Short: "Set validator availability flags",
		Run: func(cmd *cobra.Command, args []string) {
			pow, err := strconv.ParseBool(args[0])
			if err != nil {
				printOutput(map[string]any{"error": "invalid pow"})
				return
			}
			pos, err := strconv.ParseBool(args[1])
			if err != nil {
				printOutput(map[string]any{"error": "invalid pos"})
				return
			}
			poh, err := strconv.ParseBool(args[2])
			if err != nil {
				printOutput(map[string]any{"error": "invalid poh"})
				return
			}
			consensus.SetAvailability(pow, pos, poh)
			ilog.Info("cli_availability", "pow", pow, "pos", pos, "poh", poh)
			gasPrint("SetAvailability")
			printOutput(map[string]bool{"pow": pow, "pos": pos, "poh": poh})
		},
	}

	powRewardsCmd := &cobra.Command{
		Use:   "powrewards [enabled]",
		Args:  cobra.ExactArgs(1),
		Short: "Toggle PoW rewards availability",
		Run: func(cmd *cobra.Command, args []string) {
			en, err := strconv.ParseBool(args[0])
			if err != nil {
				printOutput(map[string]any{"error": "invalid flag"})
				return
			}
			consensus.SetPoWRewards(en)
			ilog.Info("cli_pow_rewards", "enabled", en)
			gasPrint("SetPoWRewards")
			printOutput(map[string]bool{"enabled": en})
		},
	}

	consensusCmd.AddCommand(mineCmd, weightsCmd, adjustCmd, thresholdCmd, transitionCmd, difficultyCmd, availabilityCmd, powRewardsCmd)
	rootCmd.AddCommand(consensusCmd)
}
