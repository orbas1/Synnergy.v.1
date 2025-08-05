package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
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
		Run: func(cmd *cobra.Command, args []string) {
			txs := make([]on.Transaction, 0, len(args))
			for _, a := range args {
				parts := strings.Split(a, ":")
				if len(parts) != 3 {
					fmt.Println("invalid tx format")
					return
				}
				var fee uint64
				var size int
				fmt.Sscanf(parts[1], "%d", &fee)
				fmt.Sscanf(parts[2], "%d", &size)
				txs = append(txs, on.Transaction{Hash: parts[0], Fee: fee, Size: size})
			}
			res := feeOpt.Optimize(txs)
			for _, tx := range res {
				fmt.Printf("%s:%d:%d\n", tx.Hash, tx.Fee, tx.Size)
			}
		},
	}
	cmd.AddCommand(feeCmd)

	rootCmd.AddCommand(cmd)
}
