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
		Short: "Manage SYN3900 benefit records",
	}

	registerCmd := &cobra.Command{
		Use:   "register <recipient> <program> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a new benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			recipient := args[0]
			if recipient == "" {
				return fmt.Errorf("recipient required")
			}
			program := args[1]
			if program == "" {
				return fmt.Errorf("program required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			approverFlag, _ := cmd.Flags().GetString("approver")
			if approverFlag != "" {
				path, password, err := parseWalletCredential(approverFlag)
				if err != nil {
					return err
				}
				if _, err := loadWallet(path, password); err != nil {
					return err
				}
			}
			id := benefitRegistry.RegisterBenefit(recipient, program, amt)
			cmd.Println(id)
			return nil
		},
	}
	registerCmd.Flags().String("approver", "", "wallet credentials path:password to validate registration")

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
			cmd.Println("claimed")
			return nil
		},
	}
	claimCmd.Flags().String("wallet", "", "path to wallet file")
	claimCmd.Flags().String("password", "", "wallet password")

	approveCmd := &cobra.Command{
		Use:   "approve <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Approve a claimed benefit",
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
			cmd.Println("approved")
			return nil
		},
	}
	approveCmd.Flags().String("wallet", "", "path to wallet file")
	approveCmd.Flags().String("password", "", "wallet password")

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
			payload, _ := json.MarshalIndent(b, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List benefits",
		Run: func(cmd *cobra.Command, args []string) {
			bs := benefitRegistry.ListBenefits()
			payload, _ := json.MarshalIndent(bs, "", "  ")
			cmd.Println(string(payload))
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show benefit telemetry",
		Run: func(cmd *cobra.Command, args []string) {
			tele := benefitRegistry.Telemetry()
			payload, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(payload))
		},
	}

	cmd.AddCommand(registerCmd, claimCmd, approveCmd, getCmd, listCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
