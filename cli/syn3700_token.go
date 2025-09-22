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
		Short: "SYN3700 index token operations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return ensureSyn3700Loaded()
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
			if controllerSpec == "" {
				return fmt.Errorf("controller required")
			}
			addr, err := walletAddressFromSpec(controllerSpec)
			if err != nil {
				return err
			}
			syn3700 = core.NewSYN3700Token(name, symbol)
			syn3700.AddController(core.Address(addr), core.Address(addr))
			if err := persistSyn3700(); err != nil {
				return err
			}
			printOutput("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("controller", "", "controller wallet path:password")
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
			weight, err := cmd.Flags().GetFloat64("weight")
			if err != nil {
				return fmt.Errorf("invalid weight")
			}
			drift, err := cmd.Flags().GetFloat64("drift")
			if err != nil {
				return fmt.Errorf("invalid drift")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if !syn3700.IsController(core.Address(wallet.Address)) {
				return fmt.Errorf("wallet not authorised")
			}
			if err := syn3700.AddComponent(args[0], weight, drift, core.Address(wallet.Address)); err != nil {
				return err
			}
			if err := persistSyn3700(); err != nil {
				return err
			}
			printOutput("component added")
			return nil
		},
	}
	addCmd.Flags().Float64("weight", 0, "component weight")
	addCmd.Flags().Float64("drift", 0, "allowed drift")
	addWalletFlags(addCmd)
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove component from index",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if !syn3700.IsController(core.Address(wallet.Address)) {
				return fmt.Errorf("wallet not authorised")
			}
			if err := syn3700.RemoveComponent(args[0], core.Address(wallet.Address)); err != nil {
				return err
			}
			if err := persistSyn3700(); err != nil {
				return err
			}
			printOutput("component removed")
			return nil
		},
	}
	addWalletFlags(removeCmd)
	cmd.AddCommand(removeCmd)

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Emit token snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			snap := syn3700.Snapshot()
			resp := map[string]interface{}{
				"name":       snap.Name,
				"symbol":     snap.Symbol,
				"components": snap.Components,
			}
			b, _ := json.MarshalIndent(resp, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(snapshotCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show component/controller counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			tele := syn3700.Telemetry()
			b, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	controllersCmd := &cobra.Command{
		Use:   "controllers",
		Short: "List controller addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			ctrls := syn3700.Controllers()
			b, _ := json.MarshalIndent(ctrls, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(controllersCmd)

	valueCmd := &cobra.Command{
		Use:   "value <token:price>...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Compute index value using token prices",
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
				if err != nil {
					return fmt.Errorf("invalid price for %s", parts[0])
				}
				prices[parts[0]] = p
			}
			val := syn3700.Value(prices)
			b, _ := json.MarshalIndent(map[string]float64{"value": val}, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(valueCmd)

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance",
		Short: "Show rebalance plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := walletFromFlags(cmd)
			if err != nil {
				return err
			}
			if !syn3700.IsController(core.Address(wallet.Address)) {
				return fmt.Errorf("wallet not authorised")
			}
			plan := syn3700.RebalancePlan()
			b, _ := json.MarshalIndent(plan, "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	addWalletFlags(rebalanceCmd)
	requireWalletFlags(rebalanceCmd)
	cmd.AddCommand(rebalanceCmd)

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show audit history",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			b, _ := json.MarshalIndent(syn3700.AuditLog(), "", "  ")
			cmd.Println(string(b))
			return nil
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}
