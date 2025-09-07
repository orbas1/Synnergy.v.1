package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var grantRegistry = core.NewGrantRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3800",
		Short: "Manage SYN3800 grant records",
	}

	createCmd := &cobra.Command{
		Use:   "create <beneficiary> <name> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Create a new grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" || args[1] == "" {
				return fmt.Errorf("beneficiary and name required")
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			id := grantRegistry.CreateGrant(args[0], args[1], amt)
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}

	releaseCmd := &cobra.Command{
		Use:   "release <id> <amount> [note]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Release funds for a grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil || amt == 0 {
				return fmt.Errorf("invalid amount")
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			if err := grantRegistry.Disburse(id, amt, note); err != nil {
				return err
			}
			cmd.Println("released")
			return nil
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant details",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			g, ok := grantRegistry.GetGrant(id)
			if !ok {
				return fmt.Errorf("not found")
			}
			b, _ := json.MarshalIndent(g, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		Run: func(cmd *cobra.Command, args []string) {
			gs := grantRegistry.ListGrants()
			b, _ := json.MarshalIndent(gs, "", "  ")
			cmd.Println(string(b))
		},
	}

	cmd.AddCommand(createCmd, releaseCmd, getCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
