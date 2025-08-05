package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var swarm = core.NewSwarm()

func init() {
	cmd := &cobra.Command{
		Use:   "swarm",
		Short: "Manage swarms of nodes",
	}

	joinCmd := &cobra.Command{
		Use:   "join <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Join a node to the swarm",
		Run: func(cmd *cobra.Command, args []string) {
			n := core.NewNode(args[0], args[0], core.NewLedger())
			swarm.Join(n)
			fmt.Println("node joined")
		},
	}

	leaveCmd := &cobra.Command{
		Use:   "leave <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a node from the swarm",
		Run: func(cmd *cobra.Command, args []string) {
			swarm.Leave(args[0])
		},
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List peer IDs",
		Run: func(cmd *cobra.Command, args []string) {
			for _, id := range swarm.Peers() {
				fmt.Println(id)
			}
		},
	}

	broadcastCmd := &cobra.Command{
		Use:   "broadcast <from> <to> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Broadcast a transaction to all members",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			tx := core.NewTransaction(args[0], args[1], amt, 0, 0)
			swarm.Broadcast(tx)
		},
	}

	consensusCmd := &cobra.Command{
		Use:   "consensus",
		Short: "Start consensus on the swarm",
		Run: func(cmd *cobra.Command, args []string) {
			blocks := swarm.StartConsensus()
			fmt.Printf("mined %d blocks\n", len(blocks))
		},
	}

	cmd.AddCommand(joinCmd, leaveCmd, peersCmd, broadcastCmd, consensusCmd)
	rootCmd.AddCommand(cmd)
}
