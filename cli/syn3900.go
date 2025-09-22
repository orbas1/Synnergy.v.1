package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "syn3900",
		Short: "Manage SYN3900 government benefits",
	}

	registerCmd := &cobra.Command{
		Use:   "register <recipient> <program> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a new benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3900Register")
			store, err := stage73State()
			if err != nil {
				return err
			}
			if args[0] == "" {
				return fmt.Errorf("recipient required")
			}
			if args[1] == "" {
				return fmt.Errorf("program required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			registry := store.Benefits()
			id := registry.RegisterBenefit(args[0], args[1], amt)
			approver, _ := cmd.Flags().GetString("approver")
			if approver != "" {
				path, password, err := parseCredentialPair(approver)
				if err != nil {
					return err
				}
				wallet, err := loadWallet(path, password)
				if err != nil {
					return err
				}
				if err := registry.AddApprover(id, wallet.Address); err != nil {
					return err
				}
			}
			markStage73Dirty()
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}
	registerCmd.Flags().String("approver", "", "optional approver credential path:password")
	cmd.AddCommand(registerCmd)

	claimCmd := &cobra.Command{
		Use:   "claim <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Claim a benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3900Claim")
			store, err := stage73State()
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			registry := store.Benefits()
			if err := registry.Claim(id, wallet.Address); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("claimed")
			return nil
		},
	}
	claimCmd.Flags().String("wallet", "", "recipient wallet path")
	claimCmd.Flags().String("password", "", "wallet password")
	cmd.AddCommand(claimCmd)

	approveCmd := &cobra.Command{
		Use:   "approve <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Approve a claimed benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3900Approve")
			store, err := stage73State()
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			registry := store.Benefits()
			_ = registry.AddApprover(id, wallet.Address)
			if err := registry.Approve(id, wallet.Address); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("approved")
			return nil
		},
	}
	approveCmd.Flags().String("wallet", "", "approver wallet path")
	approveCmd.Flags().String("password", "", "wallet password")
	cmd.AddCommand(approveCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show benefit details",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3900Get")
			store, err := stage73State()
			if err != nil {
				return err
			}
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			b, ok := store.Benefits().GetBenefit(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			data, _ := json.MarshalIndent(b, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered benefits",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3900List")
			store, err := stage73State()
			if err != nil {
				return err
			}
			list := store.Benefits().ListBenefits()
			data, _ := json.MarshalIndent(list, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show benefit telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3900Status")
			store, err := stage73State()
			if err != nil {
				return err
			}
			summary := store.Benefits().Summary()
			data, _ := json.MarshalIndent(summary, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	rootCmd.AddCommand(cmd)
}
