package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
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
			gasPrint("SYN500Create")
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
			cmd.Println("token created")
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
		Short: "Grant a usage tier",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SYN500Grant")
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tier, _ := cmd.Flags().GetInt("tier")
			max, _ := cmd.Flags().GetUint64("max")
			window, err := cmd.Flags().GetDuration("window")
			if err != nil {
				return fmt.Errorf("invalid window: %w", err)
			}
			if err := syn500Token.Grant(args[0], tier, max, window); err != nil {
				return err
			}
			cmd.Println("granted")
			return nil
		},
	}
	grantCmd.Flags().Int("tier", 0, "service tier")
	grantCmd.Flags().Uint64("max", 0, "max usage")
	grantCmd.Flags().Duration("window", 0, "rolling usage window (e.g. 1h)")
	_ = grantCmd.MarkFlagRequired("tier")
	_ = grantCmd.MarkFlagRequired("max")

	useCmd := &cobra.Command{
		Use:   "use <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Record usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SYN500Use")
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			if err := syn500Token.Use(args[0]); err != nil {
				return err
			}
			cmd.Println("usage recorded")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant status",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SYN500Status")
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tier, ok := syn500Token.Status(args[0])
			if !ok {
				return fmt.Errorf("grant not found")
			}
			resp := map[string]any{
				"tier": tier.Tier,
				"max":  tier.Max,
				"used": tier.Used,
			}
			if tier.Window > 0 {
				resp["window"] = tier.Window.String()
			}
			if !tier.LastReset.IsZero() {
				resp["last_reset"] = tier.LastReset
			}
			if len(tier.AuditTrail) > 0 {
				resp["audit"] = tier.AuditTrail
			}
			data, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return err
			}
			cmd.Println(string(data))
			return nil
		},
	}

	telemetryCmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Summarise SYN500 usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SYN500Telemetry")
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tele := syn500Token.Telemetry()
			data, err := json.MarshalIndent(tele, "", "  ")
			if err != nil {
				return err
			}
			cmd.Println(string(data))
			return nil
		},
	}

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show SYN500 audit events",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("SYN500Audit")
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			events := syn500Token.AuditLog()
			data, err := json.MarshalIndent(events, "", "  ")
			if err != nil {
				return err
			}
			cmd.Println(string(data))
			return nil
		},
	}

	cmd.AddCommand(createCmd, grantCmd, useCmd, statusCmd, telemetryCmd, auditCmd)
	rootCmd.AddCommand(cmd)
}
