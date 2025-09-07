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
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			stakePenaltyMgr.Slash(stakingNode, args[0], amt)
			printOutput(map[string]string{"status": "slashed"})
			return nil
		},
	}

	rewardCmd := &cobra.Command{
		Use:   "reward <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Reward staked tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			stakePenaltyMgr.Reward(stakingNode, args[0], amt)
			printOutput(map[string]string{"status": "rewarded"})
			return nil
		},
	}

	cmd.AddCommand(slashCmd, rewardCmd)
	rootCmd.AddCommand(cmd)
}
