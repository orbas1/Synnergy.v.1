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
			id := benefitRegistry.RegisterBenefit(args[0], args[1], amt)
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}

	claimCmd := &cobra.Command{
		Use:   "claim <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Claim a benefit",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			if err := benefitRegistry.Claim(id); err != nil {
				return err
			}
			cmd.Println("claimed")
			return nil
		},
	}

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

	cmd.AddCommand(registerCmd, claimCmd, getCmd)
	rootCmd.AddCommand(cmd)
}
