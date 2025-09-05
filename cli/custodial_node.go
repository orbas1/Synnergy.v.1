package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var custodialLedger = core.NewLedger()
var custodialNode = core.NewCustodialNode("custodian", "custodian_addr", custodialLedger)

func init() {
	custCmd := &cobra.Command{
		Use:   "custodial",
		Short: "Operate a custodial node",
	}

	var custodyJSON bool
	custodyCmd := &cobra.Command{
		Use:   "custody <user> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Custody assets for a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			custodialNode.Custody(args[0], amt)
			gas := synnergy.GasCost("Custody")
			if custodyJSON {
				enc, _ := json.Marshal(map[string]interface{}{"status": "recorded", "gas": gas})
				fmt.Println(string(enc))
				return nil
			}
			fmt.Println("recorded")
			return nil
		},
	}
	custodyCmd.Flags().BoolVar(&custodyJSON, "json", false, "output as JSON")

	var releaseJSON bool
	releaseCmd := &cobra.Command{
		Use:   "release <user> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Release assets to a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount")
			}
			if err := custodialNode.Release(args[0], amt); err != nil {
				return err
			}
			gas := synnergy.GasCost("Release")
			if releaseJSON {
				enc, _ := json.Marshal(map[string]interface{}{"status": "released", "gas": gas})
				fmt.Println(string(enc))
				return nil
			}
			fmt.Println("released")
			return nil
		},
	}
	releaseCmd.Flags().BoolVar(&releaseJSON, "json", false, "output as JSON")

	var holdingsJSON bool
	holdingsCmd := &cobra.Command{
		Use:   "holdings [user]",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Show holdings",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				bal := custodialNode.Balance(args[0])
				if holdingsJSON {
					enc, _ := json.Marshal(map[string]interface{}{"user": args[0], "balance": bal})
					fmt.Println(string(enc))
					return
				}
				fmt.Println(bal)
				return
			}
			if holdingsJSON {
				out := make(map[string]uint64)
				for u := range custodialNode.Holdings {
					out[u] = custodialNode.Balance(u)
				}
				enc, _ := json.Marshal(out)
				fmt.Println(string(enc))
				return
			}
			for u := range custodialNode.Holdings {
				fmt.Printf("%s: %d\n", u, custodialNode.Balance(u))
			}
		},
	}
	holdingsCmd.Flags().BoolVar(&holdingsJSON, "json", false, "output as JSON")

	custCmd.AddCommand(custodyCmd, releaseCmd, holdingsCmd)
	rootCmd.AddCommand(custCmd)
}
