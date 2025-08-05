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
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			stakingNode.Stake(args[0], amt)
			fmt.Println("staked")
		},
	}

	unstakeCmd := &cobra.Command{
		Use:   "unstake <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Unstake tokens",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			stakingNode.Unstake(args[0], amt)
			fmt.Println("unstaked")
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show staked balance",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(stakingNode.Balance(args[0]))
		},
	}

	totalCmd := &cobra.Command{
		Use:   "total",
		Short: "Show total staked tokens",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(stakingNode.TotalStaked())
		},
	}

	cmd.AddCommand(stakeCmd, unstakeCmd, balanceCmd, totalCmd)
	rootCmd.AddCommand(cmd)
}
