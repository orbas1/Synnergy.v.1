package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var debtReg = tokens.NewDebtRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn845",
		Short: "Debt token registry",
	}

	createCmd := &cobra.Command{
		Use:   "create-token",
		Short: "Create debt token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			supply, _ := cmd.Flags().GetUint64("supply")
			id, _ := debtReg.CreateToken(name, symbol, owner, supply)
			fmt.Println(id)
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Uint64("supply", 0, "token supply")
	cmd.AddCommand(createCmd)

	issueCmd := &cobra.Command{
		Use:   "issue <tokenID> <debtID> <borrower> <principal> <rate> <penalty> <due>",
		Short: "Issue debt instrument",
		Args:  cobra.ExactArgs(7),
		Run: func(cmd *cobra.Command, args []string) {
			var principal uint64
			var rate, penalty float64
			fmt.Sscanf(args[3], "%d", &principal)
			fmt.Sscanf(args[4], "%f", &rate)
			fmt.Sscanf(args[5], "%f", &penalty)
			due, _ := time.Parse(time.RFC3339, args[6])
			if err := debtReg.IssueDebt(args[0], args[1], args[2], principal, rate, penalty, due); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("issued")
			}
		},
	}
	cmd.AddCommand(issueCmd)

	payCmd := &cobra.Command{
		Use:   "pay <tokenID> <debtID> <amount>",
		Short: "Record debt payment",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := debtReg.RecordPayment(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("payment recorded")
			}
		},
	}
	cmd.AddCommand(payCmd)

	getCmd := &cobra.Command{
		Use:   "get <tokenID> <debtID>",
		Short: "Get debt info",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			d, err := debtReg.GetDebt(args[0], args[1])
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("%+v\n", *d)
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
