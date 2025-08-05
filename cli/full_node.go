package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
	nodes "synnergy/internal/nodes"
)

var fullNode *core.FullNode

func init() {
	cmd := &cobra.Command{
		Use:   "fullnode",
		Short: "Manage full node configuration",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a full node",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			modeStr, _ := cmd.Flags().GetString("mode")
			var mode core.FullNodeMode
			if modeStr == "archive" {
				mode = core.FullNodeModeArchive
			} else {
				mode = core.FullNodeModePruned
			}
			fullNode = core.NewFullNode(nodes.Address(id), mode)
			fmt.Println("full node created")
		},
	}
	createCmd.Flags().String("id", "", "node id")
	createCmd.Flags().String("mode", "archive", "mode: archive or pruned")
	cmd.AddCommand(createCmd)

	modeCmd := &cobra.Command{
		Use:   "mode",
		Short: "Show current node mode",
		Run: func(cmd *cobra.Command, args []string) {
			if fullNode == nil {
				fmt.Println("node not initialised")
				return
			}
			fmt.Println(fullNode.CurrentMode())
		},
	}
	cmd.AddCommand(modeCmd)

	setCmd := &cobra.Command{
		Use:   "set-mode <mode>",
		Short: "Update node mode",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if fullNode == nil {
				fmt.Println("node not initialised")
				return
			}
			if args[0] == "archive" {
				fullNode.SetMode(core.FullNodeModeArchive)
			} else {
				fullNode.SetMode(core.FullNodeModePruned)
			}
			fmt.Println("mode updated")
		},
	}
	cmd.AddCommand(setCmd)

	rootCmd.AddCommand(cmd)
}
