package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var loanApply = core.NewLoanPoolApply(loanPool)

func init() {
	appCmd := &cobra.Command{
		Use:   "loanpool_apply",
		Short: "Manage loan applications",
	}

	submitCmd := &cobra.Command{
		Use:   "submit [applicant] [amount] [termMonths] [purpose]",
		Args:  cobra.ExactArgs(4),
		Short: "Submit a loan application",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolApplySubmit")
			amt, _ := strconv.ParseUint(args[1], 10, 64)
			term, _ := strconv.ParseUint(args[2], 10, 32)
			id := loanApply.Submit(args[0], amt, uint32(term), args[3])
			printOutput(map[string]any{"status": "submitted", "id": id})
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [voter] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote on an application",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolApplyVote")
			id, _ := strconv.ParseUint(args[1], 10, 64)
			if err := loanApply.Vote(args[0], id); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput("voted")
		},
	}

	processCmd := &cobra.Command{
		Use:   "process",
		Short: "Process pending applications",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolApplyProcess")
			loanApply.Process()
			printOutput("processed")
		},
	}

	disburseCmd := &cobra.Command{
		Use:   "disburse [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Disburse an approved application",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolApplyDisburse")
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if err := loanApply.Disburse(id); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput("disbursed")
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Display an application",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolApplyGet")
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if view, ok := loanApply.ApplicationInfo(id); ok {
				printOutput(view)
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List applications",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolApplyList")
			views := loanApply.ApplicationViews()
			printOutput(views)
		},
	}

	appCmd.AddCommand(submitCmd, voteCmd, processCmd, disburseCmd, getCmd, listCmd)
	rootCmd.AddCommand(appCmd)
}
