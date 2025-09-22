package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	synn "synnergy"
	"synnergy/core"
)

var syn500Token *core.SYN500Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn500",
		Short: "SYN500 utility token",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SYN500 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			dec, _ := cmd.Flags().GetUint("dec")
			supply, _ := cmd.Flags().GetUint64("supply")
			if name == "" || symbol == "" || owner == "" {
				return fmt.Errorf("name, symbol and owner required")
			}
			if dec == 0 || supply == 0 {
				return fmt.Errorf("decimals and supply must be positive")
			}
			syn500Token = core.NewSYN500Token(name, symbol, owner, uint8(dec), supply)
			printOutput(map[string]any{
				"status": "token created",
				"symbol": symbol,
				"gas":    synn.GasCost("SYN500Create"),
			})
			return nil
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Uint("dec", 0, "decimals")
	createCmd.Flags().Uint64("supply", 0, "initial supply")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("symbol")
	_ = createCmd.MarkFlagRequired("owner")
	_ = createCmd.MarkFlagRequired("dec")
	_ = createCmd.MarkFlagRequired("supply")

	grantCmd := &cobra.Command{
		Use:   "grant <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Grant or update a usage tier",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tier, _ := cmd.Flags().GetInt("tier")
			max, _ := cmd.Flags().GetUint64("max")
			if err := syn500Token.Grant(args[0], tier, max); err != nil {
				return err
			}
			usage, _ := syn500Token.Usage(args[0])
			printOutput(map[string]any{
				"status":    "granted",
				"address":   args[0],
				"tier":      tier,
				"max":       max,
				"gas":       synn.GasCost("GrantServiceTier"),
				"remaining": usage.Max - usage.Used,
			})
			return nil
		},
	}
	grantCmd.Flags().Int("tier", 0, "service tier")
	grantCmd.Flags().Uint64("max", 0, "max usage units")
	_ = grantCmd.MarkFlagRequired("tier")
	_ = grantCmd.MarkFlagRequired("max")

	revokeCmd := &cobra.Command{
		Use:   "revoke <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Revoke a usage grant",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			syn500Token.Revoke(args[0])
			printOutput(map[string]any{
				"status":  "revoked",
				"address": args[0],
				"gas":     synn.GasCost("RevokeServiceTier"),
			})
			return nil
		},
	}

	useCmd := &cobra.Command{
		Use:   "use <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Record usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			amount, _ := cmd.Flags().GetUint64("amount")
			note, _ := cmd.Flags().GetString("note")
			record, err := syn500Token.Use(args[0], amount, note)
			if err != nil {
				return err
			}
			printOutput(map[string]any{
				"status":    "usage recorded",
				"address":   record.Address,
				"amount":    record.Amount,
				"remaining": record.Remaining,
				"digest":    record.Digest,
				"gas":       synn.GasCost("RecordServiceUsage"),
			})
			return nil
		},
	}
	useCmd.Flags().Uint64("amount", 1, "usage units")
	useCmd.Flags().String("note", "", "audit note")

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Display grant snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			snap := syn500Token.Snapshot()
			snap["gas"] = synn.GasCost("SYN500Snapshot")
			printOutput(snap)
			return nil
		},
	}

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Display recent usage entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			limit, _ := cmd.Flags().GetInt("limit")
			entries := syn500Token.AuditTrail(limit)
			printOutput(map[string]any{
				"entries": entries,
				"gas":     synn.GasCost("SYN500Audit"),
			})
			return nil
		},
	}
	auditCmd.Flags().Int("limit", 5, "number of entries to return")

	cmd.AddCommand(createCmd, grantCmd, revokeCmd, useCmd, snapshotCmd, auditCmd)
	rootCmd.AddCommand(cmd)
}
