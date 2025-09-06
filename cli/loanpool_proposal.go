package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	proposals      = make(map[uint64]*core.LoanProposal)
	nextProposalID uint64
)

func init() {
	cmd := &cobra.Command{
		Use:   "loanproposal",
		Short: "Work with standalone loan proposals",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "new [creator] [recipient] [type] [amount] [desc] [durationHours]",
		Args:  cobra.ExactArgs(6),
		Short: "Create a loan proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("LoanProposalNew")
			amt, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %s", args[3])
			}
			dur, err := strconv.Atoi(args[5])
			if err != nil {
				return fmt.Errorf("invalid duration: %s", args[5])
			}
			nextProposalID++
			p := core.NewLoanProposal(nextProposalID, args[0], args[1], args[2], amt, args[4], time.Duration(dur)*time.Hour)
			proposals[nextProposalID] = p
			printOutput(map[string]any{"status": "created", "id": nextProposalID})
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "vote [id] [voter]",
		Args:  cobra.ExactArgs(2),
		Short: "Cast a vote on a proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("LoanProposalVote")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %s", args[0])
			}
			if p, ok := proposals[id]; ok {
				p.Vote(args[1])
				printOutput(map[string]any{"status": "voted", "id": id})
				return nil
			}
			return fmt.Errorf("not found")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "votes [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Show vote count",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("LoanProposalVotes")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %s", args[0])
			}
			if p, ok := proposals[id]; ok {
				printOutput(map[string]int{"votes": p.VoteCount()})
				return nil
			}
			return fmt.Errorf("not found")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "expired [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if proposal has expired",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("LoanProposalExpired")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %s", args[0])
			}
			if p, ok := proposals[id]; ok {
				printOutput(map[string]bool{"expired": p.IsExpired(time.Now())})
				return nil
			}
			return fmt.Errorf("not found")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Show proposal details",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("LoanProposalGet")
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id: %s", args[0])
			}
			if p, ok := proposals[id]; ok {
				printOutput(p)
				return nil
			}
			return fmt.Errorf("not found")
		},
	})

	rootCmd.AddCommand(cmd)
}
