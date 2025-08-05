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
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			owner, _ := cmd.Flags().GetString("owner")
			principal, _ := cmd.Flags().GetUint64("principal")
			rate, _ := cmd.Flags().GetFloat64("rate")
			maturityUnix, _ := cmd.Flags().GetInt64("maturity")
			maturity := time.Unix(maturityUnix, 0)
			if _, err := investments.Issue(id, owner, principal, rate, maturity); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("investment issued")
		},
	}
	issueCmd.Flags().String("id", "", "investment id")
	issueCmd.Flags().String("owner", "", "owner")
	issueCmd.Flags().Uint64("principal", 0, "principal")
	issueCmd.Flags().Float64("rate", 0, "annual rate")
	issueCmd.Flags().Int64("maturity", time.Now().Unix(), "maturity unix time")
	cmd.AddCommand(issueCmd)

	accrueCmd := &cobra.Command{
		Use:   "accrue <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Accrue interest to now",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := investments.Accrue(args[0], time.Now())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(amt)
		},
	}
	cmd.AddCommand(accrueCmd)

	redeemCmd := &cobra.Command{
		Use:   "redeem <id> <owner>",
		Args:  cobra.ExactArgs(2),
		Short: "Redeem an investment",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := investments.Redeem(args[0], args[1], time.Now())
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(amt)
		},
	}
	cmd.AddCommand(redeemCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get investment info",
		Run: func(cmd *cobra.Command, args []string) {
			rec, ok := investments.Get(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%s owner:%s principal:%d accrued:%d\n", rec.ID, rec.Owner, rec.Principal, rec.Accrued)
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
