package cli

import (
	"fmt"
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
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			minStake, _ := cmd.Flags().GetUint64("minstake")
			quorum, _ := cmd.Flags().GetInt("quorum")
			validatorNode = core.NewValidatorNode(id, addr, core.NewLedger(), minStake, quorum)
			fmt.Println("validator node created")
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
			if validatorNode == nil {
				fmt.Println("node not initialised")
				return
			}
			stake, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid stake")
				return
			}
			if err := validatorNode.AddValidator(args[0], stake); err != nil {
				fmt.Println("add validator error:", err)
				return
			}
			fmt.Println("validator added")
		},
	}
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <addr>",
		Short: "Remove a validator",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if validatorNode == nil {
				fmt.Println("node not initialised")
				return
			}
			validatorNode.RemoveValidator(args[0])
			fmt.Println("validator removed")
		},
	}
	cmd.AddCommand(removeCmd)

	slashCmd := &cobra.Command{
		Use:   "slash <addr>",
		Short: "Slash a validator",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if validatorNode == nil {
				fmt.Println("node not initialised")
				return
			}
			validatorNode.SlashValidator(args[0])
			fmt.Println("validator slashed")
		},
	}
	cmd.AddCommand(slashCmd)

	quorumCmd := &cobra.Command{
		Use:   "quorum",
		Short: "Check if quorum is reached",
		Run: func(cmd *cobra.Command, args []string) {
			if validatorNode == nil {
				fmt.Println("node not initialised")
				return
			}
			fmt.Println(validatorNode.HasQuorum())
		},
	}
	cmd.AddCommand(quorumCmd)

	rootCmd.AddCommand(cmd)
}
