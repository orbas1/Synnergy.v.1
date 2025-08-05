package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var lifeReg = tokens.NewLifePolicyRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn2800",
		Short: "Life insurance policies",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue policy",
		Run: func(cmd *cobra.Command, args []string) {
			insured, _ := cmd.Flags().GetString("insured")
			beneficiary, _ := cmd.Flags().GetString("beneficiary")
			coverage, _ := cmd.Flags().GetUint64("coverage")
			premium, _ := cmd.Flags().GetUint64("premium")
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			start, _ := time.Parse(time.RFC3339, startStr)
			end, _ := time.Parse(time.RFC3339, endStr)
			p := lifeReg.IssuePolicy(insured, beneficiary, coverage, premium, start, end)
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
		Use:   "pay-premium <policyID> <amount>",
		Short: "Record premium payment",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := lifeReg.PayPremium(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("payment recorded")
			}
		},
	}
	cmd.AddCommand(payCmd)

	claimCmd := &cobra.Command{
		Use:   "file-claim <policyID> <amount>",
		Short: "File a claim",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			c, err := lifeReg.FileClaim(args[0], amt)
			if err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println(c.ClaimID)
			}
		},
	}
	cmd.AddCommand(claimCmd)

	getCmd := &cobra.Command{
		Use:   "get <policyID>",
		Short: "Get policy info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := lifeReg.GetPolicy(args[0])
			if !ok {
				fmt.Println("policy not found")
				return
			}
			fmt.Printf("%+v\n", *p)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List policies",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range lifeReg.ListPolicies() {
				fmt.Printf("%+v\n", *p)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
