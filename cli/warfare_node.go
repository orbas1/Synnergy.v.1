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
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			base := core.NewNode(id, addr, core.NewLedger())
			warfareNode = core.NewWarfareNode(base)
			fmt.Println("warfare node created")
		},
	}
	createCmd.Flags().String("id", "", "node id")
	createCmd.Flags().String("addr", "", "node address")
	cmd.AddCommand(createCmd)

	cmdCmd := &cobra.Command{
		Use:   "command <cmd>",
		Args:  cobra.ExactArgs(1),
		Short: "Execute secure command",
		Run: func(cmd *cobra.Command, args []string) {
			if warfareNode == nil {
				fmt.Println("node not initialised")
				return
			}
			if err := warfareNode.SecureCommand(args[0]); err != nil {
				fmt.Println("command error:", err)
				return
			}
			fmt.Println("command executed")
		},
	}
	cmd.AddCommand(cmdCmd)

	trackCmd := &cobra.Command{
		Use:   "track <assetID> <location> <status>",
		Args:  cobra.ExactArgs(3),
		Short: "Record logistics information",
		Run: func(cmd *cobra.Command, args []string) {
			if warfareNode == nil {
				fmt.Println("node not initialised")
				return
			}
			warfareNode.TrackLogistics(args[0], args[1], args[2])
			fmt.Println("logistics recorded")
		},
	}
	cmd.AddCommand(trackCmd)

	listCmd := &cobra.Command{
		Use:   "logistics [assetID]",
		Args:  cobra.MaximumNArgs(1),
		Short: "List logistics records",
		Run: func(cmd *cobra.Command, args []string) {
			if warfareNode == nil {
				fmt.Println("node not initialised")
				return
			}
			var recs []militarynodes.LogisticsRecord
			if len(args) == 1 {
				recs = warfareNode.LogisticsByAsset(args[0])
			} else {
				recs = warfareNode.Logistics()
			}
			for i, r := range recs {
				fmt.Printf("%d: %+v\n", i, r)
			}
		},
	}
	cmd.AddCommand(listCmd)

	shareCmd := &cobra.Command{
		Use:   "share <info>",
		Args:  cobra.ExactArgs(1),
		Short: "Share tactical information",
		Run: func(cmd *cobra.Command, args []string) {
			if warfareNode == nil {
				fmt.Println("node not initialised")
				return
			}
			warfareNode.ShareTactical(args[0])
			fmt.Println("info shared")
		},
	}
	cmd.AddCommand(shareCmd)

	rootCmd.AddCommand(cmd)
}
