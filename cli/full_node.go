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
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			modeStr, _ := cmd.Flags().GetString("mode")
			var mode core.FullNodeMode
			switch modeStr {
			case "archive":
				mode = core.FullNodeModeArchive
			case "pruned":
				mode = core.FullNodeModePruned
			default:
				return fmt.Errorf("unknown mode %s", modeStr)
			}
			fullNode = core.NewFullNode(nodes.Address(id), mode)
			printOutput(map[string]any{"status": "created", "mode": modeStr})
			return nil
		},
	}
	createCmd.Flags().String("id", "", "node id")
	createCmd.Flags().String("mode", "archive", "mode: archive or pruned")
	cmd.AddCommand(createCmd)

	modeCmd := &cobra.Command{
		Use:   "mode",
		Short: "Show current node mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fullNode == nil {
				return fmt.Errorf("node not initialised")
			}
			printOutput(fullNode.CurrentMode())
			return nil
		},
	}
	cmd.AddCommand(modeCmd)

	setCmd := &cobra.Command{
		Use:   "set-mode <mode>",
		Short: "Update node mode",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if fullNode == nil {
				return fmt.Errorf("node not initialised")
			}
			switch args[0] {
			case "archive":
				fullNode.SetMode(core.FullNodeModeArchive)
			case "pruned":
				fullNode.SetMode(core.FullNodeModePruned)
			default:
				return fmt.Errorf("unknown mode %s", args[0])
			}
			printOutput(map[string]string{"status": "mode updated", "mode": args[0]})
			return nil
		},
	}
	cmd.AddCommand(setCmd)

	rootCmd.AddCommand(cmd)
}
