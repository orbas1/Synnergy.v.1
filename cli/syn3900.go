package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var benefitRegistry = core.NewBenefitRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3900",
		Short: "Manage SYN3900 benefits",
	}

	registerCmd := &cobra.Command{
		Use:   "register <recipient> <program> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a new benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			recipient := args[0]
			program := args[1]
			if recipient == "" {
				return fmt.Errorf("recipient required")
			}
			if program == "" {
				return fmt.Errorf("program required")
			}
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amount == 0 {
				return fmt.Errorf("invalid amount")
			}
			approverSpec, _ := cmd.Flags().GetString("approver")
			var approverAddr string
			if approverSpec != "" {
				path, password, err := parseWalletDescriptor(approverSpec)
				if err != nil {
					return err
				}
				wallet, err := loadWallet(path, password)
				if err != nil {
					return err
				}
				approverAddr = wallet.Address
			}
			id, err := benefitRegistry.RegisterBenefitWithApprover(recipient, program, amount, approverAddr)
			if err != nil {
				return err
			}
			gasPrint("RegisterBenefit")
			cmd.Println(id)
			return nil
		},
	}
	registerCmd.Flags().String("approver", "", "wallet descriptor path:password for initial approval")
	cmd.AddCommand(registerCmd)

	claimCmd := &cobra.Command{
		Use:   "claim <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Claim a benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := benefitRegistry.ClaimWithWallet(id, wallet.Address); err != nil {
				return err
			}
			gasPrint("ClaimBenefit")
			cmd.Println("claimed")
			return nil
		},
	}
	claimCmd.Flags().String("wallet", "", "wallet path for recipient or approver")
	claimCmd.Flags().String("password", "", "wallet password")
	cmd.AddCommand(claimCmd)

	approveCmd := &cobra.Command{
		Use:   "approve <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Approve an additional wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := benefitRegistry.Approve(id, wallet.Address); err != nil {
				return err
			}
			gasPrint("ApproveBenefitWallet")
			cmd.Println("approved")
			return nil
		},
	}
	approveCmd.Flags().String("wallet", "", "wallet path")
	approveCmd.Flags().String("password", "", "wallet password")
	cmd.AddCommand(approveCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show benefit details",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			benefit, ok := benefitRegistry.GetBenefit(id)
			if !ok {
				return fmt.Errorf("benefit not found")
			}
			return printJSON(cmd, benefit)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List benefits",
		RunE: func(cmd *cobra.Command, args []string) error {
			benefits := benefitRegistry.ListBenefits()
			return printJSON(cmd, benefits)
		},
	}
	cmd.AddCommand(listCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Summarise benefit lifecycle counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			counts := benefitRegistry.StatusSummary()
			out := struct {
				Total    int `json:"total"`
				Pending  int `json:"pending"`
				Approved int `json:"approved"`
				Claimed  int `json:"claimed"`
			}{
				Total:    counts["total"],
				Pending:  counts[string(core.BenefitStatusPending)],
				Approved: counts[string(core.BenefitStatusApproved)],
				Claimed:  counts[string(core.BenefitStatusClaimed)],
			}
			return printJSON(cmd, out)
		},
	}
	cmd.AddCommand(statusCmd)

	auditCmd := &cobra.Command{
		Use:   "audit <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show benefit audit events",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			events, err := benefitRegistry.AuditTrail(id)
			if err != nil {
				return err
			}
			return printJSON(cmd, events)
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}
