package cli

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
	ierr "synnergy/internal/errors"
)

var validatorMgr = core.NewValidatorManager(100)

func init() {
	cmd := &cobra.Command{
		Use:   "validator",
		Short: "Manage consensus validators",
	}

	addCmd := &cobra.Command{
		Use:   "add [addr] [stake]",
		Args:  cobra.ExactArgs(2),
		Short: "Register validator with stake",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("AddValidator")
			stake, _ := strconv.ParseUint(args[1], 10, 64)
			if err := validatorMgr.Add(context.Background(), args[0], stake); err != nil {
				if e, ok := err.(*ierr.Error); ok {
					printOutput(map[string]any{"error": e.Message, "code": e.Code})
				} else {
					printOutput(map[string]any{"error": err.Error()})
				}
				return
			}
			printOutput(map[string]any{"status": "added", "address": args[0]})
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove validator",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("RemoveValidator")
			validatorMgr.Remove(context.Background(), args[0])
			printOutput(map[string]any{"status": "removed", "address": args[0]})
		},
	}

	slashCmd := &cobra.Command{
		Use:   "slash [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Slash validator stake",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SlashValidator")
			validatorMgr.Slash(context.Background(), args[0])
			printOutput(map[string]any{"status": "slashed", "address": args[0]})
		},
	}

	eligibleCmd := &cobra.Command{
		Use:   "eligible",
		Short: "List eligible validators",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Eligible")
			printOutput(validatorMgr.Eligible())
		},
	}

	stakeCmd := &cobra.Command{
		Use:   "stake [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show validator stake",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Stake")
			printOutput(map[string]uint64{"stake": validatorMgr.Stake(args[0])})
		},
	}

	setMinCmd := &cobra.Command{
		Use:   "set-min [stake]",
		Args:  cobra.ExactArgs(1),
		Short: "Set minimum stake requirement",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("SetMinStake")
			s, _ := strconv.ParseUint(args[0], 10, 64)
			validatorMgr = core.NewValidatorManager(s)
			printOutput(map[string]any{"minStake": s})
		},
	}

	cmd.AddCommand(addCmd, removeCmd, slashCmd, eligibleCmd, stakeCmd, setMinCmd)
	rootCmd.AddCommand(cmd)
}
