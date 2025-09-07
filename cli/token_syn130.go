package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	syn130Registry = core.NewTangibleAssetRegistry()
	syn130IDs      []string
)

func init() {
	cmd := &cobra.Command{
		Use:   "token_syn130",
		Short: "Manage SYN130 tangible assets",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <owner> <meta> <value>",
		Args:  cobra.ExactArgs(4),
		Short: "Register a tangible asset",
		Run: func(cmd *cobra.Command, args []string) {
			val, _ := strconv.ParseUint(args[3], 10, 64)
			if _, err := syn130Registry.Register(args[0], args[1], args[2], val); err != nil {
				fmt.Println("error:", err)
			} else {
				syn130IDs = append(syn130IDs, args[0])
			}
		},
	}

	valueCmd := &cobra.Command{
		Use:   "value <id> <val>",
		Args:  cobra.ExactArgs(2),
		Short: "Update valuation",
		Run: func(cmd *cobra.Command, args []string) {
			val, _ := strconv.ParseUint(args[1], 10, 64)
			if err := syn130Registry.UpdateValuation(args[0], val); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	saleCmd := &cobra.Command{
		Use:   "sale <id> <buyer> <price>",
		Args:  cobra.ExactArgs(3),
		Short: "Record a sale",
		Run: func(cmd *cobra.Command, args []string) {
			price, _ := strconv.ParseUint(args[2], 10, 64)
			if err := syn130Registry.RecordSale(args[0], args[1], price); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	leaseCmd := &cobra.Command{
		Use:   "lease <id> <lessee> <pay> <start> <end>",
		Args:  cobra.ExactArgs(5),
		Short: "Start a lease",
		Run: func(cmd *cobra.Command, args []string) {
			pay, _ := strconv.ParseUint(args[2], 10, 64)
			start, _ := strconv.ParseInt(args[3], 10, 64)
			end, _ := strconv.ParseInt(args[4], 10, 64)
			if err := syn130Registry.StartLease(args[0], args[1], pay, time.Unix(start, 0), time.Unix(end, 0)); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	endLeaseCmd := &cobra.Command{
		Use:   "endlease <id>",
		Args:  cobra.ExactArgs(1),
		Short: "End a lease",
		Run: func(cmd *cobra.Command, args []string) {
			if err := syn130Registry.EndLease(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get asset info",
		Run: func(cmd *cobra.Command, args []string) {
			if a, ok := syn130Registry.Get(args[0]); ok {
				b, _ := json.MarshalIndent(a, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all tangible assets",
		Run: func(cmd *cobra.Command, args []string) {
			var assets []*core.TangibleAsset
			for _, id := range syn130IDs {
				if a, ok := syn130Registry.Get(id); ok {
					assets = append(assets, a)
				}
			}
			b, _ := json.MarshalIndent(assets, "", "  ")
			fmt.Println(string(b))
		},
	}

	cmd.AddCommand(registerCmd, valueCmd, saleCmd, leaseCmd, endLeaseCmd, infoCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
