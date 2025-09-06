package cli

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupSubmit")
			id, err := rollupAgg.SubmitBatch(args)
			if err != nil {
				return err
			}
			printOutput(map[string]string{"id": id})
			return nil
		},
	}

	challengeCmd := &cobra.Command{
		Use:   "challenge [batchID] [txIdx] [proof]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a fraud proof for a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupChallenge")
			idx, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			if err := rollupAgg.ChallengeBatch(args[0], idx, []byte(args[2])); err != nil {
				return err
			}
			printOutput(map[string]string{"status": "challenged"})
			return nil
		},
	}

	finalizeCmd := &cobra.Command{
		Use:   "finalize [batchID] [valid]",
		Args:  cobra.ExactArgs(2),
		Short: "Finalize or revert a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupFinalize")
			valid, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			if err := rollupAgg.FinalizeBatch(args[0], valid); err != nil {
				return err
			}
			status := "reverted"
			if valid {
				status = "finalized"
			}
			printOutput(map[string]string{"status": status})
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [batchID]",
		Args:  cobra.ExactArgs(1),
		Short: "Display batch header and state",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupInfo")
			b, ok := rollupAgg.BatchInfo(args[0])
			if !ok {
				printOutput(map[string]string{"error": "not found"})
				return nil
			}
			printOutput(b)
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List recent batches",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupList")
			bs := rollupAgg.ListBatches()
			printOutput(bs)
			return nil
		},
	}

	txsCmd := &cobra.Command{
		Use:   "txs [batchID]",
		Args:  cobra.ExactArgs(1),
		Short: "List transactions in a batch",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RollupTxs")
			txs, err := rollupAgg.BatchTransactions(args[0])
			if err != nil {
				return err
			}
			printOutput(txs)
			return nil
		},
	}

	cmd.AddCommand(submitCmd, challengeCmd, finalizeCmd, infoCmd, listCmd, txsCmd)
	rootCmd.AddCommand(cmd)
}
