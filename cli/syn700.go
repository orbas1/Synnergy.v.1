package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	synn "synnergy"
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
			asset, err := ipRegistry.Register(args[0], args[1], args[2], args[3], args[4])
			if err != nil {
				return err
			}
			printOutput(map[string]any{
				"status": "registered",
				"id":     asset.TokenID,
				"owner":  asset.Owner,
				"gas":    synn.GasCost("RegisterIPAsset"),
			})
			return nil
		},
	}

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <newOwner>",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer ownership",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ipRegistry.Transfer(args[0], args[1]); err != nil {
				return err
			}
			printOutput(map[string]any{
				"status": "transferred",
				"id":     args[0],
				"owner":  args[1],
				"gas":    synn.GasCost("TransferIPAsset"),
			})
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
			expiryStr, _ := cmd.Flags().GetString("expiry")
			if expiryStr == "" {
				return fmt.Errorf("expiry required")
			}
			exp, err := time.Parse(time.RFC3339, expiryStr)
			if err != nil {
				return fmt.Errorf("invalid expiry")
			}
			if err := ipRegistry.CreateLicense(args[0], args[1], args[2], args[3], royalty, exp); err != nil {
				return err
			}
			printOutput(map[string]any{
				"status":  "license created",
				"token":   args[0],
				"license": args[1],
				"expires": exp,
				"gas":     synn.GasCost("CreateIPLicense"),
			})
			return nil
		},
	}
	licenseCmd.Flags().String("expiry", "", "RFC3339 expiry timestamp")
	_ = licenseCmd.MarkFlagRequired("expiry")

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
			printOutput(map[string]any{
				"status":  "royalty recorded",
				"token":   args[0],
				"license": args[1],
				"amount":  amt,
				"gas":     synn.GasCost("RecordIPRoyalty"),
			})
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <tokenID>",
		Args:  cobra.ExactArgs(1),
		Short: "Show token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			if t, ok := ipRegistry.Get(args[0]); ok {
				snap := map[string]any{}
				b, _ := json.Marshal(t)
				_ = json.Unmarshal(b, &snap)
				snap["gas"] = synn.GasCost("DescribeIPAsset")
				printOutput(snap)
				return nil
			}
			return fmt.Errorf("token not found")
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered IP assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			ids := ipRegistry.List()
			printOutput(map[string]any{
				"ids": ids,
				"gas": synn.GasCost("ListIPAssets"),
			})
			return nil
		},
	}

	summaryCmd := &cobra.Command{
		Use:   "royalties <tokenID> <licenseID>",
		Args:  cobra.ExactArgs(2),
		Short: "Summarise royalties for a license",
		RunE: func(cmd *cobra.Command, args []string) error {
			total, err := ipRegistry.RoyaltySummary(args[0], args[1])
			if err != nil {
				return err
			}
			printOutput(map[string]any{
				"token":   args[0],
				"license": args[1],
				"total":   total,
				"gas":     synn.GasCost("SummariseIPRoyalties"),
			})
			return nil
		},
	}

	cmd.AddCommand(registerCmd, transferCmd, licenseCmd, royaltyCmd, infoCmd, listCmd, summaryCmd)
	rootCmd.AddCommand(cmd)
}
