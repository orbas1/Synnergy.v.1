package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var bills = core.NewBillRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3200",
		Short: "Bill registry operations",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a bill",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			issuer, _ := cmd.Flags().GetString("issuer")
			payer, _ := cmd.Flags().GetString("payer")
			amt, _ := cmd.Flags().GetUint64("amount")
			dueStr, _ := cmd.Flags().GetString("due")
			meta, _ := cmd.Flags().GetString("meta")
			due, _ := time.Parse(time.RFC3339, dueStr)
			if _, err := bills.Create(id, issuer, payer, amt, due, meta); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("bill created")
			}
		},
	}
	createCmd.Flags().String("id", "", "bill id")
	createCmd.Flags().String("issuer", "", "issuer")
	createCmd.Flags().String("payer", "", "payer")
	createCmd.Flags().Uint64("amount", 0, "amount")
	createCmd.Flags().String("due", time.Now().Add(24*time.Hour).Format(time.RFC3339), "due time")
	createCmd.Flags().String("meta", "", "metadata")
	cmd.AddCommand(createCmd)

	payCmd := &cobra.Command{
		Use:   "pay <id> <payer> <amt>",
		Short: "Record a payment",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := bills.Pay(args[0], args[1], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("payment recorded")
			}
		},
	}
	cmd.AddCommand(payCmd)

	adjustCmd := &cobra.Command{
		Use:   "adjust <id> <amt>",
		Short: "Adjust bill amount",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var amt uint64
			fmt.Sscanf(args[1], "%d", &amt)
			if err := bills.Adjust(args[0], amt); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("bill adjusted")
			}
		},
	}
	cmd.AddCommand(adjustCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get bill info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			b, ok := bills.Get(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%+v\n", *b)
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
