package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var shardMgr = core.NewShardManager(2)

func init() {
	shardCmd := &cobra.Command{
		Use:   "sharding",
		Short: "Manage sharded network state",
	}

	leaderCmd := &cobra.Command{Use: "leader"}

	leaderGetCmd := &cobra.Command{
		Use:   "get <shardID>",
		Args:  cobra.ExactArgs(1),
		Short: "Show leader for a shard",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			addr, ok := shardMgr.GetLeader(id)
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Println(addr)
		},
	}

	leaderSetCmd := &cobra.Command{
		Use:   "set <shardID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Set leader for a shard",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			shardMgr.SetLeader(id, args[1])
		},
	}

	leaderCmd.AddCommand(leaderGetCmd, leaderSetCmd)

	mapCmd := &cobra.Command{
		Use:   "map",
		Short: "List shard leaders",
		Run: func(cmd *cobra.Command, args []string) {
			for id, addr := range shardMgr.LeaderMap() {
				fmt.Printf("%d %s\n", id, addr)
			}
		},
	}

	submitCmd := &cobra.Command{
		Use:   "submit <fromShard> <toShard> <txHash>",
		Args:  cobra.ExactArgs(3),
		Short: "Submit cross-shard tx header",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := strconv.Atoi(args[0])
			to, _ := strconv.Atoi(args[1])
			shardMgr.SubmitCrossShardTx(from, to, args[2])
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull <shardID>",
		Args:  cobra.ExactArgs(1),
		Short: "Pull receipts for shard",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			for _, tx := range shardMgr.PullReceipts(id) {
				fmt.Println(tx)
			}
		},
	}

	reshardCmd := &cobra.Command{
		Use:   "reshard <bits>",
		Args:  cobra.ExactArgs(1),
		Short: "Set shard bit-size",
		Run: func(cmd *cobra.Command, args []string) {
			bits, _ := strconv.Atoi(args[0])
			shardMgr.Reshard(uint8(bits))
		},
	}

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance <threshold>",
		Args:  cobra.ExactArgs(1),
		Short: "List shards exceeding load",
		Run: func(cmd *cobra.Command, args []string) {
			th, _ := strconv.Atoi(args[0])
			for _, id := range shardMgr.Rebalance(th) {
				fmt.Println(id)
			}
		},
	}

	shardCmd.AddCommand(leaderCmd, mapCmd, submitCmd, pullCmd, reshardCmd, rebalanceCmd)
	rootCmd.AddCommand(shardCmd)
}
