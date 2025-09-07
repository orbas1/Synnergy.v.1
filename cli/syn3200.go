package cli

import (
	"fmt"
	"strconv"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			issuer, _ := cmd.Flags().GetString("issuer")
			payer, _ := cmd.Flags().GetString("payer")
			amt, _ := cmd.Flags().GetUint64("amount")
			dueStr, _ := cmd.Flags().GetString("due")
			meta, _ := cmd.Flags().GetString("meta")
			if id == "" || issuer == "" || payer == "" {
				return fmt.Errorf("id, issuer and payer required")
			}
			due, err := time.Parse(time.RFC3339, dueStr)
			if err != nil {
				return fmt.Errorf("invalid due date: %w", err)
			}
			if _, err := bills.Create(id, issuer, payer, amt, due, meta); err != nil {
				return err
			}
			cmd.Println("bill created")
			return nil
		},
	}
	createCmd.Flags().String("id", "", "bill id")
	createCmd.Flags().String("issuer", "", "issuer")
	createCmd.Flags().String("payer", "", "payer")
	createCmd.Flags().Uint64("amount", 0, "amount")
	createCmd.Flags().String("due", time.Now().Add(24*time.Hour).Format(time.RFC3339), "due time")
	createCmd.Flags().String("meta", "", "metadata")
	_ = createCmd.MarkFlagRequired("id")
	_ = createCmd.MarkFlagRequired("issuer")
	_ = createCmd.MarkFlagRequired("payer")
	_ = createCmd.MarkFlagRequired("amount")
	cmd.AddCommand(createCmd)

	payCmd := &cobra.Command{
		Use:   "pay <id> <payer> <amt>",
		Short: "Record a payment",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			if err := bills.Pay(args[0], args[1], amt); err != nil {
				return err
			}
			cmd.Println("payment recorded")
			return nil
		},
	}
	cmd.AddCommand(payCmd)

	adjustCmd := &cobra.Command{
		Use:   "adjust <id> <amt>",
		Short: "Adjust bill amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			if err := bills.Adjust(args[0], amt); err != nil {
				return err
			}
			cmd.Println("bill adjusted")
			return nil
		},
	}
	cmd.AddCommand(adjustCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get bill info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			b, ok := bills.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			cmd.Printf("%+v\n", *b)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
