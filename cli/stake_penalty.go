package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var stakePenaltyMgr = core.NewStakePenaltyManager()

func init() {
	cmd := &cobra.Command{
		Use:   "stake_penalty",
		Short: "Apply staking penalties or rewards",
	}

	slashCmd := &cobra.Command{
		Use:   "slash <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Slash staked tokens",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			stakePenaltyMgr.Slash(stakingNode, args[0], amt)
			fmt.Println("slashed")
		},
	}

	rewardCmd := &cobra.Command{
		Use:   "reward <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Reward staked tokens",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			stakePenaltyMgr.Reward(stakingNode, args[0], amt)
			fmt.Println("rewarded")
		},
	}

	cmd.AddCommand(slashCmd, rewardCmd)
	rootCmd.AddCommand(cmd)
}
