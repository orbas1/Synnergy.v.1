package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
	nodes "synnergy/internal/nodes"
)

var historicalNode = core.NewHistoricalNode()

func init() {
	histCmd := &cobra.Command{
		Use:   "historical",
		Short: "Historical node operations",
	}

	archiveCmd := &cobra.Command{
		Use:   "archive [height] [hash]",
		Args:  cobra.ExactArgs(2),
		Short: "Archive a block summary",
		Run: func(cmd *cobra.Command, args []string) {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid height:", err)
				return
			}
			summary := nodes.BlockSummary{Height: h, Hash: args[1], Timestamp: time.Now()}
			if err := historicalNode.ArchiveBlock(summary); err != nil {
				fmt.Println("archive error:", err)
			}
		},
	}

	heightCmd := &cobra.Command{
		Use:   "height [n]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch block by height",
		Run: func(cmd *cobra.Command, args []string) {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid height:", err)
				return
			}
			b, ok := historicalNode.GetBlockByHeight(h)
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%d %s %s\n", b.Height, b.Hash, b.Timestamp.Format(time.RFC3339))
		},
	}

	hashCmd := &cobra.Command{
		Use:   "hash [hash]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch block by hash",
		Run: func(cmd *cobra.Command, args []string) {
			b, ok := historicalNode.GetBlockByHash(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%d %s %s\n", b.Height, b.Hash, b.Timestamp.Format(time.RFC3339))
		},
	}

	totalCmd := &cobra.Command{
		Use:   "total",
		Short: "Show total archived blocks",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(historicalNode.TotalBlocks())
		},
	}

	histCmd.AddCommand(archiveCmd, heightCmd, hashCmd, totalCmd)
	rootCmd.AddCommand(histCmd)
}
