package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
)

var regManager = synnergy.NewRegulatoryManager()

func init() {
	regCmd := &cobra.Command{
		Use:   "regulation",
		Short: "Manage regulatory rules",
	}

	addCmd := &cobra.Command{
		Use:   "add [id] [maxAmount] [jurisdiction] [description]",
		Args:  cobra.ExactArgs(4),
		Short: "Add a regulation",
		RunE: func(cmd *cobra.Command, args []string) error {
			maxAmt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			reg := synnergy.Regulation{ID: args[0], MaxAmount: maxAmt, Jurisdiction: args[2], Description: args[3]}
			return regManager.AddRegulation(reg)
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a regulation",
		Run: func(cmd *cobra.Command, args []string) {
			regManager.RemoveRegulation(args[0])
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get a regulation",
		Run: func(cmd *cobra.Command, args []string) {
			reg, ok := regManager.GetRegulation(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%s: max=%d jurisdiction=%s description=%s\n", reg.ID, reg.MaxAmount, reg.Jurisdiction, reg.Description)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List regulations",
		Run: func(cmd *cobra.Command, args []string) {
			regs := regManager.ListRegulations()
			for _, r := range regs {
				fmt.Printf("%s: max=%d jurisdiction=%s description=%s\n", r.ID, r.MaxAmount, r.Jurisdiction, r.Description)
			}
		},
	}

	evalCmd := &cobra.Command{
		Use:   "eval [from] [to] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Evaluate a transaction",
		Run: func(cmd *cobra.Command, args []string) {
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			tx := synnergy.Transaction{From: args[0], To: args[1], Amount: amount}
			v := regManager.EvaluateTransaction(tx)
			if len(v) == 0 {
				fmt.Println("no violations")
			} else {
				fmt.Println("violations:", v)
			}
		},
	}

	regCmd.AddCommand(addCmd, removeCmd, getCmd, listCmd, evalCmd)
	rootCmd.AddCommand(regCmd)
}
