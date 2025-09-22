package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn500Token *core.SYN500Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn500",
		Short: "SYN500 utility token",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureStage73Loaded()
		},
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
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("token created")
			return nil
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().Uint("dec", 0, "decimals")
	createCmd.Flags().Uint64("supply", 0, "initial supply")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("owner")
	createCmd.MarkFlagRequired("dec")
	createCmd.MarkFlagRequired("supply")

	grantCmd := &cobra.Command{
		Use:   "grant <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Grant a usage tier",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tier, _ := cmd.Flags().GetInt("tier")
			max, _ := cmd.Flags().GetUint64("max")
			windowStr, _ := cmd.Flags().GetString("window")
			if tier <= 0 || max == 0 {
				return fmt.Errorf("tier and max must be positive")
			}
			window := time.Hour
			if windowStr != "" {
				parsed, err := time.ParseDuration(windowStr)
				if err != nil {
					return fmt.Errorf("invalid window: %w", err)
				}
				window = parsed
			}
			syn500Token.Grant(args[0], tier, max, window)
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("granted")
			return nil
		},
	}
	grantCmd.Flags().Int("tier", 0, "service tier")
	grantCmd.Flags().Uint64("max", 0, "max usage within the window")
	grantCmd.Flags().String("window", "1h", "usage window (e.g. 30m, 1h, 24h)")
	grantCmd.MarkFlagRequired("tier")
	grantCmd.MarkFlagRequired("max")

	useCmd := &cobra.Command{
		Use:   "use <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Record usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			if err := syn500Token.Use(args[0]); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("usage recorded")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show address usage status",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tier, ok := syn500Token.Status(args[0])
			if !ok {
				return fmt.Errorf("address not granted")
			}
			data, _ := json.MarshalIndent(tier, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	telemetryCmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Show aggregate grant telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn500Token == nil {
				return fmt.Errorf("token not created")
			}
			tele := syn500Token.Telemetry()
			data, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	cmd.AddCommand(createCmd, grantCmd, useCmd, statusCmd, telemetryCmd)
	rootCmd.AddCommand(cmd)
}
