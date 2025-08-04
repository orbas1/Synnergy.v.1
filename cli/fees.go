package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
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
