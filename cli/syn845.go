package cli

import (
	"fmt"
	"strconv"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			supply, _ := cmd.Flags().GetUint64("supply")
			if name == "" || symbol == "" || owner == "" {
				return fmt.Errorf("name, symbol and owner required")
			}
			id, _ := debtRegistry.CreateToken(name, symbol, owner, supply)
			fmt.Fprintf(cmd.OutOrStdout(), "token created %s\n", id)
			return nil
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Uint64("supply", 0, "initial supply")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("owner")
	cmd.AddCommand(createCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <token> <debtID> <borrower> <principal> <rate> <penalty> <due>",
		Short: "Issue a debt instrument",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			principal, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid principal")
			}
			rate, err := strconv.ParseFloat(args[4], 64)
			if err != nil {
				return fmt.Errorf("invalid rate")
			}
			penalty, err := strconv.ParseFloat(args[5], 64)
			if err != nil {
				return fmt.Errorf("invalid penalty")
			}
			due, err := time.Parse(time.RFC3339, args[6])
			if err != nil {
				return fmt.Errorf("invalid due date")
			}
			if err := debtRegistry.IssueDebt(args[0], args[1], args[2], principal, rate, penalty, due); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "debt issued")
			return nil
		},
	}
	cmd.AddCommand(issueCmd)

	payCmd := &cobra.Command{
		Use:   "pay <token> <debtID> <amount>",
		Short: "Record a payment",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			if err := debtRegistry.RecordPayment(args[0], args[1], amt); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "payment recorded")
			return nil
		},
	}
	cmd.AddCommand(payCmd)

	infoCmd := &cobra.Command{
		Use:   "info <token> <debtID>",
		Short: "Show debt info",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := debtRegistry.GetDebt(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "DebtID:%s Borrower:%s Principal:%d Paid:%d Rate:%f Penalty:%f Due:%s\n", d.DebtID, d.Borrower, d.Principal, d.Paid, d.Rate, d.Penalty, d.Due.Format(time.RFC3339))
			return nil
		},
	}
	cmd.AddCommand(infoCmd)

	rootCmd.AddCommand(cmd)
}
