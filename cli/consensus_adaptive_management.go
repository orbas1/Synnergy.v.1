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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Adjust")
			d, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return err
			}
			s, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return err
			}
			w := adaptiveManager.Adjust(d, s)
			fmt.Printf("PoW: %.2f PoS: %.2f PoH: %.2f\n", w.PoW, w.PoS, w.PoH)
			return nil
		},
	}

	thresholdCmd := &cobra.Command{
		Use:   "threshold [demand] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Compute switching threshold",
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
			fmt.Println(adaptiveManager.Threshold(d, s))
			return nil
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
