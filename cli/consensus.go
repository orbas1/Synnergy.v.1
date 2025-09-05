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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("MineBlock")
			diff, err := strconv.ParseUint(args[0], 10, 8)
			if err != nil {
				return err
			}
			sb := core.NewSubBlock([]*core.Transaction{}, "validator")
			b := core.NewBlock([]*core.SubBlock{sb}, "")
			consensus.MineBlock(b, uint8(diff))
			ilog.Info("cli_mine", "nonce", b.Nonce)
			fmt.Println("block mined with nonce", b.Nonce)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("AdjustWeights")
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			s, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			consensus.AdjustWeights(d, s)
			ilog.Info("cli_adjust", "pow", consensus.Weights.PoW, "pos", consensus.Weights.PoS, "poh", consensus.Weights.PoH)
			fmt.Printf("new weights -> PoW: %.2f PoS: %.2f PoH: %.2f\n", consensus.Weights.PoW, consensus.Weights.PoS, consensus.Weights.PoH)
			return nil
		},
	}

	thresholdCmd := &cobra.Command{
		Use:   "threshold [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Calculate switching threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Threshold")
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			s, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			th := consensus.Threshold(d, s)
			ilog.Info("cli_threshold", "value", th)
			fmt.Println("threshold:", th)
			return nil
		},
	}

	transitionCmd := &cobra.Command{
		Use:   "transition [demand] [threat] [stake]",
		Args:  cobra.ExactArgs(3),
		Short: "Compute full transition threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("TransitionThreshold")
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			t, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			s, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return err
			}
			thr := consensus.TransitionThreshold(d, t, s)
			ilog.Info("cli_transition", "value", thr)
			fmt.Println("transition threshold:", thr)
			return nil
		},
	}

	difficultyCmd := &cobra.Command{
		Use:   "difficulty [old] [actual] [expected]",
		Args:  cobra.ExactArgs(3),
		Short: "Adjust mining difficulty",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("DifficultyAdjust")
			old, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			actual, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			expected, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				return err
			}
			nd := consensus.DifficultyAdjust(old, actual, expected)
			ilog.Info("cli_difficulty", "value", nd)
			fmt.Println("new difficulty:", nd)
			return nil
		},
	}

	availabilityCmd := &cobra.Command{
		Use:   "availability [pow] [pos] [poh]",
		Args:  cobra.ExactArgs(3),
		Short: "Set validator availability flags",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SetAvailability")
			pow, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			pos, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			poh, err := strconv.ParseBool(args[2])
			if err != nil {
				return err
			}
			consensus.SetAvailability(pow, pos, poh)
			ilog.Info("cli_availability", "pow", pow, "pos", pos, "poh", poh)
			return nil
		},
	}

	powRewardsCmd := &cobra.Command{
		Use:   "powrewards [enabled]",
		Args:  cobra.ExactArgs(1),
		Short: "Toggle PoW rewards availability",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SetPoWRewards")
			en, err := strconv.ParseBool(args[0])
			if err != nil {
				return err
			}
			consensus.SetPoWRewards(en)
			ilog.Info("cli_pow_rewards", "enabled", en)
			return nil
		},
	}

	consensusCmd.AddCommand(mineCmd, weightsCmd, adjustCmd, thresholdCmd, transitionCmd, difficultyCmd, availabilityCmd, powRewardsCmd)
	rootCmd.AddCommand(consensusCmd)
}
