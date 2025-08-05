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
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			id := grantRegistry.CreateGrant(args[0], args[1], amt)
			fmt.Println("grant created", id)
		},
	}

	releaseCmd := &cobra.Command{
		Use:   "release <id> <amount> [note]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Release funds for a grant",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			note := ""
			if len(args) == 3 {
				note = args[2]
			}
			if err := grantRegistry.Disburse(id, amt, note); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant details",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := strconv.ParseUint(args[0], 10, 64)
			if g, ok := grantRegistry.GetGrant(id); ok {
				b, _ := json.MarshalIndent(g, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List grants",
		Run: func(cmd *cobra.Command, args []string) {
			gs := grantRegistry.ListGrants()
			b, _ := json.MarshalIndent(gs, "", "  ")
			fmt.Println(string(b))
		},
	}

	cmd.AddCommand(createCmd, releaseCmd, getCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
