package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var debtRegistry = tokens.NewDebtRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn845",
		Short: "Debt instrument tokens",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a debt token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			supply, _ := cmd.Flags().GetUint64("supply")
			id, _ := debtRegistry.CreateToken(name, symbol, owner, supply)
			fmt.Println(id)
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Uint64("supply", 0, "initial supply")
	cmd.AddCommand(createCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <token> <debtID> <borrower> <principal> <rate> <penalty> <due>",
		Short: "Issue a debt instrument",
		Args:  cobra.ExactArgs(7),
		Run: func(cmd *cobra.Command, args []string) {
			var principal uint64
			var rate, penalty float64
			fmt.Sscanf(args[3], "%d", &principal)
			fmt.Sscanf(args[4], "%f", &rate)
			fmt.Sscanf(args[5], "%f", &penalty)
			due, _ := time.Parse(time.RFC3339, args[6])
			if err := debtRegistry.IssueDebt(args[0], args[1], args[2], principal, rate, penalty, due); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(issueCmd)

	payCmd := &cobra.Command{
		Use:   "pay <token> <debtID> <amount>",
		Short: "Record a payment",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := debtRegistry.RecordPayment(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(payCmd)

	infoCmd := &cobra.Command{
		Use:   "info <token> <debtID>",
		Short: "Show debt info",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			d, err := debtRegistry.GetDebt(args[0], args[1])
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("DebtID:%s Borrower:%s Principal:%d Paid:%d Rate:%f Penalty:%f Due:%s\n", d.DebtID, d.Borrower, d.Principal, d.Paid, d.Rate, d.Penalty, d.Due.Format(time.RFC3339))
		},
	}
	cmd.AddCommand(infoCmd)

	rootCmd.AddCommand(cmd)
}
