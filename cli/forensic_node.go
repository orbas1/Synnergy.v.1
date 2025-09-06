package cli

import (
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
	nodes "synnergy/internal/nodes"
)

var forensic = core.NewForensicNode()

func init() {
	cmd := &cobra.Command{
		Use:   "forensic",
		Short: "Record transactions and network traces",
	}

	txCmd := &cobra.Command{
		Use:   "record-tx",
		Short: "Record a minimal transaction",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ForensicRecordTx")
			hash, _ := cmd.Flags().GetString("hash")
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			value, _ := cmd.Flags().GetUint64("value")
			tx := nodes.TransactionLite{Hash: hash, From: from, To: to, Value: value, Timestamp: time.Now()}
			if err := forensic.RecordTransaction(tx); err != nil {
				printOutput(err.Error())
			} else {
				printOutput("transaction recorded")
			}
		},
	}
	txCmd.Flags().String("hash", "", "transaction hash")
	txCmd.Flags().String("from", "", "from address")
	txCmd.Flags().String("to", "", "to address")
	txCmd.Flags().Uint64("value", 0, "value")
	cmd.AddCommand(txCmd)

	traceCmd := &cobra.Command{
		Use:   "record-trace",
		Short: "Record a network trace",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ForensicRecordTrace")
			peer, _ := cmd.Flags().GetString("peer")
			event, _ := cmd.Flags().GetString("event")
			trace := nodes.NetworkTrace{PeerID: peer, Event: event, Timestamp: time.Now()}
			if err := forensic.RecordNetworkTrace(trace); err != nil {
				printOutput(err.Error())
			} else {
				printOutput("trace recorded")
			}
		},
	}
	traceCmd.Flags().String("peer", "", "peer id")
	traceCmd.Flags().String("event", "", "event description")
	cmd.AddCommand(traceCmd)

	listTx := &cobra.Command{
		Use:   "txs",
		Short: "List recorded transactions",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ForensicListTx")
			printOutput(forensic.Transactions())
		},
	}
	cmd.AddCommand(listTx)

	listTrace := &cobra.Command{
		Use:   "traces",
		Short: "List recorded network traces",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ForensicListTrace")
			printOutput(forensic.NetworkTraces())
		},
	}
	cmd.AddCommand(listTrace)

	rootCmd.AddCommand(cmd)
}
