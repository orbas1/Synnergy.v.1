package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var ipRegistry = core.NewIPRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn700",
		Short: "Manage SYN700 IP tokens",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <title> <desc> <creator> <owner>",
		Args:  cobra.ExactArgs(5),
		Short: "Register an IP asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := ipRegistry.Register(args[0], args[1], args[2], args[3], args[4]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "registered")
			return nil
		},
	}

	licenseCmd := &cobra.Command{
		Use:   "license <tokenID> <licID> <type> <licensee> <royalty>",
		Args:  cobra.ExactArgs(5),
		Short: "Create a license",
		RunE: func(cmd *cobra.Command, args []string) error {
			royalty, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid royalty")
			}
			if err := ipRegistry.CreateLicense(args[0], args[1], args[2], args[3], royalty); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "license created")
			return nil
		},
	}

	royaltyCmd := &cobra.Command{
		Use:   "royalty <tokenID> <licID> <licensee> <amount>",
		Args:  cobra.ExactArgs(4),
		Short: "Record a royalty payment",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			if err := ipRegistry.RecordRoyalty(args[0], args[1], args[2], amt); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "royalty recorded")
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <tokenID>",
		Args:  cobra.ExactArgs(1),
		Short: "Show token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			if t, ok := ipRegistry.Get(args[0]); ok {
				b, _ := json.MarshalIndent(t, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				return nil
			}
			return fmt.Errorf("token not found")
		},
	}

	cmd.AddCommand(registerCmd, licenseCmd, royaltyCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
