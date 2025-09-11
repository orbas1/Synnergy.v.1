package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn2100 = core.NewTradeFinanceToken()

func parseTime(t string) time.Time {
	if t == "" {
		return time.Now()
	}
	tt, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return time.Now()
	}
	return tt
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn2100",
		Short: "Trade finance token operations",
	}

	regCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a financial document",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			issuer, _ := cmd.Flags().GetString("issuer")
			recipient, _ := cmd.Flags().GetString("recipient")
			amount, _ := cmd.Flags().GetUint64("amount")
			if id == "" || issuer == "" || recipient == "" || amount == 0 {
				return fmt.Errorf("id, issuer, recipient and amount are required")
			}
			issueStr, _ := cmd.Flags().GetString("issue")
			dueStr, _ := cmd.Flags().GetString("due")
			desc, _ := cmd.Flags().GetString("desc")
			if err := syn2100.RegisterDocument(id, issuer, recipient, amount, parseTime(issueStr), parseTime(dueStr), desc); err != nil {
				return err
			}
			fmt.Println("document registered")
			return nil
		},
	}
	regCmd.Flags().String("id", "", "document id")
	regCmd.Flags().String("issuer", "", "issuer")
	regCmd.Flags().String("recipient", "", "recipient")
	regCmd.Flags().Uint64("amount", 0, "amount")
	regCmd.Flags().String("issue", "", "issue time RFC3339")
	regCmd.Flags().String("due", "", "due time RFC3339")
	regCmd.Flags().String("desc", "", "description")
	cmd.AddCommand(regCmd)

	finCmd := &cobra.Command{
		Use:   "finance <docID> <financier>",
		Short: "Finance a document",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := syn2100.FinanceDocument(args[0], args[1]); err != nil {
				return err
			}
			fmt.Println("document financed")
			return nil
		},
	}
	cmd.AddCommand(finCmd)

	getCmd := &cobra.Command{
		Use:   "get <docID>",
		Short: "Get document info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			d, ok := syn2100.GetDocument(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%+v\n", *d)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List documents",
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range syn2100.ListDocuments() {
				fmt.Printf("%s %s->%s %d financed:%v\n", d.DocID, d.Issuer, d.Recipient, d.Amount, d.Financed)
			}
		},
	}
	cmd.AddCommand(listCmd)

	addLiq := &cobra.Command{
		Use:   "add-liquidity <addr> <amt>",
		Short: "Add liquidity",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			amt := uint64(0)
			fmt.Sscanf(args[1], "%d", &amt)
			if amt == 0 {
				return fmt.Errorf("amount must be greater than zero")
			}
			syn2100.AddLiquidity(args[0], amt)
			fmt.Println("liquidity added")
			return nil
		},
	}
	cmd.AddCommand(addLiq)

	remLiq := &cobra.Command{
		Use:   "remove-liquidity <addr> <amt>",
		Short: "Remove liquidity",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			amt := uint64(0)
			fmt.Sscanf(args[1], "%d", &amt)
			if err := syn2100.RemoveLiquidity(args[0], amt); err != nil {
				return err
			}
			fmt.Println("liquidity removed")
			return nil
		},
	}
	cmd.AddCommand(remLiq)

	liqCmd := &cobra.Command{
		Use:   "liquidity",
		Short: "List liquidity pools",
		Run: func(cmd *cobra.Command, args []string) {
			var sb strings.Builder
			for addr, amt := range syn2100.Liquidity {
				sb.WriteString(fmt.Sprintf("%s:%d ", addr, amt))
			}
			fmt.Println(sb.String())
		},
	}
	cmd.AddCommand(liqCmd)

	rootCmd.AddCommand(cmd)
}
