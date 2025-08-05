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
		Short: "Index token management",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise index token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			syn3700 = core.NewSYN3700Token(name, symbol)
			fmt.Println("token initialised")
		},
	}
	initCmd.Flags().String("name", "", "name")
	initCmd.Flags().String("symbol", "", "symbol")
	cmd.AddCommand(initCmd)

	addCmd := &cobra.Command{
		Use:   "add <token> <weight>",
		Short: "Add component",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn3700 == nil {
				fmt.Println("token not initialised")
				return
			}
			w, _ := strconv.ParseFloat(args[1], 64)
			syn3700.AddComponent(args[0], w)
			fmt.Println("component added")
		},
	}
	cmd.AddCommand(addCmd)

	removeCmd := &cobra.Command{
		Use:   "remove <token>",
		Short: "Remove component",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn3700 == nil {
				fmt.Println("token not initialised")
				return
			}
			if err := syn3700.RemoveComponent(args[0]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("component removed")
			}
		},
	}
	cmd.AddCommand(removeCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List components",
		Run: func(cmd *cobra.Command, args []string) {
			if syn3700 == nil {
				fmt.Println("token not initialised")
				return
			}
			for _, c := range syn3700.ListComponents() {
				fmt.Printf("%s %.2f\n", c.Token, c.Weight)
			}
		},
	}
	cmd.AddCommand(listCmd)

	valueCmd := &cobra.Command{
		Use:   "value",
		Short: "Compute index value from price list",
		Run: func(cmd *cobra.Command, args []string) {
			if syn3700 == nil {
				fmt.Println("token not initialised")
				return
			}
			priceStr, _ := cmd.Flags().GetString("prices")
			m := make(map[string]float64)
			for _, kv := range strings.Split(priceStr, ",") {
				parts := strings.SplitN(kv, "=", 2)
				if len(parts) != 2 {
					continue
				}
				v, _ := strconv.ParseFloat(parts[1], 64)
				m[parts[0]] = v
			}
			fmt.Println(syn3700.Value(m))
		},
	}
	valueCmd.Flags().String("prices", "", "token=price comma-separated")
	cmd.AddCommand(valueCmd)

	rootCmd.AddCommand(cmd)
}
