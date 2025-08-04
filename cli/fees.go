package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	feesCmd := &cobra.Command{
		Use:   "fees",
		Short: "Estimate transaction fees and provide feedback",
	}

	estimateCmd := &cobra.Command{
		Use:   "estimate",
		Short: "Estimate fees for a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			txTypeStr, _ := cmd.Flags().GetString("type")
			units, _ := cmd.Flags().GetUint64("units")
			tip, _ := cmd.Flags().GetUint64("tip")
			load, _ := cmd.Flags().GetFloat64("load")

			var txType core.TransactionType
			switch txTypeStr {
			case "transfer":
				txType = core.TxTypeTransfer
			case "purchase":
				txType = core.TxTypePurchase
			case "token":
				txType = core.TxTypeTokenInteraction
			case "contract":
				txType = core.TxTypeContract
			case "wallet":
				txType = core.TxTypeWalletVerification
			default:
				fmt.Printf("unknown type %s\n", txTypeStr)
				return
			}

			// Placeholder for recent fee data which would be synced from the network.
			recent := []uint64{1, 2, 1, 2, 1}
			base := core.CalculateBaseFee(recent, load)
			base, variable := core.AdjustFeeRates(base, 1, load)
			fb := core.EstimateFee(txType, units, base, variable, tip)
			fmt.Printf("Base: %d\nVariable: %d\nPriority: %d\nTotal: %d\n", fb.Base, fb.Variable, fb.Priority, fb.Total)
		},
	}
	estimateCmd.Flags().String("type", "transfer", "transaction type (transfer,purchase,token,contract,wallet)")
	estimateCmd.Flags().Uint64("units", 1, "units for fee calculation")
	estimateCmd.Flags().Uint64("tip", 0, "priority fee")
	estimateCmd.Flags().Float64("load", 0, "network load factor")
	feesCmd.AddCommand(estimateCmd)

	feedbackCmd := &cobra.Command{
		Use:   "feedback",
		Short: "Submit feedback about fee estimates",
		Run: func(cmd *cobra.Command, args []string) {
			message, _ := cmd.Flags().GetString("message")
			fmt.Printf("feedback received: %s\n", message)
		},
	}
	feedbackCmd.Flags().String("message", "", "feedback message")
	feesCmd.AddCommand(feedbackCmd)

	rootCmd.AddCommand(feesCmd)
	feeCmd := &cobra.Command{
		Use:   "fees",
		Short: "Fee utilities",
	}

	shareCmd := &cobra.Command{
		Use:   "share [total] [validatorWeight] [minerWeight]",
		Args:  cobra.ExactArgs(3),
		Short: "Compute proportional validator and miner fee shares",
		Run: func(cmd *cobra.Command, args []string) {
			total, _ := strconv.ParseUint(args[0], 10, 64)
			v, _ := strconv.ParseUint(args[1], 10, 64)
			m, _ := strconv.ParseUint(args[2], 10, 64)
			shares := core.ShareProportional(total, map[string]uint64{"validator": v, "miner": m})
			fmt.Printf("shares: %v\n", shares)
		},
	}

	feeCmd.AddCommand(shareCmd)
	rootCmd.AddCommand(feeCmd)
}
