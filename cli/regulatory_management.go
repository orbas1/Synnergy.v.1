package cli

import (
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var regManager = core.NewRegulatoryManager()

func init() {
	cmd := &cobra.Command{
		Use:   "regulator",
		Short: "Manage regulations",
	}

	addCmd := &cobra.Command{
		Use:   "add [id] [jurisdiction] [description] [maxAmount]",
		Args:  cobra.ExactArgs(4),
		Short: "Add a new regulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegulatorAdd")
			amt, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			reg := core.Regulation{ID: args[0], Jurisdiction: args[1], Description: args[2], MaxAmount: amt}
			if err := regManager.AddRegulation(reg); err != nil {
				return err
			}
			printOutput(map[string]any{"status": "added", "id": reg.ID})
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a regulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegulatorRemove")
			regManager.RemoveRegulation(args[0])
			printOutput(map[string]string{"status": "removed", "id": args[0]})
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all regulations",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegulatorList")
			regs := regManager.ListRegulations()
			printOutput(regs)
			return nil
		},
	}

	evalCmd := &cobra.Command{
		Use:   "evaluate [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Check amount against regulations",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegulatorEvaluate")
			amt, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			tx := core.Transaction{Amount: amt}
			v := regManager.EvaluateTransaction(tx)
			if len(v) == 0 {
				printOutput(map[string]string{"status": "ok"})
			} else {
				printOutput(v)
			}
			return nil
		},
	}

	cmd.AddCommand(addCmd, removeCmd, listCmd, evalCmd)
	rootCmd.AddCommand(cmd)
}
