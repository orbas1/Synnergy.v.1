package cli

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			insured, _ := cmd.Flags().GetString("insured")
			beneficiary, _ := cmd.Flags().GetString("beneficiary")
			coverage, _ := cmd.Flags().GetUint64("coverage")
			premium, _ := cmd.Flags().GetUint64("premium")
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			if insured == "" || beneficiary == "" || coverage == 0 || premium == 0 {
				return errors.New("insured, beneficiary, coverage and premium are required")
			}
			start, err := time.Parse(time.RFC3339, startStr)
			if err != nil {
				return err
			}
			end, err := time.Parse(time.RFC3339, endStr)
			if err != nil {
				return err
			}
			p, err := lifeRegistry.IssuePolicy(insured, beneficiary, coverage, premium, start, end)
			if err != nil {
				return err
			}
			fmt.Println(p.PolicyID)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var amt uint64
			if _, err := fmt.Sscanf(args[1], "%d", &amt); err != nil {
				return err
			}
			return lifeRegistry.PayPremium(args[0], amt)
		},
	}
	cmd.AddCommand(payCmd)

	claimCmd := &cobra.Command{
		Use:   "claim <policy> <amount>",
		Short: "File a claim",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var amt uint64
			if _, err := fmt.Sscanf(args[1], "%d", &amt); err != nil {
				return err
			}
			_, err := lifeRegistry.FileClaim(args[0], amt)
			return err
		},
	}
	cmd.AddCommand(claimCmd)

	getCmd := &cobra.Command{
		Use:   "get <policy>",
		Short: "Get policy info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, ok := lifeRegistry.GetPolicy(args[0])
			if !ok {
				return errors.New("policy not found")
			}
			fmt.Printf("ID:%s Insured:%s Beneficiary:%s Coverage:%d Premium:%d Paid:%d\n", p.PolicyID, p.Insured, p.Beneficiary, p.Coverage, p.Premium, p.PaidPremium)
			for _, c := range p.Claims {
				fmt.Printf("claim %s %d %s settled:%t\n", c.ClaimID, c.Amount, c.Time.Format(time.RFC3339), c.Settled)
			}
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List policies",
		RunE: func(cmd *cobra.Command, args []string) error {
			policies := lifeRegistry.ListPolicies()
			for _, p := range policies {
				fmt.Printf("%s %s %s %d %d\n", p.PolicyID, p.Insured, p.Beneficiary, p.Coverage, p.Premium)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	deactivateCmd := &cobra.Command{
		Use:   "deactivate <policy>",
		Short: "Deactivate a life policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lifeRegistry.Deactivate(args[0])
		},
	}
	cmd.AddCommand(deactivateCmd)

	rootCmd.AddCommand(cmd)
}
