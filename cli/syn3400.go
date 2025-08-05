package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var forexRegistry = tokens.NewForexRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3400",
		Short: "Forex pair registry",
	}

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a forex pair",
		Run: func(cmd *cobra.Command, args []string) {
			base, _ := cmd.Flags().GetString("base")
			quote, _ := cmd.Flags().GetString("quote")
			rate, _ := cmd.Flags().GetFloat64("rate")
			p := forexRegistry.Register(base, quote, rate)
			fmt.Println(p.PairID)
		},
	}
	registerCmd.Flags().String("base", "", "base currency")
	registerCmd.Flags().String("quote", "", "quote currency")
	registerCmd.Flags().Float64("rate", 1.0, "exchange rate")
	cmd.AddCommand(registerCmd)

	updateCmd := &cobra.Command{
		Use:   "update <id> <rate>",
		Short: "Update exchange rate",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var r float64
			fmt.Sscanf(args[1], "%f", &r)
			if err := forexRegistry.UpdateRate(args[0], r); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(updateCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get pair info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := forexRegistry.Get(args[0])
			if !ok {
				fmt.Println("pair not found")
				return
			}
			fmt.Printf("ID:%s %s/%s Rate:%f\n", p.PairID, p.Base, p.Quote, p.Rate)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List forex pairs",
		Run: func(cmd *cobra.Command, args []string) {
			pairs := forexRegistry.List()
			for _, p := range pairs {
				fmt.Printf("%s %s/%s %f\n", p.PairID, p.Base, p.Quote, p.Rate)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
