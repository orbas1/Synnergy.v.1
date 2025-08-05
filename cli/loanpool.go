package cli

import (
	"encoding/json"
	"fmt"
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
			amt, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			id, err := loanPool.SubmitProposal(args[0], args[1], args[2], amt, args[4])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("proposal submitted", id)
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [voter] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote on a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[1], 10, 64)
			if err := loanPool.VoteProposal(args[0], id); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	disburseCmd := &cobra.Command{
		Use:   "disburse [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Disburse an approved proposal",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if err := loanPool.Disburse(id); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	tickCmd := &cobra.Command{
		Use:   "tick",
		Short: "Process proposals and update approvals",
		Run: func(cmd *cobra.Command, args []string) {
			loanPool.Tick()
			fmt.Println("processed")
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Display a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if view, ok := loanPool.ProposalInfo(id); ok {
				b, _ := json.MarshalIndent(view, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		Run: func(cmd *cobra.Command, args []string) {
			views := loanPool.ProposalViews()
			b, _ := json.MarshalIndent(views, "", "  ")
			fmt.Println(string(b))
		},
	}

	cancelCmd := &cobra.Command{
		Use:   "cancel [creator] [id]",
		Args:  cobra.ExactArgs(2),
		Short: "Cancel an active proposal",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[1], 10, 64)
			if err := loanPool.CancelProposal(args[0], id); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	extendCmd := &cobra.Command{
		Use:   "extend [creator] [id] [hrs]",
		Args:  cobra.ExactArgs(3),
		Short: "Extend voting deadline",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[1], 10, 64)
			hrs, _ := strconv.Atoi(args[2])
			if err := loanPool.ExtendProposal(args[0], id, hrs); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	lpCmd.AddCommand(submitCmd, voteCmd, disburseCmd, tickCmd, getCmd, listCmd, cancelCmd, extendCmd)
	rootCmd.AddCommand(lpCmd)
}
