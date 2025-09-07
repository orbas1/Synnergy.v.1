package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var validatorNode *core.ValidatorNode

func init() {
	cmd := &cobra.Command{
		Use:   "validatornode",
		Short: "Operations for validator nodes",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a validator node",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ValidatorNodeCreate")
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			minStake, _ := cmd.Flags().GetUint64("minstake")
			quorum, _ := cmd.Flags().GetInt("quorum")
			validatorNode = core.NewValidatorNode(id, addr, core.NewLedger(), minStake, quorum)
			printOutput(map[string]any{"status": "created", "id": id})
		},
	}
	createCmd.Flags().String("id", "", "node id")
	createCmd.Flags().String("addr", "", "node address")
	createCmd.Flags().Uint64("minstake", 0, "minimum stake")
	createCmd.Flags().Int("quorum", 1, "quorum requirement")
	cmd.AddCommand(createCmd)

	addCmd := &cobra.Command{
		Use:   "add <addr> <stake>",
		Short: "Add a validator",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ValidatorNodeAdd")
			if validatorNode == nil {
				printOutput(map[string]any{"error": "node not initialised"})
				return
			}
			stake, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				printOutput(map[string]any{"error": "invalid stake"})
				return
			}
			if err := validatorNode.AddValidator(args[0], stake); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]any{"status": "added", "address": args[0]})
		},
	}
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <addr>",
		Short: "Remove a validator",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ValidatorNodeRemove")
			if validatorNode == nil {
				printOutput(map[string]any{"error": "node not initialised"})
				return
			}
			validatorNode.RemoveValidator(args[0])
			printOutput(map[string]any{"status": "removed", "address": args[0]})
		},
	}
	cmd.AddCommand(removeCmd)

	slashCmd := &cobra.Command{
		Use:   "slash <addr>",
		Short: "Slash a validator",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ValidatorNodeSlash")
			if validatorNode == nil {
				printOutput(map[string]any{"error": "node not initialised"})
				return
			}
			validatorNode.SlashValidator(args[0])
			printOutput(map[string]any{"status": "slashed", "address": args[0]})
		},
	}
	cmd.AddCommand(slashCmd)

	quorumCmd := &cobra.Command{
		Use:   "quorum",
		Short: "Check if quorum is reached",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ValidatorNodeQuorum")
			if validatorNode == nil {
				printOutput(map[string]any{"error": "node not initialised"})
				return
			}
			printOutput(map[string]bool{"quorum": validatorNode.HasQuorum()})
		},
	}
	cmd.AddCommand(quorumCmd)

	rootCmd.AddCommand(cmd)
}
