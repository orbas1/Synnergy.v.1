package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var custodialLedger = core.NewLedger()
var custodialNode = core.NewCustodialNode("custodian", "custodian_addr", custodialLedger)

func init() {
	custCmd := &cobra.Command{
		Use:   "custodial",
		Short: "Operate a custodial node",
	}

	custodyCmd := &cobra.Command{
		Use:   "custody <user> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Custody assets for a user",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			custodialNode.Custody(args[0], amt)
			fmt.Println("recorded")
		},
	}

	releaseCmd := &cobra.Command{
		Use:   "release <user> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Release assets to a user",
		Run: func(cmd *cobra.Command, args []string) {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid amount")
				return
			}
			if err := custodialNode.Release(args[0], amt); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("released")
		},
	}

	holdingsCmd := &cobra.Command{
		Use:   "holdings [user]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Show holdings",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				fmt.Println(custodialNode.Balance(args[0]))
				return
			}
			for u := range custodialNode.Holdings {
				fmt.Printf("%s: %d\n", u, custodialNode.Balance(u))
			}
		},
	}

	custCmd.AddCommand(custodyCmd, releaseCmd, holdingsCmd)
	rootCmd.AddCommand(custCmd)
}
