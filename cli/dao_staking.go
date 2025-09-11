package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var daoStaking = core.NewDAOStaking(daoMgr)

func init() {
	stakingCmd := &cobra.Command{
		Use:   "dao-stake",
		Short: "DAO staking operations",
	}

	var stakeJSON bool
	var stakePub, stakeMsg, stakeSig string
	stakeCmd := &cobra.Command{
		Use:   "stake <daoID> <addr> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Stake tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_Stake")
			ok, err := VerifySignature(stakePub, stakeMsg, stakeSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid amount")
				return
			}
			if err := daoStaking.Stake(args[0], args[1], amt); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if stakeJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "staked"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "staked")
		},
	}
	stakeCmd.Flags().BoolVar(&stakeJSON, "json", false, "output as JSON")
	stakeCmd.Flags().StringVar(&stakePub, "pub", "", "hex encoded public key")
	stakeCmd.Flags().StringVar(&stakeMsg, "msg", "", "hex encoded message")
	stakeCmd.Flags().StringVar(&stakeSig, "sig", "", "hex encoded signature")

	var unstakeJSON bool
	var unstakePub, unstakeMsg, unstakeSig string
	unstakeCmd := &cobra.Command{
		Use:   "unstake <daoID> <addr> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Unstake tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_Unstake")
			ok, err := VerifySignature(unstakePub, unstakeMsg, unstakeSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid amount")
				return
			}
			if err := daoStaking.Unstake(args[0], args[1], amt); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if unstakeJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "unstaked"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "unstaked")
		},
	}
	unstakeCmd.Flags().BoolVar(&unstakeJSON, "json", false, "output as JSON")
	unstakeCmd.Flags().StringVar(&unstakePub, "pub", "", "hex encoded public key")
	unstakeCmd.Flags().StringVar(&unstakeMsg, "msg", "", "hex encoded message")
	unstakeCmd.Flags().StringVar(&unstakeSig, "sig", "", "hex encoded signature")

	var balanceJSON bool
	balanceCmd := &cobra.Command{
		Use:   "balance <daoID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Show staked balance",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_Staked")
			bal := daoStaking.Balance(args[0], args[1])
			if balanceJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]uint64{"balance": bal})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), bal)
		},
	}
	balanceCmd.Flags().BoolVar(&balanceJSON, "json", false, "output as JSON")

	var totalJSON bool
	totalCmd := &cobra.Command{
		Use:   "total <daoID>",
		Args:  cobra.ExactArgs(1),
		Short: "Show total staked tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAO_TotalStaked")
			t := daoStaking.TotalStaked(args[0])
			if totalJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]uint64{"total": t})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), t)
		},
	}
	totalCmd.Flags().BoolVar(&totalJSON, "json", false, "output as JSON")

	stakingCmd.AddCommand(stakeCmd, unstakeCmd, balanceCmd, totalCmd)
	rootCmd.AddCommand(stakingCmd)
}
