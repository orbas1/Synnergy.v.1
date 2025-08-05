package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var forexReg = tokens.NewForexRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn3400",
		Short: "Forex pair registry",
	}

	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register forex pair",
		Run: func(cmd *cobra.Command, args []string) {
			base, _ := cmd.Flags().GetString("base")
			quote, _ := cmd.Flags().GetString("quote")
			rate, _ := cmd.Flags().GetFloat64("rate")
			p := forexReg.Register(base, quote, rate)
			fmt.Println(p.PairID)
		},
	}
	registerCmd.Flags().String("base", "", "base currency")
	registerCmd.Flags().String("quote", "", "quote currency")
	registerCmd.Flags().Float64("rate", 1.0, "exchange rate")
	cmd.AddCommand(registerCmd)

	updateCmd := &cobra.Command{
		Use:   "update-rate <pairID> <rate>",
		Short: "Update exchange rate",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var rate float64
			fmt.Sscanf(args[1], "%f", &rate)
			if err := forexReg.UpdateRate(args[0], rate); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("updated")
			}
		},
	}
	cmd.AddCommand(updateCmd)

	getCmd := &cobra.Command{
		Use:   "get <pairID>",
		Short: "Get pair info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, ok := forexReg.Get(args[0])
			if !ok {
				fmt.Println("pair not found")
				return
			}
			fmt.Printf("%+v\n", *p)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List pairs",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range forexReg.List() {
				fmt.Printf("%+v\n", *p)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
