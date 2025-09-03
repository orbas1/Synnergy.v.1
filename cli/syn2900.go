package cli

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			holder, _ := cmd.Flags().GetString("holder")
			coverage, _ := cmd.Flags().GetString("coverage")
			premium, _ := cmd.Flags().GetUint64("premium")
			payout, _ := cmd.Flags().GetUint64("payout")
			deductible, _ := cmd.Flags().GetUint64("deductible")
			limit, _ := cmd.Flags().GetUint64("limit")
			startStr, _ := cmd.Flags().GetString("start")
			endStr, _ := cmd.Flags().GetString("end")
			start, _ := time.Parse(time.RFC3339, startStr)
			end, _ := time.Parse(time.RFC3339, endStr)
			p, err := insuranceRegistry.IssuePolicy(holder, coverage, premium, payout, deductible, limit, start, end)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Println(p.PolicyID)
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
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if _, err := insuranceRegistry.FileClaim(args[0], args[1], amt); err != nil {
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
			p, ok := insuranceRegistry.GetPolicy(args[0])
			if !ok {
				fmt.Println("policy not found")
				return
			}
			fmt.Printf("ID:%s Holder:%s Coverage:%s Premium:%d Payout:%d Active:%t\n", p.PolicyID, p.Holder, p.Coverage, p.Premium, p.Payout, p.Active)
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
			policies := insuranceRegistry.ListPolicies()
			for _, p := range policies {
				fmt.Printf("%s %s %s %d %d\n", p.PolicyID, p.Holder, p.Coverage, p.Premium, p.Payout)
			}
		},
	}
	cmd.AddCommand(listCmd)

	deactivateCmd := &cobra.Command{
		Use:   "deactivate <policy>",
		Short: "Deactivate a policy",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := insuranceRegistry.Deactivate(args[0]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(deactivateCmd)

	rootCmd.AddCommand(cmd)
}
