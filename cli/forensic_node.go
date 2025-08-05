package cli

import (
	"fmt"
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
			hash, _ := cmd.Flags().GetString("hash")
			from, _ := cmd.Flags().GetString("from")
			to, _ := cmd.Flags().GetString("to")
			value, _ := cmd.Flags().GetUint64("value")
			tx := nodes.TransactionLite{Hash: hash, From: from, To: to, Value: value, Timestamp: time.Now()}
			if err := forensic.RecordTransaction(tx); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("transaction recorded")
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
			peer, _ := cmd.Flags().GetString("peer")
			event, _ := cmd.Flags().GetString("event")
			trace := nodes.NetworkTrace{PeerID: peer, Event: event, Timestamp: time.Now()}
			if err := forensic.RecordNetworkTrace(trace); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("trace recorded")
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
			for _, tx := range forensic.Transactions() {
				fmt.Printf("%s %s->%s %d\n", tx.Hash, tx.From, tx.To, tx.Value)
			}
		},
	}
	cmd.AddCommand(listTx)

	listTrace := &cobra.Command{
		Use:   "traces",
		Short: "List recorded network traces",
		Run: func(cmd *cobra.Command, args []string) {
			for _, tr := range forensic.NetworkTraces() {
				fmt.Printf("%s %s %s\n", tr.PeerID, tr.Event, tr.Timestamp.Format(time.RFC3339))
			}
		},
	}
	cmd.AddCommand(listTrace)

	rootCmd.AddCommand(cmd)
}
