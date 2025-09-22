package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func parseCredentialPair(value string) (string, string, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("expected path:password format")
	}
	return parts[0], parts[1], nil
}

func requireControllerWallet(cmd *cobra.Command) (*core.Wallet, error) {
	path, _ := cmd.Flags().GetString("wallet")
	password, _ := cmd.Flags().GetString("password")
	if path == "" || password == "" {
		return nil, fmt.Errorf("wallet and password required")
	}
	return loadWallet(path, password)
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn3700",
		Short: "SYN3700 index token",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the index token",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Init")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			controllerArg, _ := cmd.Flags().GetString("controller")
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			if controllerArg == "" {
				return fmt.Errorf("controller required")
			}
			path, password, err := parseCredentialPair(controllerArg)
			if err != nil {
				return err
			}
			wallet, err := loadWallet(path, password)
			if err != nil {
				return err
			}
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := core.NewSYN3700Token(name, symbol)
			token.AddController(wallet.Address)
			store.SetIndex(token)
			markStage73Dirty()
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("controller", "", "controller credential path:password")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("symbol")
	_ = initCmd.MarkFlagRequired("controller")
	cmd.AddCommand(initCmd)

	addCmd := &cobra.Command{
		Use:   "add <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Add component to index",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Add")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			if !token.HasController(wallet.Address) {
				return fmt.Errorf("wallet not authorised")
			}
			weight, _ := cmd.Flags().GetFloat64("weight")
			drift, _ := cmd.Flags().GetFloat64("drift")
			if err := token.AddComponent(wallet.Address, strings.ToUpper(args[0]), weight, drift); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("component added")
			return nil
		},
	}
	addCmd.Flags().Float64("weight", 0, "component weight")
	addCmd.Flags().Float64("drift", 0, "allowed drift ratio")
	addCmd.Flags().String("wallet", "", "controller wallet path")
	addCmd.Flags().String("password", "", "wallet password")
	_ = addCmd.MarkFlagRequired("weight")
	_ = addCmd.MarkFlagRequired("wallet")
	_ = addCmd.MarkFlagRequired("password")
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <token>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove component from index",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Remove")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			if !token.HasController(wallet.Address) {
				return fmt.Errorf("wallet not authorised")
			}
			if err := token.RemoveComponent(wallet.Address, strings.ToUpper(args[0])); err != nil {
				return err
			}
			markStage73Dirty()
			cmd.Println("component removed")
			return nil
		},
	}
	removeCmd.Flags().String("wallet", "", "controller wallet path")
	removeCmd.Flags().String("password", "", "wallet password")
	_ = removeCmd.MarkFlagRequired("wallet")
	_ = removeCmd.MarkFlagRequired("password")
	cmd.AddCommand(removeCmd)

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Show current index snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Snapshot")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			snap := token.Snapshot()
			data, _ := json.MarshalIndent(snap, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(snapshotCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show telemetry for the index",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Status")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			tele := token.Telemetry()
			data, _ := json.MarshalIndent(tele, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	controllersCmd := &cobra.Command{
		Use:   "controllers",
		Short: "List controller addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Controllers")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			list := token.Controllers()
			data, _ := json.MarshalIndent(list, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(controllersCmd)

	valueCmd := &cobra.Command{
		Use:   "value <token:price>...",
		Args:  cobra.MinimumNArgs(1),
		Short: "Compute index value using token prices",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Value")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
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
			val := token.Value(prices)
			payload := map[string]float64{"value": val}
			data, _ := json.MarshalIndent(payload, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(valueCmd)

	rebalanceCmd := &cobra.Command{
		Use:   "rebalance",
		Short: "Rebalance index weights to targets",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Rebalance")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			wallet, err := requireControllerWallet(cmd)
			if err != nil {
				return err
			}
			if !token.HasController(wallet.Address) {
				return fmt.Errorf("wallet not authorised")
			}
			updates := token.Rebalance(wallet.Address)
			markStage73Dirty()
			data, _ := json.MarshalIndent(updates, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	rebalanceCmd.Flags().String("wallet", "", "controller wallet path")
	rebalanceCmd.Flags().String("password", "", "wallet password")
	_ = rebalanceCmd.MarkFlagRequired("wallet")
	_ = rebalanceCmd.MarkFlagRequired("password")
	cmd.AddCommand(rebalanceCmd)

	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "Show audit log",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("Syn3700Audit")
			store, err := stage73State()
			if err != nil {
				return err
			}
			token := store.Index()
			if token == nil {
				return fmt.Errorf("token not initialised")
			}
			events := token.AuditTrail()
			data, _ := json.MarshalIndent(events, "", "  ")
			cmd.Println(string(data))
			return nil
		},
	}
	cmd.AddCommand(auditCmd)

	rootCmd.AddCommand(cmd)
}
