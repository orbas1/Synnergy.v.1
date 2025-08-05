package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var shardMgr = core.NewShardManager(2)

func init() {
	cmd := &cobra.Command{
		Use:   "sharding",
		Short: "Manage shard assignments",
	}

	leaderCmd := &cobra.Command{
		Use:   "leader",
		Short: "Query or set shard leaders",
	}

	leaderGet := &cobra.Command{
		Use:   "get [shardID]",
		Args:  cobra.ExactArgs(1),
		Short: "Show the leader for a shard",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			if addr, ok := shardMgr.GetLeader(id); ok {
				fmt.Println(addr)
			} else {
				fmt.Println("not found")
			}
		},
	}

	leaderSet := &cobra.Command{
		Use:   "set [shardID] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Set the leader address for a shard",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			shardMgr.SetLeader(id, args[1])
		},
	}

	leaderCmd.AddCommand(leaderGet, leaderSet)

	mapCmd := &cobra.Command{
		Use:   "map",
		Short: "List shard-to-leader mappings",
		Run: func(cmd *cobra.Command, args []string) {
			m := shardMgr.LeaderMap()
			out, _ := json.MarshalIndent(m, "", "  ")
			fmt.Println(string(out))
		},
	}

	submitCmd := &cobra.Command{
		Use:   "submit [fromShard] [toShard] [txHash]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a cross-shard transaction header",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := strconv.Atoi(args[0])
			to, _ := strconv.Atoi(args[1])
			shardMgr.SubmitCrossShardTx(from, to, args[2])
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull [shardID]",
		Args:  cobra.ExactArgs(1),
		Short: "Pull receipts for a shard",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.Atoi(args[0])
			rec := shardMgr.PullReceipts(id)
			for _, tx := range rec {
				fmt.Println(tx)
			}
		},
	}

	reshardCmd := &cobra.Command{
		Use:   "reshard [newBits]",
		Args:  cobra.ExactArgs(1),
		Short: "Increase the shard count",
		Run: func(cmd *cobra.Command, args []string) {
			bits, _ := strconv.Atoi(args[0])
			shardMgr.Reshard(uint8(bits))
		},
	}

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance [threshold]",
		Args:  cobra.ExactArgs(1),
		Short: "List shards exceeding load threshold",
		Run: func(cmd *cobra.Command, args []string) {
			th, _ := strconv.Atoi(args[0])
			heavy := shardMgr.Rebalance(th)
			if len(heavy) == 0 {
				fmt.Println("none")
				return
			}
			for _, id := range heavy {
				fmt.Println(id)
			}
		},
	}

	cmd.AddCommand(leaderCmd, mapCmd, submitCmd, pullCmd, reshardCmd, rebalanceCmd)
	rootCmd.AddCommand(cmd)
}
