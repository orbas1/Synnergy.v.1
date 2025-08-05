package cli

import (
	"fmt"
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
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[3], 10, 64)
			reg := core.Regulation{ID: args[0], Jurisdiction: args[1], Description: args[2], MaxAmount: amt}
			if err := regManager.AddRegulation(reg); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a regulation",
		Run:   func(cmd *cobra.Command, args []string) { regManager.RemoveRegulation(args[0]) },
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all regulations",
		Run: func(cmd *cobra.Command, args []string) {
			for _, r := range regManager.ListRegulations() {
				fmt.Printf("%s %s %d\n", r.ID, r.Jurisdiction, r.MaxAmount)
			}
		},
	}

	evalCmd := &cobra.Command{
		Use:   "evaluate [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Check amount against regulations",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[0], 10, 64)
			tx := core.Transaction{Amount: amt}
			v := regManager.EvaluateTransaction(tx)
			if len(v) == 0 {
				fmt.Println("ok")
				return
			}
			fmt.Println(v)
		},
	}

	cmd.AddCommand(addCmd, removeCmd, listCmd, evalCmd)
	rootCmd.AddCommand(cmd)
}
