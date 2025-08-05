package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var lifeRegistry = tokens.NewLifePolicyRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn2800",
		Short: "Life insurance policies",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a life policy",
		Run: func(cmd *cobra.Command, args []string) {
			insured, _ := cmd.Flags().GetString("insured")
			beneficiary, _ := cmd.Flags().GetString("beneficiary")
			coverage, _ := cmd.Flags().GetUint64("coverage")
			premium, _ := cmd.Flags().GetUint64("premium")
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			start, _ := time.Parse(time.RFC3339, startStr)
			end, _ := time.Parse(time.RFC3339, endStr)
			p := lifeRegistry.IssuePolicy(insured, beneficiary, coverage, premium, start, end)
			fmt.Println(p.PolicyID)
		},
	}
	issueCmd.Flags().String("insured", "", "insured party")
	issueCmd.Flags().String("beneficiary", "", "beneficiary")
	issueCmd.Flags().Uint64("coverage", 0, "coverage amount")
	issueCmd.Flags().Uint64("premium", 0, "premium amount")
	issueCmd.Flags().String("start", time.Now().Format(time.RFC3339), "start time")
	issueCmd.Flags().String("end", time.Now().Add(24*time.Hour).Format(time.RFC3339), "end time")
	cmd.AddCommand(issueCmd)

	payCmd := &cobra.Command{
		Use:   "pay <policy> <amount>",
		Short: "Record premium payment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := lifeRegistry.PayPremium(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(payCmd)

	claimCmd := &cobra.Command{
		Use:   "claim <policy> <amount>",
		Short: "File a claim",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if _, err := lifeRegistry.FileClaim(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(claimCmd)

	getCmd := &cobra.Command{
		Use:   "get <policy>",
		Short: "Get policy info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := lifeRegistry.GetPolicy(args[0])
			if !ok {
				fmt.Println("policy not found")
				return
			}
			fmt.Printf("ID:%s Insured:%s Beneficiary:%s Coverage:%d Premium:%d Paid:%d\n", p.PolicyID, p.Insured, p.Beneficiary, p.Coverage, p.Premium, p.PaidPremium)
			for _, c := range p.Claims {
				fmt.Printf("claim %s %d %s settled:%t\n", c.ClaimID, c.Amount, c.Time.Format(time.RFC3339), c.Settled)
			}
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List policies",
		Run: func(cmd *cobra.Command, args []string) {
			policies := lifeRegistry.ListPolicies()
			for _, p := range policies {
				fmt.Printf("%s %s %s %d %d\n", p.PolicyID, p.Insured, p.Beneficiary, p.Coverage, p.Premium)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
