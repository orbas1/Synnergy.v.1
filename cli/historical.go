package cli

import (
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
				printOutput("invalid height: " + err.Error())
				return
			}
			summary := nodes.BlockSummary{Height: h, Hash: args[1], Timestamp: time.Now()}
			if err := historicalNode.ArchiveBlock(summary); err != nil {
				printOutput("archive error: " + err.Error())
				return
			}
			gasPrint("HistoricalArchive")
			printOutput(map[string]any{"status": "archived", "height": h, "hash": args[1]})
		},
	}

	heightCmd := &cobra.Command{
		Use:   "height [n]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch block by height",
		Run: func(cmd *cobra.Command, args []string) {
			h, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				printOutput("invalid height: " + err.Error())
				return
			}
			b, ok := historicalNode.GetBlockByHeight(h)
			if !ok {
				printOutput("not found")
				return
			}
			gasPrint("HistoricalHeight")
			printOutput(map[string]any{"height": b.Height, "hash": b.Hash, "timestamp": b.Timestamp.Format(time.RFC3339)})
		},
	}

	hashCmd := &cobra.Command{
		Use:   "hash [hash]",
		Args:  cobra.ExactArgs(1),
		Short: "Fetch block by hash",
		Run: func(cmd *cobra.Command, args []string) {
			b, ok := historicalNode.GetBlockByHash(args[0])
			if !ok {
				printOutput("not found")
				return
			}
			gasPrint("HistoricalHash")
			printOutput(map[string]any{"height": b.Height, "hash": b.Hash, "timestamp": b.Timestamp.Format(time.RFC3339)})
		},
	}

	totalCmd := &cobra.Command{
		Use:   "total",
		Short: "Show total archived blocks",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("HistoricalTotal")
			printOutput(map[string]int{"total": historicalNode.TotalBlocks()})
		},
	}

	histCmd.AddCommand(archiveCmd, heightCmd, hashCmd, totalCmd)
	rootCmd.AddCommand(histCmd)
}
