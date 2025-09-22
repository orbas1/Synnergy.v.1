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
			return ensureBenefitRegistryLoaded()
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
			approver := ""
			if approverSpec != "" {
				addr, err := walletAddressFromSpec(approverSpec)
				if err != nil {
					return err
				}
				approver = addr
			}
			id, err := benefitRegistry.RegisterBenefit(args[0], args[1], amt, core.Address(approver))
			if err != nil {
				return err
			}
			if err := persistBenefitRegistry(); err != nil {
				return err
			}
			cmd.Println(id)
			return nil
		},
	}
	registerCmd.Flags().String("approver", "", "approver wallet path:password")

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
			if err := benefitRegistry.Claim(id, core.Address(wallet.Address)); err != nil {
				return err
			}
			if err := persistBenefitRegistry(); err != nil {
				return err
			}
			printOutput("claimed")
			return nil
		},
	}
	addWalletFlags(claimCmd)

	approveCmd := &cobra.Command{
		Use:   "approve <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Approve a benefit payout",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := benefitRegistry.Approve(id, core.Address(wallet.Address)); err != nil {
				return err
			}
			if err := persistBenefitRegistry(); err != nil {
				return err
			}
			printOutput("approved")
			return nil
		},
	}
	addWalletFlags(approveCmd)

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
		Short: "List registered benefits",
		Run: func(cmd *cobra.Command, args []string) {
			bs := benefitRegistry.ListBenefits()
			data, _ := json.MarshalIndent(bs, "", "  ")
			cmd.Println(string(data))
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show benefit registry status",
		Run: func(cmd *cobra.Command, args []string) {
			summary := benefitRegistry.StatusSummary()
			data, _ := json.MarshalIndent(summary, "", "  ")
			cmd.Println(string(data))
		},
	}

	cmd.AddCommand(registerCmd, claimCmd, approveCmd, getCmd, listCmd, statusCmd)
	rootCmd.AddCommand(cmd)
}
