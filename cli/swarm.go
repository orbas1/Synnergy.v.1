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
		RunE: func(cmd *cobra.Command, args []string) error {
			n := core.NewNode(args[0], args[0], core.NewLedger())
			swarm.Join(n)
			printOutput(map[string]string{"status": "joined", "id": args[0]})
			return nil
		},
	}

	leaveCmd := &cobra.Command{
		Use:   "leave <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a node from the swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			swarm.Leave(args[0])
			printOutput(map[string]string{"status": "left", "id": args[0]})
			return nil
		},
	}

	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "List peer IDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			printOutput(swarm.Peers())
			return nil
		},
	}

	broadcastCmd := &cobra.Command{
		Use:   "broadcast <from> <to> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Broadcast a transaction to all members",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			tx := core.NewTransaction(args[0], args[1], amt, 0, 0)
			swarm.Broadcast(tx)
			printOutput(map[string]string{"status": "broadcast"})
			return nil
		},
	}

	consensusCmd := &cobra.Command{
		Use:   "consensus",
		Short: "Start consensus on the swarm",
		RunE: func(cmd *cobra.Command, args []string) error {
			blocks := swarm.StartConsensus()
			printOutput(map[string]int{"mined": len(blocks)})
			return nil
		},
	}

	cmd.AddCommand(joinCmd, leaveCmd, peersCmd, broadcastCmd, consensusCmd)
	rootCmd.AddCommand(cmd)
}
