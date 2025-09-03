package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	on "synnergy/internal/nodes/optimization_nodes"
)

var feeOpt on.FeeOptimizer

func init() {
	cmd := &cobra.Command{
		Use:   "optimize",
		Short: "Transaction optimisation utilities",
	}

	feeCmd := &cobra.Command{
		Use:   "fee <tx...>",
		Short: "Order transactions by fee density",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txs := make([]on.Transaction, 0, len(args))
			for _, a := range args {
				parts := strings.Split(a, ":")
				if len(parts) != 3 {
					return fmt.Errorf("invalid tx format")
				}
				var fee uint64
				var size int
				fmt.Sscanf(parts[1], "%d", &fee)
				fmt.Sscanf(parts[2], "%d", &size)
				txs = append(txs, on.Transaction{Hash: parts[0], Fee: fee, Size: size})
			}
			res := feeOpt.Optimize(txs)
			printOutput(res)
			return nil
		},
	}
	cmd.AddCommand(feeCmd)

	rootCmd.AddCommand(cmd)
}
