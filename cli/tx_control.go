package cli

import (
	"encoding/hex"
	"fmt"
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
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			exec, _ := strconv.ParseInt(args[5], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			st := core.ScheduleTransaction(tx, time.Unix(exec, 0))
			fmt.Printf("scheduled tx %s for %s\n", st.Tx.ID, st.ExecuteAt)
		},
	}

	cancelCmd := &cobra.Command{
		Use:   "cancel [from] [to] [amount] [fee] [nonce] [execUnix]",
		Args:  cobra.ExactArgs(6),
		Short: "Schedule and then cancel a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			exec, _ := strconv.ParseInt(args[5], 10, 64)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			st := core.ScheduleTransaction(tx, time.Unix(exec, 0))
			canceled := core.CancelTransaction(st)
			fmt.Printf("canceled: %v\n", canceled)
		},
	}

	reverseCmd := &cobra.Command{
		Use:   "reverse [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(5),
		Short: "Apply then reverse a transaction on a fresh ledger",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			l := core.NewLedger()
			l.Credit(args[0], amt+fee)
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			if err := l.ApplyTransaction(tx); err != nil {
				fmt.Println("apply err:", err)
				return
			}
			if err := core.ReverseTransaction(l, tx); err != nil {
				fmt.Println("reverse err:", err)
				return
			}
			fmt.Printf("balances after reverse: %s=%d %s=%d\n", args[0], l.GetBalance(args[0]), args[1], l.GetBalance(args[1]))
		},
	}

	privateCmd := &cobra.Command{
		Use:   "private [from] [to] [amount] [fee] [nonce] [keyhex]",
		Args:  cobra.ExactArgs(6),
		Short: "Convert a transaction to private and decrypt it",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[2], 10, 64)
			fee, _ := strconv.ParseUint(args[3], 10, 64)
			nonce, _ := strconv.ParseUint(args[4], 10, 64)
			key, _ := hex.DecodeString(args[5])
			tx := core.NewTransaction(args[0], args[1], amt, fee, nonce)
			pt, err := core.ConvertToPrivate(tx, key)
			if err != nil {
				fmt.Println("private err:", err)
				return
			}
			dec, err := pt.Decrypt(key)
			if err != nil {
				fmt.Println("decrypt err:", err)
				return
			}
			fmt.Printf("private payload %x decrypted tx %s->%s\n", pt.Payload, dec.From, dec.To)
		},
	}

	receiptCmd := &cobra.Command{
		Use:   "receipt [txid] [status] [details]",
		Args:  cobra.ExactArgs(3),
		Short: "Generate and store a transaction receipt",
		Run: func(cmd *cobra.Command, args []string) {
			r := core.GenerateReceipt(args[0], args[1], args[2])
			store := core.NewReceiptStore()
			store.Store(r)
			got, _ := store.Get(args[0])
			fmt.Printf("receipt: %+v\n", got)
		},
	}

	controlCmd.AddCommand(scheduleCmd, cancelCmd, reverseCmd, privateCmd, receiptCmd)
	txCmd.AddCommand(controlCmd)
}
