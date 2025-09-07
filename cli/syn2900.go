package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var insuranceRegistry = tokens.NewInsuranceRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn2900",
		Short: "General insurance policies",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			holder, _ := cmd.Flags().GetString("holder")
			coverage, _ := cmd.Flags().GetString("coverage")
			premium, _ := cmd.Flags().GetUint64("premium")
			payout, _ := cmd.Flags().GetUint64("payout")
			deductible, _ := cmd.Flags().GetUint64("deductible")
			limit, _ := cmd.Flags().GetUint64("limit")
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			if holder == "" || coverage == "" || premium == 0 || payout == 0 {
				return errors.New("holder, coverage, premium and payout are required")
			}
			start, err := time.Parse(time.RFC3339, startStr)
			if err != nil {
				return err
			}
			end, err := time.Parse(time.RFC3339, endStr)
			if err != nil {
				return err
			}
			p, err := insuranceRegistry.IssuePolicy(holder, coverage, premium, payout, deductible, limit, start, end)
			if err != nil {
				return err
			}
			fmt.Println(p.PolicyID)
			return nil
		},
	}
	issueCmd.Flags().String("holder", "", "policy holder")
	issueCmd.Flags().String("coverage", "", "coverage type")
	issueCmd.Flags().Uint64("premium", 0, "premium amount")
	issueCmd.Flags().Uint64("payout", 0, "payout amount")
	issueCmd.Flags().Uint64("deductible", 0, "deductible")
	issueCmd.Flags().Uint64("limit", 0, "limit")
	issueCmd.Flags().String("start", time.Now().Format(time.RFC3339), "start time")
	issueCmd.Flags().String("end", time.Now().Add(24*time.Hour).Format(time.RFC3339), "end time")
	cmd.AddCommand(issueCmd)

	claimCmd := &cobra.Command{
		Use:   "claim <policy> <desc> <amount>",
		Short: "File a claim",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			var amt uint64
			if _, err := fmt.Sscanf(args[2], "%d", &amt); err != nil {
				return err
			}
			_, err := insuranceRegistry.FileClaim(args[0], args[1], amt)
			return err
		},
	}
	cmd.AddCommand(claimCmd)

	getCmd := &cobra.Command{
		Use:   "get <policy>",
		Short: "Get policy info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, ok := insuranceRegistry.GetPolicy(args[0])
			if !ok {
				return errors.New("policy not found")
			}
			fmt.Printf("ID:%s Holder:%s Coverage:%s Premium:%d Payout:%d Active:%t\n", p.PolicyID, p.Holder, p.Coverage, p.Premium, p.Payout, p.Active)
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
			policies := insuranceRegistry.ListPolicies()
			for _, p := range policies {
				fmt.Printf("%s %s %s %d %d\n", p.PolicyID, p.Holder, p.Coverage, p.Premium, p.Payout)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	deactivateCmd := &cobra.Command{
		Use:   "deactivate <policy>",
		Short: "Deactivate a policy",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return insuranceRegistry.Deactivate(args[0])
		},
	}
	cmd.AddCommand(deactivateCmd)

	rootCmd.AddCommand(cmd)
}
