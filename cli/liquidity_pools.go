package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var poolRegistry = core.NewLiquidityPoolRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "liquidity_pools",
		Short: "Manage liquidity pools",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create [tokenA] [tokenB] [feeBps]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Create a new liquidity pool",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityPoolCreate")
			fee := uint16(30)
			if len(args) == 3 {
				f, err := strconv.Atoi(args[2])
				if err != nil {
					printOutput(map[string]any{"error": "invalid fee"})
					return
				}
				fee = uint16(f)
			}
			id := fmt.Sprintf("%s-%s", args[0], args[1])
			if _, err := poolRegistry.Create(id, args[0], args[1], fee); err != nil {
				printOutput(map[string]any{"error": err.Error()})
			} else {
				printOutput(map[string]any{"status": "created", "id": id})
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "add [poolID] [provider] [amtA] [amtB]",
		Args:  cobra.ExactArgs(4),
		Short: "Add liquidity to a pool",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityPoolAdd")
			p, ok := poolRegistry.Get(args[0])
			if !ok {
				printOutput(map[string]any{"error": "not found"})
				return
			}
			a, _ := strconv.ParseUint(args[2], 10, 64)
			b, _ := strconv.ParseUint(args[3], 10, 64)
			lp, err := p.AddLiquidity(args[1], a, b)
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(lp)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "swap [poolID] [tokenIn] [amtIn] [minOut]",
		Args:  cobra.ExactArgs(4),
		Short: "Swap tokens within a pool",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityPoolSwap")
			p, ok := poolRegistry.Get(args[0])
			if !ok {
				printOutput(map[string]any{"error": "not found"})
				return
			}
			amtIn, _ := strconv.ParseUint(args[2], 10, 64)
			minOut, _ := strconv.ParseUint(args[3], 10, 64)
			out, err := p.Swap(args[1], amtIn, minOut)
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(out)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "remove [poolID] [provider] [lpTokens]",
		Args:  cobra.ExactArgs(3),
		Short: "Remove liquidity from a pool",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityPoolRemove")
			p, ok := poolRegistry.Get(args[0])
			if !ok {
				printOutput(map[string]any{"error": "not found"})
				return
			}
			lp, _ := strconv.ParseUint(args[2], 10, 64)
			a, b, err := p.RemoveLiquidity(args[1], lp)
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]uint64{"amountA": a, "amountB": b})
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "info [poolID]",
		Args:  cobra.ExactArgs(1),
		Short: "Show pool state",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityPoolInfo")
			if view, ok := poolRegistry.PoolInfo(args[0]); ok {
				printOutput(view)
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all pools",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LiquidityPoolList")
			views := poolRegistry.PoolViews()
			printOutput(views)
		},
	})

	rootCmd.AddCommand(cmd)
}
