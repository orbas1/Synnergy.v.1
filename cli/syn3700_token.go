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
		Short: "Manage the SYN3700 institutional index",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the SYN3700 index",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			controllerSpec, _ := cmd.Flags().GetString("controller")
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			if controllerSpec == "" {
				return fmt.Errorf("controller flag requires path:password")
			}
			path, password, err := parseWalletDescriptor(controllerSpec)
			if err != nil {
				return err
			}
			wallet, err := loadWallet(path, password)
			if err != nil {
				return err
			}
			syn3700 = core.NewSYN3700Token(name, symbol)
			syn3700.RegisterController(wallet.Address)
			gasPrint("CreateIndexToken")
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("controller", "", "wallet descriptor path:password")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("symbol")
	_ = initCmd.MarkFlagRequired("controller")
	cmd.AddCommand(initCmd)

	addCmd := &cobra.Command{
		Use:   "add <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Add a component token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			weight, _ := cmd.Flags().GetFloat64("weight")
			drift, _ := cmd.Flags().GetFloat64("drift")
			if err := syn3700.AddComponentControlled(wallet.Address, args[0], weight, drift); err != nil {
				if strings.Contains(err.Error(), "component drift") {
					return err
				}
				return fmt.Errorf("add component: %w", err)
			}
			gasPrint("IndexAddComponent")
			cmd.Println("component added")
			return nil
		},
	}
	addCmd.Flags().Float64("weight", 0, "target portfolio weight")
	addCmd.Flags().Float64("drift", 0, "allowed drift before rebalancing")
	addCmd.Flags().String("wallet", "", "controller wallet path")
	addCmd.Flags().String("password", "", "controller wallet password")
	_ = addCmd.MarkFlagRequired("weight")
	_ = addCmd.MarkFlagRequired("wallet")
	_ = addCmd.MarkFlagRequired("password")
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a component token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if err := syn3700.RemoveComponentControlled(wallet.Address, args[0]); err != nil {
				return err
			}
			gasPrint("IndexRemoveComponent")
			cmd.Println("component removed")
			return nil
		},
	}
	removeCmd.Flags().String("wallet", "", "controller wallet path")
	removeCmd.Flags().String("password", "", "controller wallet password")
	_ = removeCmd.MarkFlagRequired("wallet")
	_ = removeCmd.MarkFlagRequired("password")
	cmd.AddCommand(removeCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List components",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			comps := syn3700.ListComponents()
			return printJSON(cmd, comps)
		},
	}
	cmd.AddCommand(listCmd)

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Show a telemetry snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			snap := syn3700.Snapshot()
			return printJSON(cmd, snap)
		},
	}
	cmd.AddCommand(snapshotCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show controller and component counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			status := syn3700.Status()
			return printJSON(cmd, status)
		},
	}
	cmd.AddCommand(statusCmd)

	controllersCmd := &cobra.Command{
		Use:   "controllers",
		Short: "List controller wallets",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			controllers := syn3700.Controllers()
			return printJSON(cmd, controllers)
		},
	}
	cmd.AddCommand(controllersCmd)

	valueCmd := &cobra.Command{
		Use:   "value <token:price>...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Calculate index value from token prices",
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
				price, err := strconv.ParseFloat(parts[1], 64)
				if err != nil {
					return fmt.Errorf("invalid price for %s", parts[0])
				}
				prices[parts[0]] = price
			}
			value := syn3700.Value(prices)
			return printJSON(cmd, map[string]float64{"value": value})
		},
	}
	cmd.AddCommand(valueCmd)

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance",
		Short: "Rebalance drifted components",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			changes, err := syn3700.RebalanceControlled(wallet.Address)
			if err != nil {
				return err
			}
			gasPrint("IndexRebalance")
			return printJSON(cmd, changes)
		},
	}
	rebalanceCmd.Flags().String("wallet", "", "controller wallet path")
	rebalanceCmd.Flags().String("password", "", "controller wallet password")
	_ = rebalanceCmd.MarkFlagRequired("wallet")
	_ = rebalanceCmd.MarkFlagRequired("password")
	cmd.AddCommand(rebalanceCmd)

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Display the administrative audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			return printJSON(cmd, syn3700.AuditTrail())
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}

func parseWalletDescriptor(spec string) (string, string, error) {
	parts := strings.SplitN(spec, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("wallet descriptor must be path:password")
	}
	return parts[0], parts[1], nil
}

func walletFromFlags(cmd *cobra.Command) (*core.Wallet, error) {
	path, _ := cmd.Flags().GetString("wallet")
	password, _ := cmd.Flags().GetString("password")
	if path == "" || password == "" {
		return nil, fmt.Errorf("wallet and password required")
	}
	return loadWallet(path, password)
}

func printJSON(cmd *cobra.Command, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	cmd.Println(string(data))
	return nil
}
