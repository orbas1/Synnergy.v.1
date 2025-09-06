package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var loanPool = core.NewLoanPool(1_000_000)

func init() {
	lpCmd := &cobra.Command{
		Use:   "loanpool",
		Short: "Manage loan pool proposals",
	}

	submitCmd := &cobra.Command{
		Use:   "submit [creator] [recipient] [type] [amount] [desc]",
		Args:  cobra.ExactArgs(5),
		Short: "Submit a loan proposal",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolSubmit")
			amt, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid amount"})
				return
			}
			id, err := loanPool.SubmitProposal(args[0], args[1], args[2], amt, args[4])
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "submitted", "id": id})
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [voter] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote on a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolVote")
			id, _ := strconv.ParseUint(args[1], 10, 64)
			if err := loanPool.VoteProposal(args[0], id); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput("voted")
		},
	}

	disburseCmd := &cobra.Command{
		Use:   "disburse [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Disburse an approved proposal",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolDisburse")
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if err := loanPool.Disburse(id); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput("disbursed")
		},
	}

	tickCmd := &cobra.Command{
		Use:   "tick",
		Short: "Process proposals and update approvals",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolTick")
			loanPool.Tick()
			printOutput("processed")
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Display a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolGet")
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if view, ok := loanPool.ProposalInfo(id); ok {
				printOutput(view)
			} else {
				printOutput(map[string]any{"error": "not found"})
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolList")
			views := loanPool.ProposalViews()
			printOutput(views)
		},
	}

	cancelCmd := &cobra.Command{
		Use:   "cancel [creator] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Cancel an active proposal",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolCancel")
			id, _ := strconv.ParseUint(args[1], 10, 64)
			if err := loanPool.CancelProposal(args[0], id); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput("cancelled")
		},
	}

	extendCmd := &cobra.Command{
		Use:   "extend [creator] [id] [hrs]",
		Args:  cobra.ExactArgs(3),
		Short: "Extend voting deadline",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanpoolExtend")
			id, _ := strconv.ParseUint(args[1], 10, 64)
			hrs, _ := strconv.Atoi(args[2])
			if err := loanPool.ExtendProposal(args[0], id, hrs); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput("extended")
		},
	}

	lpCmd.AddCommand(submitCmd, voteCmd, disburseCmd, tickCmd, getCmd, listCmd, cancelCmd, extendCmd)
	rootCmd.AddCommand(lpCmd)
}
