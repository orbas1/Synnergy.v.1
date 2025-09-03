package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

// adaptiveManager tunes consensus weights using network metrics.
var adaptiveManager = core.NewAdaptiveManager(consensus)

func init() {
	cmd := &cobra.Command{
		Use:   "consensus-adaptive",
		Short: "Adaptive consensus weight management",
	}

	adjustCmd := &cobra.Command{
		Use:   "adjust [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Adjust weights and show result",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Adjust")
			d, _ := strconv.ParseFloat(args[0], 64)
			s, _ := strconv.ParseFloat(args[1], 64)
			w := adaptiveManager.Adjust(d, s)
			fmt.Printf("PoW: %.2f PoS: %.2f PoH: %.2f\n", w.PoW, w.PoS, w.PoH)
		},
	}

	thresholdCmd := &cobra.Command{
		Use:   "threshold [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Compute switching threshold",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Threshold")
			d, _ := strconv.ParseFloat(args[0], 64)
			s, _ := strconv.ParseFloat(args[1], 64)
			fmt.Println(adaptiveManager.Threshold(d, s))
		},
	}

	weightsCmd := &cobra.Command{
		Use:   "weights",
		Short: "Show current consensus weights",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Weights")
			w := adaptiveManager.Weights()
			fmt.Printf("PoW: %.2f PoS: %.2f PoH: %.2f\n", w.PoW, w.PoS, w.PoH)
		},
	}

	cmd.AddCommand(adjustCmd, thresholdCmd, weightsCmd)
	rootCmd.AddCommand(cmd)
}
