package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var stakingNode = core.NewStakingNode()

func init() {
	cmd := &cobra.Command{
		Use:   "staking_node",
		Short: "Manage staking balances",
	}

	stakeCmd := &cobra.Command{
		Use:   "stake <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Stake tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			stakingNode.Stake(args[0], amt)
			printOutput(map[string]any{"status": "staked", "address": args[0], "amount": amt})
			return nil
		},
	}

	unstakeCmd := &cobra.Command{
		Use:   "unstake <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Unstake tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			stakingNode.Unstake(args[0], amt)
			printOutput(map[string]any{"status": "unstaked", "address": args[0], "amount": amt})
			return nil
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show staked balance",
		RunE: func(cmd *cobra.Command, args []string) error {
			printOutput(stakingNode.Balance(args[0]))
			return nil
		},
	}

	totalCmd := &cobra.Command{
		Use:   "total",
		Short: "Show total staked tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			printOutput(stakingNode.TotalStaked())
			return nil
		},
	}

	cmd.AddCommand(stakeCmd, unstakeCmd, balanceCmd, totalCmd)
	rootCmd.AddCommand(cmd)
}
