package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	cmd := &cobra.Command{
		Use:   "syn500",
		Short: "SYN500 utility token",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SYN500 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn500Create")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			dec, _ := cmd.Flags().GetUint("dec")
			supply, _ := cmd.Flags().GetUint64("supply")
			if name == "" {
				return fmt.Errorf("name required")
			}
			if symbol == "" {
				return fmt.Errorf("symbol required")
			}
			if owner == "" {
				return fmt.Errorf("owner required")
			}
			if dec == 0 {
				return fmt.Errorf("decimals must be positive")
			}
			if supply == 0 {
				return fmt.Errorf("supply must be positive")
			}
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := core.NewSYN500Token(name, symbol, owner, uint8(dec), supply)
			store.SetUtility(token)
			markStage73Dirty()
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
			gasPrint("Syn500Grant")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Utility()
			if token == nil {
				return fmt.Errorf("token not created")
			}
			tier, _ := cmd.Flags().GetInt("tier")
			max, _ := cmd.Flags().GetUint64("max")
			windowStr, _ := cmd.Flags().GetString("window")
			if tier <= 0 {
				return fmt.Errorf("tier must be positive")
			}
			if max == 0 {
				return fmt.Errorf("max must be positive")
			}
			var window time.Duration
			if windowStr != "" {
				d, err := time.ParseDuration(windowStr)
				if err != nil {
					return fmt.Errorf("invalid window duration")
				}
				window = d
			}
			token.Grant(args[0], tier, max, window)
			markStage73Dirty()
			cmd.Println("granted")
			return nil
		},
	}
	grantCmd.Flags().Int("tier", 0, "service tier")
	grantCmd.Flags().Uint64("max", 0, "max usage")
	grantCmd.Flags().String("window", "1h", "usage reset window")
	_ = grantCmd.MarkFlagRequired("tier")
	_ = grantCmd.MarkFlagRequired("max")

	useCmd := &cobra.Command{
		Use:   "use <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Record usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn500Use")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Utility()
			if token == nil {
				return fmt.Errorf("token not created")
			}
			if err := token.Use(args[0]); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("usage recorded")
			return nil
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Show grant usage status",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn500Status")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Utility()
			if token == nil {
				return fmt.Errorf("token not created")
			}
			st, ok := token.Status(args[0])
			if !ok {
				return fmt.Errorf("no tier granted")
			}
			data, _ := json.MarshalIndent(st, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	telemetryCmd := &cobra.Command{
		Use:   "telemetry",
		Short: "Show aggregate usage telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn500Telemetry")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Utility()
			if token == nil {
				return fmt.Errorf("token not created")
			}
			tele := token.Telemetry()
			data, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}

	cmd.AddCommand(createCmd, grantCmd, useCmd, statusCmd, telemetryCmd)
	rootCmd.AddCommand(cmd)
}
