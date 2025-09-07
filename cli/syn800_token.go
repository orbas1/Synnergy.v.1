package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var assetRegistry = core.NewAssetRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn800_token",
		Short: "Manage SYN800 asset tokens",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <desc> <valuation> <loc> <type> <cert>",
		Args:  cobra.ExactArgs(6),
		Short: "Register an asset",
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid valuation")
			}
			if _, err := assetRegistry.Register(args[0], args[1], val, args[3], args[4], args[5]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "asset registered")
			return nil
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update <id> <valuation>",
		Args:  cobra.ExactArgs(2),
		Short: "Update asset valuation",
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid valuation")
			}
			if err := assetRegistry.UpdateValuation(args[0], val); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "valuation updated")
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Display asset information",
		RunE: func(cmd *cobra.Command, args []string) error {
			if a, ok := assetRegistry.Get(args[0]); ok {
				b, _ := json.MarshalIndent(a, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				return nil
			}
			return fmt.Errorf("asset not found")
		},
	}

	cmd.AddCommand(registerCmd, updateCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
