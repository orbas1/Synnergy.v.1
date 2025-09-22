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
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the index token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			controller, _ := cmd.Flags().GetString("controller")
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			if controller == "" {
				return fmt.Errorf("controller credentials required")
			}
			path, password, err := parseWalletCredential(controller)
			if err != nil {
				return err
			}
			wallet, err := loadWallet(path, password)
			if err != nil {
				return err
			}
			syn3700 = core.NewSYN3700Token(name, symbol)
			syn3700.AddController(wallet.Address)
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("controller", "", "controller wallet credentials path:password")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("symbol")
	_ = initCmd.MarkFlagRequired("controller")
	cmd.AddCommand(initCmd)

	addCmd := &cobra.Command{
		Use:   "add <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Add component to index",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			weight, _ := cmd.Flags().GetFloat64("weight")
			drift, _ := cmd.Flags().GetFloat64("drift")
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath == "" || password == "" {
				return fmt.Errorf("wallet and password required")
			}
			wallet, err := loadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if err := syn3700.AddComponent(strings.ToUpper(args[0]), weight, drift, wallet.Address); err != nil {
				return err
			}
			cmd.Println("component added")
			return nil
		},
	}
	addCmd.Flags().Float64("weight", 0, "component weight")
	addCmd.Flags().Float64("drift", 0, "acceptable drift")
	addCmd.Flags().String("wallet", "", "path to controller wallet")
	addCmd.Flags().String("password", "", "controller wallet password")
	_ = addCmd.MarkFlagRequired("weight")
	_ = addCmd.MarkFlagRequired("wallet")
	_ = addCmd.MarkFlagRequired("password")
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
			if err := syn3700.RemoveComponent(strings.ToUpper(args[0]), wallet.Address); err != nil {
				return err
			}
			cmd.Println("component removed")
			return nil
		},
	}
	removeCmd.Flags().String("wallet", "", "path to controller wallet")
	removeCmd.Flags().String("password", "", "controller wallet password")
	_ = removeCmd.MarkFlagRequired("wallet")
	_ = removeCmd.MarkFlagRequired("password")
	cmd.AddCommand(removeCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List index components",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			comps := syn3700.ListComponents()
			payload, _ := json.MarshalIndent(comps, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Snapshot index state",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			out := struct {
				Symbol     string           `json:"symbol"`
				Components []core.Component `json:"components"`
			}{Symbol: syn3700.Symbol, Components: syn3700.ListComponents()}
			payload, _ := json.MarshalIndent(out, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	cmd.AddCommand(snapshotCmd)

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
				prices[strings.ToUpper(parts[0])] = p
			}
			val := syn3700.Value(prices)
			payload, _ := json.MarshalIndent(map[string]float64{"value": val}, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	cmd.AddCommand(valueCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show component and controller counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			payload, _ := json.MarshalIndent(struct {
				ComponentCount  int `json:"component_count"`
				ControllerCount int `json:"controller_count"`
			}{ComponentCount: syn3700.ComponentCount(), ControllerCount: syn3700.ControllerCount()}, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	controllersCmd := &cobra.Command{
		Use:   "controllers",
		Short: "List registered controllers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			payload, _ := json.MarshalIndent(syn3700.Controllers(), "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	cmd.AddCommand(controllersCmd)

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance",
		Short: "Rebalance index weights",
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
			changes, err := syn3700.Rebalance(wallet.Address)
			if err != nil {
				return err
			}
			payload, _ := json.MarshalIndent(changes, "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	rebalanceCmd.Flags().String("wallet", "", "path to controller wallet")
	rebalanceCmd.Flags().String("password", "", "controller wallet password")
	_ = rebalanceCmd.MarkFlagRequired("wallet")
	_ = rebalanceCmd.MarkFlagRequired("password")
	cmd.AddCommand(rebalanceCmd)

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show audit trail",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			payload, _ := json.MarshalIndent(syn3700.Audit(), "", "  ")
			cmd.Println(string(payload))
			return nil
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}
