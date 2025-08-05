package cli

import (
	"encoding/json"
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
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[3], 10, 64)
			dur, _ := strconv.Atoi(args[5])
			nextProposalID++
			p := core.NewLoanProposal(nextProposalID, args[0], args[1], args[2], amt, args[4], time.Duration(dur)*time.Hour)
			proposals[nextProposalID] = p
			fmt.Println(nextProposalID)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "vote [id] [voter]",
		Args:  cobra.ExactArgs(2),
		Short: "Cast a vote on a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if p, ok := proposals[id]; ok {
				p.Vote(args[1])
			} else {
				fmt.Println("not found")
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "votes [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Show vote count",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if p, ok := proposals[id]; ok {
				fmt.Println(p.VoteCount())
			} else {
				fmt.Println("not found")
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "expired [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if proposal has expired",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if p, ok := proposals[id]; ok {
				fmt.Println(p.IsExpired(time.Now()))
			} else {
				fmt.Println("not found")
			}
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Show proposal details",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if p, ok := proposals[id]; ok {
				b, _ := json.MarshalIndent(p, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	})

	rootCmd.AddCommand(cmd)
}
