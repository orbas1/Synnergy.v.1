package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var benefitRegistry = core.NewBenefitRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3900",
		Short: "Manage SYN3900 government benefits",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureStage73Loaded()
		},
	}

	registerCmd := &cobra.Command{
		Use:   "register <recipient> <program> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a new benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" || args[1] == "" {
				return fmt.Errorf("recipient and program required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			approverSpec, _ := cmd.Flags().GetString("approver")
			var approver string
			if approverSpec != "" {
				wallet, err := loadWalletSpec(approverSpec)
				if err != nil {
					return err
				}
				approver = wallet.Address
			}
			id := benefitRegistry.RegisterBenefit(args[0], args[1], amt, approver)
			if err := persistStage73(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}
	registerCmd.Flags().String("approver", "", "Wallet path:password used to pre-approve the benefit")

	claimCmd := &cobra.Command{
		Use:   "claim <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Claim a benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := benefitRegistry.Claim(id, wallet.Address); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("claimed")
			return nil
		},
	}
	claimCmd.Flags().String("wallet", "", "Claimant wallet path")
	claimCmd.Flags().String("password", "", "Claimant wallet password")

	approveCmd := &cobra.Command{
		Use:   "approve <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Approve a benefit payout",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := benefitRegistry.Approve(id, wallet.Address); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("approved")
			return nil
		},
	}
	approveCmd.Flags().String("wallet", "", "Approver wallet path")
	approveCmd.Flags().String("password", "", "Approver wallet password")

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show benefit details",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			b, ok := benefitRegistry.GetBenefit(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			data, _ := json.MarshalIndent(b, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List benefits",
		RunE: func(cmd *cobra.Command, args []string) error {
			bs := benefitRegistry.ListBenefits()
			data, _ := json.MarshalIndent(bs, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show benefit telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			tele := benefitRegistry.Telemetry()
			data, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	cmd.AddCommand(registerCmd, claimCmd, approveCmd, getCmd, listCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
