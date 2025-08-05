package cli

import (
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
	rollupsCmd := &cobra.Command{
		Use:   "rollups",
		Short: "Manage rollup batches",
	}

	submitCmd := &cobra.Command{
		Use:   "submit [tx1 tx2 ...]",
		Args:  cobra.MinimumNArgs(1),
		Short: "Submit a new rollup batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := rollupAgg.SubmitBatch(args)
			if err != nil {
				return err
			}
			fmt.Println(id)
			return nil
		},
	}

	challengeCmd := &cobra.Command{
		Use:   "challenge <batchID> <txIdx> <proof>",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a fraud proof for a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			idx, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			return rollupAgg.ChallengeBatch(args[0], idx, []byte(args[2]))
		},
	}

	finalizeCmd := &cobra.Command{
		Use:   "finalize <batchID> <valid>",
		Args:  cobra.ExactArgs(2),
		Short: "Finalize or revert a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			valid, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			return rollupAgg.FinalizeBatch(args[0], valid)
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <batchID>",
		Args:  cobra.ExactArgs(1),
		Short: "Display batch header and state",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, ok := rollupAgg.BatchInfo(args[0])
			if !ok {
				return fmt.Errorf("batch not found")
			}
			fmt.Printf("%s %s %d\n", b.ID, b.Status, len(b.Transactions))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List recent batches",
		Run: func(cmd *cobra.Command, args []string) {
			for _, b := range rollupAgg.ListBatches() {
				fmt.Printf("%s %s %d\n", b.ID, b.Status, len(b.Transactions))
			}
		},
	}

	txsCmd := &cobra.Command{
		Use:   "txs <batchID>",
		Args:  cobra.ExactArgs(1),
		Short: "List transactions in a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			txs, err := rollupAgg.BatchTransactions(args[0])
			if err != nil {
				return err
			}
			for _, tx := range txs {
				fmt.Println(tx)
			}
			return nil
		},
	}

	rollupsCmd.AddCommand(submitCmd, challengeCmd, finalizeCmd, infoCmd, listCmd, txsCmd)
	rootCmd.AddCommand(rollupsCmd)
}
