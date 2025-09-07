package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
	militarynodes "synnergy/internal/nodes/military_nodes"
)

var warfareNode *core.WarfareNode

func init() {
	cmd := &cobra.Command{
		Use:   "warfare",
		Short: "Interact with a warfare node",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create warfare node",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			if id == "" || addr == "" {
				return fmt.Errorf("id and addr required")
			}
			base := core.NewNode(id, addr, core.NewLedger())
			warfareNode = core.NewWarfareNode(base)
			printOutput("warfare node created")
			return nil
		},
	}
	createCmd.Flags().String("id", "", "node id")
	createCmd.Flags().String("addr", "", "node address")
	cmd.AddCommand(createCmd)

	cmdCmd := &cobra.Command{
		Use:   "command <cmd>",
		Args:  cobra.ExactArgs(1),
		Short: "Execute secure command",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return fmt.Errorf("node not initialised")
			}
			if err := warfareNode.SecureCommand(args[0]); err != nil {
				return err
			}
			printOutput("command executed")
			return nil
		},
	}
	cmd.AddCommand(cmdCmd)

	trackCmd := &cobra.Command{
		Use:   "track <assetID> <location> <status>",
		Args:  cobra.ExactArgs(3),
		Short: "Record logistics information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return fmt.Errorf("node not initialised")
			}
			warfareNode.TrackLogistics(args[0], args[1], args[2])
			printOutput("logistics recorded")
			return nil
		},
	}
	cmd.AddCommand(trackCmd)

	listCmd := &cobra.Command{
		Use:   "logistics [assetID]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List logistics records",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return fmt.Errorf("node not initialised")
			}
			var recs []militarynodes.LogisticsRecord
			if len(args) == 1 {
				recs = warfareNode.LogisticsByAsset(args[0])
			} else {
				recs = warfareNode.Logistics()
			}
			printOutput(recs)
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	shareCmd := &cobra.Command{
		Use:   "share <info>",
		Args:  cobra.ExactArgs(1),
		Short: "Share tactical information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if warfareNode == nil {
				return fmt.Errorf("node not initialised")
			}
			warfareNode.ShareTactical(args[0])
			printOutput("info shared")
			return nil
		},
	}
	cmd.AddCommand(shareCmd)

	rootCmd.AddCommand(cmd)
}
