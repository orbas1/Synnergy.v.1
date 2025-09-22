package cli

import (
	"fmt"
	"strconv"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			base, _ := cmd.Flags().GetString("base")
			quote, _ := cmd.Flags().GetString("quote")
			rate, _ := cmd.Flags().GetFloat64("rate")
			if base == "" || quote == "" {
				return fmt.Errorf("base and quote are required")
			}
			if rate <= 0 {
				return fmt.Errorf("rate must be positive")
			}
			p, err := forexRegistry.Register(base, quote, rate)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), p.PairID)
			return nil
		},
	}
	registerCmd.Flags().String("base", "", "base currency")
	registerCmd.Flags().String("quote", "", "quote currency")
	registerCmd.Flags().Float64("rate", 0, "exchange rate")
	registerCmd.MarkFlagRequired("base")
	registerCmd.MarkFlagRequired("quote")
	registerCmd.MarkFlagRequired("rate")
	cmd.AddCommand(registerCmd)

	updateCmd := &cobra.Command{
		Use:   "update <id> <rate>",
		Short: "Update exchange rate",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := strconv.ParseFloat(args[1], 64)
			if err != nil || r <= 0 {
				return fmt.Errorf("invalid rate")
			}
			if err := forexRegistry.UpdateRate(args[0], r); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(updateCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get pair info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := forexRegistry.Get(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ID:%s %s/%s Rate:%f\n", p.PairID, p.Base, p.Quote, p.Rate)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List forex pairs",
		Run: func(cmd *cobra.Command, args []string) {
			pairs := forexRegistry.List()
			for _, p := range pairs {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s/%s %f\n", p.PairID, p.Base, p.Quote, p.Rate)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
