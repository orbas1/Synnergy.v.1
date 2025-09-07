package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var investments = core.NewInvestmentRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn1401",
		Short: "SYN1401 investment tokens",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new investment",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			owner, _ := cmd.Flags().GetString("owner")
			principal, _ := cmd.Flags().GetUint64("principal")
			rate, _ := cmd.Flags().GetFloat64("rate")
			maturityUnix, _ := cmd.Flags().GetInt64("maturity")
			if id == "" || owner == "" || principal == 0 || rate == 0 || maturityUnix == 0 {
				return fmt.Errorf("id, owner, principal, rate and maturity must be provided")
			}
			maturity := time.Unix(maturityUnix, 0)
			if _, err := investments.Issue(id, owner, principal, rate, maturity); err != nil {
				return err
			}
			cmd.Println("investment issued")
			return nil
		},
	}
	issueCmd.Flags().String("id", "", "investment id")
	issueCmd.Flags().String("owner", "", "owner")
	issueCmd.Flags().Uint64("principal", 0, "principal")
	issueCmd.Flags().Float64("rate", 0, "annual rate")
	issueCmd.Flags().Int64("maturity", 0, "maturity unix time")
	issueCmd.MarkFlagRequired("id")
	issueCmd.MarkFlagRequired("owner")
	issueCmd.MarkFlagRequired("principal")
	issueCmd.MarkFlagRequired("rate")
	issueCmd.MarkFlagRequired("maturity")
	cmd.AddCommand(issueCmd)

	accrueCmd := &cobra.Command{
		Use:   "accrue <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Accrue interest to now",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := investments.Accrue(args[0], time.Now())
			if err != nil {
				return err
			}
			cmd.Println(amt)
			return nil
		},
	}
	cmd.AddCommand(accrueCmd)

	redeemCmd := &cobra.Command{
		Use:   "redeem <id> <owner>",
		Args:  cobra.ExactArgs(2),
		Short: "Redeem an investment",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := investments.Redeem(args[0], args[1], time.Now())
			if err != nil {
				return err
			}
			cmd.Println(amt)
			return nil
		},
	}
	cmd.AddCommand(redeemCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get investment info",
		RunE: func(cmd *cobra.Command, args []string) error {
			rec, ok := investments.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			cmd.Printf("%s owner:%s principal:%d accrued:%d\n", rec.ID, rec.Owner, rec.Principal, rec.Accrued)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
