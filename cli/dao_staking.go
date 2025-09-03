package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var daoStaking = core.NewDAOStaking()

func init() {
	stakingCmd := &cobra.Command{
		Use:   "dao-stake",
		Short: "DAO staking operations",
	}

	stakeCmd := &cobra.Command{
		Use:   "stake <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Stake tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_Stake")
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			daoStaking.Stake(args[0], amt)
		},
	}

	unstakeCmd := &cobra.Command{
		Use:   "unstake <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Unstake tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_Unstake")
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := daoStaking.Unstake(args[0], amt); err != nil {
				fmt.Println(err)
			}
		},
	}

	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show staked balance",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_Staked")
			fmt.Println(daoStaking.Balance(args[0]))
		},
	}

	totalCmd := &cobra.Command{
		Use:   "total",
		Short: "Show total staked tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_TotalStaked")
			fmt.Println(daoStaking.TotalStaked())
		},
	}

	stakingCmd.AddCommand(stakeCmd, unstakeCmd, balanceCmd, totalCmd)
	rootCmd.AddCommand(stakingCmd)
}
