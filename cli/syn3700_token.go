package cli

import (
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
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			syn3700 = core.NewSYN3700Token(name, symbol)
			cmd.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	_ = initCmd.MarkFlagRequired("name")
	_ = initCmd.MarkFlagRequired("symbol")
	cmd.AddCommand(initCmd)

	addCmd := &cobra.Command{
		Use:   "add <token> <weight>",
		Short: "Add component to index",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			weight, err := strconv.ParseFloat(args[1], 64)
			if err != nil || weight <= 0 {
				return fmt.Errorf("invalid weight")
			}
			if args[0] == "" {
				return fmt.Errorf("token symbol required")
			}
			syn3700.AddComponent(args[0], weight)
			cmd.Println("component added")
			return nil
		},
	}
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <token>",
		Short: "Remove component from index",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			if err := syn3700.RemoveComponent(args[0]); err != nil {
				return err
			}
			cmd.Println("component removed")
			return nil
		},
	}
	cmd.AddCommand(removeCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List index components",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn3700 == nil {
				return fmt.Errorf("token not initialised")
			}
			comps := syn3700.ListComponents()
			for _, c := range comps {
				cmd.Printf("%s %.2f\n", c.Token, c.Weight)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

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
			cmd.Printf("%.2f\n", val)
			return nil
		},
	}
	cmd.AddCommand(valueCmd)

	rootCmd.AddCommand(cmd)
}
