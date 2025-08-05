package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	rollupAgg = core.NewRollupAggregator()
	rollupMgr = core.NewRollupManager(rollupAgg)
)

func init() {
	cmd := &cobra.Command{
		Use:   "rollups",
		Short: "Manage rollup batches",
	}

	submitCmd := &cobra.Command{
		Use:   "submit [tx ...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Submit a new rollup batch",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := rollupAgg.SubmitBatch(args)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(id)
		},
	}

	challengeCmd := &cobra.Command{
		Use:   "challenge [batchID] [txIdx] [proof]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a fraud proof for a batch",
		Run: func(cmd *cobra.Command, args []string) {
			idx, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println("invalid index")
				return
			}
			if err := rollupAgg.ChallengeBatch(args[0], idx, []byte(args[2])); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	finalizeCmd := &cobra.Command{
		Use:   "finalize [batchID] [valid]",
		Args:  cobra.ExactArgs(2),
		Short: "Finalize or revert a batch",
		Run: func(cmd *cobra.Command, args []string) {
			valid, err := strconv.ParseBool(args[1])
			if err != nil {
				fmt.Println("invalid bool")
				return
			}
			if err := rollupAgg.FinalizeBatch(args[0], valid); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [batchID]",
		Args:  cobra.ExactArgs(1),
		Short: "Display batch header and state",
		Run: func(cmd *cobra.Command, args []string) {
			b, ok := rollupAgg.BatchInfo(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			out, _ := json.MarshalIndent(b, "", "  ")
			fmt.Println(string(out))
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List recent batches",
		Run: func(cmd *cobra.Command, args []string) {
			bs := rollupAgg.ListBatches()
			out, _ := json.MarshalIndent(bs, "", "  ")
			fmt.Println(string(out))
		},
	}

	txsCmd := &cobra.Command{
		Use:   "txs [batchID]",
		Args:  cobra.ExactArgs(1),
		Short: "List transactions in a batch",
		Run: func(cmd *cobra.Command, args []string) {
			txs, err := rollupAgg.BatchTransactions(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			for _, tx := range txs {
				fmt.Println(tx)
			}
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause the rollup aggregator",
		Run: func(cmd *cobra.Command, args []string) {
			rollupMgr.Pause()
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume the rollup aggregator",
		Run: func(cmd *cobra.Command, args []string) {
			rollupMgr.Resume()
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show aggregator status",
		Run: func(cmd *cobra.Command, args []string) {
			if rollupMgr.Status() {
				fmt.Println("paused")
			} else {
				fmt.Println("running")
			}
		},
	}

	cmd.AddCommand(submitCmd, challengeCmd, finalizeCmd, infoCmd, listCmd, txsCmd, pauseCmd, resumeCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
