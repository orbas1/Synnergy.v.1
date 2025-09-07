package cli

import (
	"encoding/hex"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	controlCmd := &cobra.Command{
		Use:   "control",
		Short: "Advanced transaction controls",
	}

	scheduleCmd := &cobra.Command{
		Use:   "schedule [from] [to] [amount] [fee] [nonce] [execUnix]",
		Args:  cobra.ExactArgs(6),
		Short: "Schedule a transaction for future execution",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("TxSchedule")
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			exec, _ := strconv.ParseInt(args[5], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			st := core.ScheduleTransaction(tx, time.Unix(exec, 0))
			printOutput(map[string]any{"txID": st.Tx.ID, "executeAt": st.ExecuteAt})
		},
	}

	cancelCmd := &cobra.Command{
		Use:   "cancel [from] [to] [amount] [fee] [nonce] [execUnix]",
		Args:  cobra.ExactArgs(6),
		Short: "Schedule and then cancel a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("TxCancel")
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			exec, _ := strconv.ParseInt(args[5], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			st := core.ScheduleTransaction(tx, time.Unix(exec, 0))
			canceled := core.CancelTransaction(st)
			printOutput(map[string]any{"canceled": canceled})
		},
	}

	reverseCmd := &cobra.Command{
		Use:   "reverse [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Apply then reverse a transaction on a fresh ledger",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("TxReverse")
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			l := core.NewLedger()
			l.Credit(args[0], amt+fee)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			if err := l.ApplyTransaction(tx); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			if err := core.ReverseTransaction(l, tx); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"balance": map[string]uint64{args[0]: l.GetBalance(args[0]), args[1]: l.GetBalance(args[1])}})
		},
	}

	privateCmd := &cobra.Command{
		Use:   "private [from] [to] [amount] [fee] [nonce] [keyhex]",
		Args:  cobra.ExactArgs(6),
		Short: "Convert a transaction to private and decrypt it",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("TxPrivate")
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			key, _ := hex.DecodeString(args[5])
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			pt, err := core.ConvertToPrivate(tx, key)
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			dec, err := pt.Decrypt(key)
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"payload": hex.EncodeToString(pt.Payload), "from": dec.From, "to": dec.To})
		},
	}

	receiptCmd := &cobra.Command{
		Use:   "receipt [txid] [status] [details]",
		Args:  cobra.ExactArgs(3),
		Short: "Generate and store a transaction receipt",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("TxReceipt")
			r := core.GenerateReceipt(args[0], args[1], args[2])
			store := core.NewReceiptStore()
			store.Store(r)
			got, _ := store.Get(args[0])
			printOutput(got)
		},
	}

	controlCmd.AddCommand(scheduleCmd, cancelCmd, reverseCmd, privateCmd, receiptCmd)
	txCmd.AddCommand(controlCmd)
}
