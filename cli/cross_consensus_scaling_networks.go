package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var consensusNetMgr = core.NewConsensusNetworkManager()

func init() {
	ccsnCmd := &cobra.Command{
		Use:   "cross-consensus",
		Short: "Manage cross-consensus scaling networks",
	}

	registerCmd := &cobra.Command{
		Use:   "register <source> <target>",
		Args:  cobra.ExactArgs(2),
		Short: "Register a new network",
		Run: func(cmd *cobra.Command, args []string) {
			id := consensusNetMgr.RegisterNetwork(args[0], args[1])
			fmt.Println(id)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List networks",
		Run: func(cmd *cobra.Command, args []string) {
			nets := consensusNetMgr.ListNetworks()
			for _, n := range nets {
				fmt.Printf("%d %s->%s\n", n.ID, n.SourceConsensus, n.TargetConsensus)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get network by ID",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("invalid id")
				return
			}
			n, err := consensusNetMgr.GetNetwork(id)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%d %s->%s\n", n.ID, n.SourceConsensus, n.TargetConsensus)
		},
	}

	ccsnCmd.AddCommand(registerCmd, listCmd, getCmd)
	rootCmd.AddCommand(ccsnCmd)
}
