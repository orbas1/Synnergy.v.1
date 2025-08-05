package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	optnodes "synnergy/internal/nodes/extra/optimization_nodes"
)

var optimizer = &optnodes.SimpleOptimizer{}

func init() {
	cmd := &cobra.Command{
		Use:   "optimization",
		Short: "Optimization node utilities",
	}

	suggestCmd := &cobra.Command{
		Use:   "suggest [cpu] [mem] [latency] [throughput]",
		Args:  cobra.ExactArgs(4),
		Short: "Get scaling suggestion",
		Run: func(cmd *cobra.Command, args []string) {
			cpu, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				fmt.Println("invalid cpu usage:", err)
				return
			}
			mem, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				fmt.Println("invalid memory usage:", err)
				return
			}
			lat, err := strconv.ParseFloat(args[2], 64)
			if err != nil {
				fmt.Println("invalid latency:", err)
				return
			}
			th, err := strconv.ParseFloat(args[3], 64)
			if err != nil {
				fmt.Println("invalid throughput:", err)
				return
			}
			s := optimizer.Optimize(optnodes.Metrics{CPUUsage: cpu, MemoryUsage: mem, LatencyMS: lat, Throughput: th})
			fmt.Printf("scale=%v notes=%s\n", s.ScaleResources, s.Notes)
		},
	}

	cmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(cmd)
}
