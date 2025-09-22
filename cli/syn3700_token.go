package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn3700 *core.SYN3700Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn3700",
		Short: "SYN3700 index token",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureStage73Loaded()
		},
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the index token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			controllerSpec, _ := cmd.Flags().GetString("controller")
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			syn3700 = core.NewSYN3700Token(name, symbol)
			if controllerSpec != "" {
				wallet, err := loadWalletSpec(controllerSpec)
				if err != nil {
					return err
				}
				syn3700.AddController(wallet.Address)
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("controller", "", "initial controller wallet path:password")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("symbol")
	cmd.AddCommand(initCmd)

	addCmd := &cobra.Command{
		Use:   "add <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Add component to index",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			weight, err := cmd.Flags().GetFloat64("weight")
			if err != nil || weight <= 0 {
				return fmt.Errorf("invalid weight")
			}
			drift, err := cmd.Flags().GetFloat64("drift")
			if err != nil {
				return fmt.Errorf("invalid drift")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := syn3700.AddComponent(args[0], weight, drift, wallet.Address); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("component added")
			return nil
		},
	}
	addCmd.Flags().Float64("weight", 0, "component weight")
	addCmd.Flags().Float64("drift", 0.1, "allowed drift ratio (0-1)")
	addCmd.Flags().String("wallet", "", "controller wallet path")
	addCmd.Flags().String("password", "", "controller wallet password")
	addCmd.MarkFlagRequired("weight")
	addCmd.MarkFlagRequired("wallet")
	addCmd.MarkFlagRequired("password")
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove component from index",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := syn3700.RemoveComponent(args[0], wallet.Address); err != nil {
				return err
			}
			if err := persistStage73(); err != nil {
				return err
			}
			cmd.Println("component removed")
			return nil
		},
	}
	removeCmd.Flags().String("wallet", "", "controller wallet path")
	removeCmd.Flags().String("password", "", "controller wallet password")
	removeCmd.MarkFlagRequired("wallet")
	removeCmd.MarkFlagRequired("password")
	cmd.AddCommand(removeCmd)

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Show index snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			snap := syn3700.Snapshot()
			data, _ := json.MarshalIndent(struct {
				Symbol     string           `json:"symbol"`
				Components []core.Component `json:"components"`
			}{Symbol: syn3700.Symbol, Components: snap.Components}, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(snapshotCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show telemetry",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			tele := syn3700.Telemetry()
			data, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	controllersCmd := &cobra.Command{
		Use:   "controllers",
		Short: "List controllers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			list := syn3700.ControllersList()
			data, _ := json.MarshalIndent(list, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(controllersCmd)

	valueCmd := &cobra.Command{
		Use:   "value <token:price>...",
		Short: "Compute index value using token prices",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			prices := make(map[string]float64)
			for _, pair := range args {
				parts := strings.Split(pair, ":")
				if len(parts) != 2 {
					return fmt.Errorf("invalid price pair %s", pair)
				}
				p, err := strconv.ParseFloat(parts[1], 64)
				if err != nil || p < 0 {
					return fmt.Errorf("invalid price for %s", parts[0])
				}
				prices[parts[0]] = p
			}
			val := syn3700.Value(prices)
			data, _ := json.MarshalIndent(map[string]float64{"value": val}, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(valueCmd)

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance",
		Short: "Compute rebalance guidance",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			report := syn3700.Rebalance(wallet.Address)
			data, _ := json.MarshalIndent(report, "", "  ")
			cmd.Println(string(data))
			if err := persistStage73(); err != nil {
				return err
			}
			return nil
		},
	}
	rebalanceCmd.Flags().String("wallet", "", "controller wallet path")
	rebalanceCmd.Flags().String("password", "", "controller wallet password")
	rebalanceCmd.MarkFlagRequired("wallet")
	rebalanceCmd.MarkFlagRequired("password")
	cmd.AddCommand(rebalanceCmd)

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			data, _ := json.MarshalIndent(syn3700.Audit(), "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}
