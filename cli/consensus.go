package cli

import (
	"fmt"
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
			gasPrint("MineBlock")
			sb := core.NewSubBlock([]*core.Transaction{}, "validator")
			b := core.NewBlock([]*core.SubBlock{sb}, "")
			diff, _ := strconv.ParseUint(args[0], 10, 8)
			consensus.MineBlock(b, uint8(diff))
			ilog.Info("cli_mine", "nonce", b.Nonce)
			fmt.Println("block mined with nonce", b.Nonce)
		},
	}

	weightsCmd := &cobra.Command{
		Use:   "weights",
		Short: "Show current consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Weights")
			ilog.Info("cli_weights", "pow", consensus.Weights.PoW, "pos", consensus.Weights.PoS, "poh", consensus.Weights.PoH)
			fmt.Printf("PoW: %.2f PoS: %.2f PoH: %.2f\n", consensus.Weights.PoW, consensus.Weights.PoS, consensus.Weights.PoH)
		},
	}

	adjustCmd := &cobra.Command{
		Use:   "adjust [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Adjust consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("AdjustWeights")
			d, _ := strconv.ParseFloat(args[0], 64)
			s, _ := strconv.ParseFloat(args[1], 64)
			consensus.AdjustWeights(d, s)
			ilog.Info("cli_adjust", "pow", consensus.Weights.PoW, "pos", consensus.Weights.PoS, "poh", consensus.Weights.PoH)
			fmt.Printf("new weights -> PoW: %.2f PoS: %.2f PoH: %.2f\n", consensus.Weights.PoW, consensus.Weights.PoS, consensus.Weights.PoH)
		},
	}

	thresholdCmd := &cobra.Command{
		Use:   "threshold [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Calculate switching threshold",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Threshold")
			d, _ := strconv.ParseFloat(args[0], 64)
			s, _ := strconv.ParseFloat(args[1], 64)
			th := consensus.Threshold(d, s)
			ilog.Info("cli_threshold", "value", th)
			fmt.Println("threshold:", th)
		},
	}

	transitionCmd := &cobra.Command{
		Use:   "transition [demand] [threat] [stake]",
		Args:  cobra.ExactArgs(3),
		Short: "Compute full transition threshold",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("TransitionThreshold")
			d, _ := strconv.ParseFloat(args[0], 64)
			t, _ := strconv.ParseFloat(args[1], 64)
			s, _ := strconv.ParseFloat(args[2], 64)
			thr := consensus.TransitionThreshold(d, t, s)
			ilog.Info("cli_transition", "value", thr)
			fmt.Println("transition threshold:", thr)
		},
	}

	difficultyCmd := &cobra.Command{
		Use:   "difficulty [old] [actual] [expected]",
		Args:  cobra.ExactArgs(3),
		Short: "Adjust mining difficulty",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DifficultyAdjust")
			old, _ := strconv.ParseFloat(args[0], 64)
			actual, _ := strconv.ParseFloat(args[1], 64)
			expected, _ := strconv.ParseFloat(args[2], 64)
			nd := consensus.DifficultyAdjust(old, actual, expected)
			ilog.Info("cli_difficulty", "value", nd)
			fmt.Println("new difficulty:", nd)
		},
	}

	availabilityCmd := &cobra.Command{
		Use:   "availability [pow] [pos] [poh]",
		Args:  cobra.ExactArgs(3),
		Short: "Set validator availability flags",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SetAvailability")
			pow := args[0] == "true"
			pos := args[1] == "true"
			poh := args[2] == "true"
			consensus.SetAvailability(pow, pos, poh)
			ilog.Info("cli_availability", "pow", pow, "pos", pos, "poh", poh)
		},
	}

	powRewardsCmd := &cobra.Command{
		Use:   "powrewards [enabled]",
		Args:  cobra.ExactArgs(1),
		Short: "Toggle PoW rewards availability",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SetPoWRewards")
			en := args[0] == "true"
			consensus.SetPoWRewards(en)
			ilog.Info("cli_pow_rewards", "enabled", en)
		},
	}

	consensusCmd.AddCommand(mineCmd, weightsCmd, adjustCmd, thresholdCmd, transitionCmd, difficultyCmd, availabilityCmd, powRewardsCmd)
	rootCmd.AddCommand(consensusCmd)
}
