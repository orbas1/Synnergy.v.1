package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
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
			sb := core.SubBlock{Transactions: []*core.Transaction{}}
			b := core.Block{SubBlocks: []core.SubBlock{sb}}
			diff, _ := strconv.ParseUint(args[0], 10, 64)
			consensus.MineBlock(&b, diff)
			fmt.Println("block mined with nonce", b.Nonce)
		},
	}

	weightsCmd := &cobra.Command{
		Use:   "weights",
		Short: "Show current consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("PoW: %.2f PoS: %.2f PoH: %.2f\n", consensus.Weights.PoW, consensus.Weights.PoS, consensus.Weights.PoH)
		},
	}

	adjustCmd := &cobra.Command{
		Use:   "adjust [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Adjust consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			d, _ := strconv.ParseFloat(args[0], 64)
			s, _ := strconv.ParseFloat(args[1], 64)
			consensus.AdjustWeights(d, s)
			fmt.Printf("new weights -> PoW: %.2f PoS: %.2f PoH: %.2f\n", consensus.Weights.PoW, consensus.Weights.PoS, consensus.Weights.PoH)
		},
	}

	thresholdCmd := &cobra.Command{
		Use:   "threshold [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Calculate switching threshold",
		Run: func(cmd *cobra.Command, args []string) {
			d, _ := strconv.ParseFloat(args[0], 64)
			s, _ := strconv.ParseFloat(args[1], 64)
			fmt.Println("threshold:", consensus.Threshold(d, s))
		},
	}

	transitionCmd := &cobra.Command{
		Use:   "transition [demand] [threat] [stake]",
		Args:  cobra.ExactArgs(3),
		Short: "Compute full transition threshold",
		Run: func(cmd *cobra.Command, args []string) {
			d, _ := strconv.ParseFloat(args[0], 64)
			t, _ := strconv.ParseFloat(args[1], 64)
			s, _ := strconv.ParseFloat(args[2], 64)
			fmt.Println("transition threshold:", consensus.TransitionThreshold(d, t, s))
		},
	}

	difficultyCmd := &cobra.Command{
		Use:   "difficulty [old] [actual] [expected]",
		Args:  cobra.ExactArgs(3),
		Short: "Adjust mining difficulty",
		Run: func(cmd *cobra.Command, args []string) {
			old, _ := strconv.ParseFloat(args[0], 64)
			actual, _ := strconv.ParseFloat(args[1], 64)
			expected, _ := strconv.ParseFloat(args[2], 64)
			fmt.Println("new difficulty:", consensus.DifficultyAdjust(old, actual, expected))
		},
	}

	availabilityCmd := &cobra.Command{
		Use:   "availability [pow] [pos] [poh]",
		Args:  cobra.ExactArgs(3),
		Short: "Set validator availability flags",
		Run: func(cmd *cobra.Command, args []string) {
			pow := args[0] == "true"
			pos := args[1] == "true"
			poh := args[2] == "true"
			consensus.SetAvailability(pow, pos, poh)
		},
	}

	powRewardsCmd := &cobra.Command{
		Use:   "powrewards [enabled]",
		Args:  cobra.ExactArgs(1),
		Short: "Toggle PoW rewards availability",
		Run: func(cmd *cobra.Command, args []string) {
			en := args[0] == "true"
			consensus.SetPoWRewards(en)
		},
	}

	consensusCmd.AddCommand(mineCmd, weightsCmd, adjustCmd, thresholdCmd, transitionCmd, difficultyCmd, availabilityCmd, powRewardsCmd)
	rootCmd.AddCommand(consensusCmd)
}
