package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
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
			stake, _ := strconv.ParseUint(args[1], 10, 64)
			if err := validatorMgr.Add(args[0], stake); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove validator",
		Run: func(cmd *cobra.Command, args []string) {
			validatorMgr.Remove(args[0])
		},
	}

	slashCmd := &cobra.Command{
		Use:   "slash [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Slash validator stake",
		Run: func(cmd *cobra.Command, args []string) {
			validatorMgr.Slash(args[0])
		},
	}

	eligibleCmd := &cobra.Command{
		Use:   "eligible",
		Short: "List eligible validators",
		Run: func(cmd *cobra.Command, args []string) {
			for addr, stake := range validatorMgr.Eligible() {
				fmt.Printf("%s:%d\n", addr, stake)
			}
		},
	}

	stakeCmd := &cobra.Command{
		Use:   "stake [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show validator stake",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(validatorMgr.Stake(args[0]))
		},
	}

	setMinCmd := &cobra.Command{
		Use:   "set-min [stake]",
		Args:  cobra.ExactArgs(1),
		Short: "Set minimum stake requirement",
		Run: func(cmd *cobra.Command, args []string) {
			s, _ := strconv.ParseUint(args[0], 10, 64)
			validatorMgr = core.NewValidatorManager(s)
		},
	}

	cmd.AddCommand(addCmd, removeCmd, slashCmd, eligibleCmd, stakeCmd, setMinCmd)
	rootCmd.AddCommand(cmd)
}
