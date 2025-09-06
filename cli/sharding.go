package cli

import (
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
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingLeaderGet")
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			if addr, ok := shardMgr.GetLeader(id); ok {
				printOutput(map[string]string{"leader": addr})
			} else {
				printOutput(map[string]string{"error": "not found"})
			}
			return nil
		},
	}

	leaderSet := &cobra.Command{
		Use:   "set [shardID] [addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Set the leader address for a shard",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingLeaderSet")
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			shardMgr.SetLeader(id, args[1])
			printOutput(map[string]any{"status": "set", "id": id})
			return nil
		},
	}

	leaderCmd.AddCommand(leaderGet, leaderSet)

	mapCmd := &cobra.Command{
		Use:   "map",
		Short: "List shard-to-leader mappings",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingMap")
			m := shardMgr.LeaderMap()
			printOutput(m)
			return nil
		},
	}

	submitCmd := &cobra.Command{
		Use:   "submit [fromShard] [toShard] [txHash]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit a cross-shard transaction header",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingSubmit")
			from, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			to, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			shardMgr.SubmitCrossShardTx(from, to, args[2])
			printOutput(map[string]string{"status": "submitted"})
			return nil
		},
	}

	pullCmd := &cobra.Command{
		Use:   "pull [shardID]",
		Args:  cobra.ExactArgs(1),
		Short: "Pull receipts for a shard",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingPull")
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			rec := shardMgr.PullReceipts(id)
			printOutput(rec)
			return nil
		},
	}

	reshardCmd := &cobra.Command{
		Use:   "reshard [newBits]",
		Args:  cobra.ExactArgs(1),
		Short: "Increase the shard count",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingReshard")
			bits, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			shardMgr.Reshard(uint8(bits))
			printOutput(map[string]any{"status": "resharded", "bits": bits})
			return nil
		},
	}

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance [threshold]",
		Args:  cobra.ExactArgs(1),
		Short: "List shards exceeding load threshold",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("ShardingRebalance")
			th, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			heavy := shardMgr.Rebalance(th)
			printOutput(heavy)
			return nil
		},
	}

	cmd.AddCommand(leaderCmd, mapCmd, submitCmd, pullCmd, reshardCmd, rebalanceCmd)
	rootCmd.AddCommand(cmd)
}
