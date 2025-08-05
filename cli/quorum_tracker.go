package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	quorumTracker  *core.QuorumTracker
	quorumRequired int
)

func init() {
	cmd := &cobra.Command{
		Use:   "quorum",
		Short: "Manage quorum tracker",
	}

	initCmd := &cobra.Command{
		Use:   "init [total] [threshold]",
		Args:  cobra.ExactArgs(2),
		Short: "Initialise a quorum tracker",
		Run: func(cmd *cobra.Command, args []string) {
			t, _ := strconv.Atoi(args[1])
			quorumRequired = t
			quorumTracker = core.NewQuorumTracker(t)
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Record a vote from an address",
		Run: func(cmd *cobra.Command, args []string) {
			if quorumTracker != nil {
				quorumTracker.Join(args[0])
			}
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check",
		Short: "Check if quorum is reached",
		Run: func(cmd *cobra.Command, args []string) {
			if quorumTracker == nil {
				fmt.Println(false)
				return
			}
			fmt.Println(quorumTracker.Reached())
		},
	}

	resetCmd := &cobra.Command{
		Use:   "reset",
		Short: "Clear all recorded votes",
		Run: func(cmd *cobra.Command, args []string) {
			if quorumTracker != nil {
				quorumTracker = core.NewQuorumTracker(quorumRequired)
			}
		},
	}

	cmd.AddCommand(initCmd, voteCmd, checkCmd, resetCmd)
	rootCmd.AddCommand(cmd)
}
