package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var policy *core.TokenInsurancePolicy

func init() {
	cmd := &cobra.Command{
		Use:   "syn2900",
		Short: "Token insurance policies",
	}

	issueCmd := &cobra.Command{
		Use:   "issue",
		Short: "Issue a new policy",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
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
			policy = core.NewTokenInsurancePolicy(id, holder, coverage, premium, payout, deductible, limit, start, end)
			fmt.Println("policy issued")
		},
	}
	issueCmd.Flags().String("id", "", "policy id")
	issueCmd.Flags().String("holder", "", "policy holder")
	issueCmd.Flags().String("coverage", "", "coverage type")
	issueCmd.Flags().Uint64("premium", 0, "premium amount")
	issueCmd.Flags().Uint64("payout", 0, "payout amount")
	issueCmd.Flags().Uint64("deductible", 0, "deductible")
	issueCmd.Flags().Uint64("limit", 0, "limit")
	issueCmd.Flags().String("start", time.Now().Format(time.RFC3339), "start time")
	issueCmd.Flags().String("end", time.Now().Add(24*time.Hour).Format(time.RFC3339), "end time")
	cmd.AddCommand(issueCmd)

	activeCmd := &cobra.Command{
		Use:   "active",
		Short: "Check if policy is active",
		Run: func(cmd *cobra.Command, args []string) {
			if policy == nil {
				fmt.Println("no policy")
				return
			}
			fmt.Println(policy.IsActive(time.Now()))
		},
	}
	cmd.AddCommand(activeCmd)

	claimCmd := &cobra.Command{
		Use:   "claim",
		Short: "Claim the policy",
		Run: func(cmd *cobra.Command, args []string) {
			if policy == nil {
				fmt.Println("no policy")
				return
			}
			amt, err := policy.Claim(time.Now())
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("payout %d\n", amt)
		},
	}
	cmd.AddCommand(claimCmd)

	rootCmd.AddCommand(cmd)
}
