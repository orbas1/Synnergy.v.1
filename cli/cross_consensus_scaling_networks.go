package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var consensusNetMgr = core.NewConsensusNetworkManager()

func init() {
	ccsnCmd := &cobra.Command{
		Use:   "cross-consensus",
		Short: "Manage cross-consensus scaling networks",
	}

	var registerJSON bool
	registerCmd := &cobra.Command{
		Use:   "register <source> <target>",
		Args:  cobra.ExactArgs(2),
		Short: "Register a new network",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := consensusNetMgr.RegisterNetwork(args[0], args[1])
			gas := synnergy.GasCost("RegisterConsensusNetwork")
			if registerJSON {
				enc, _ := json.Marshal(map[string]interface{}{"id": id, "gas": gas})
				fmt.Println(string(enc))
				return nil
			}
			fmt.Println(id)
			return nil
		},
	}
	registerCmd.Flags().BoolVar(&registerJSON, "json", false, "output as JSON")

	var listJSON bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List networks",
		Run: func(cmd *cobra.Command, args []string) {
			nets := consensusNetMgr.ListNetworks()
			if listJSON {
				enc, _ := json.Marshal(nets)
				fmt.Println(string(enc))
				return
			}
			for _, n := range nets {
				fmt.Printf("%d %s->%s\n", n.ID, n.SourceConsensus, n.TargetConsensus)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	var getJSON bool
	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get network by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			n, err := consensusNetMgr.GetNetwork(id)
			if err != nil {
				return err
			}
			if getJSON {
				enc, _ := json.Marshal(n)
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("%d %s->%s\n", n.ID, n.SourceConsensus, n.TargetConsensus)
			return nil
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	var removeJSON bool
	removeCmd := &cobra.Command{
		Use:   "remove <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove network",
		Run: func(cmd *cobra.Command, args []string) {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("invalid id")
				return
			}
			if err := consensusNetMgr.RemoveNetwork(id); err != nil {
				fmt.Println(err)
				return
			}
			gas := synnergy.GasCost("RemoveConsensusNetwork")
			if removeJSON {
				enc, _ := json.Marshal(map[string]interface{}{"gas": gas})
				fmt.Println(string(enc))
				return
			}
			fmt.Println("removed")
		},
	}
	removeCmd.Flags().BoolVar(&removeJSON, "json", false, "output as JSON")

	ccsnCmd.AddCommand(registerCmd, listCmd, getCmd, removeCmd)
	rootCmd.AddCommand(ccsnCmd)
}
