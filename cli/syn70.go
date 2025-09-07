package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn70 *tokens.SYN70Token

func init() {
	cmd := &cobra.Command{
		Use:   "syn70",
		Short: "SYN70 game asset token",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise SYN70 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint32("decimals")
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			id := tokenRegistry.NextID()
			syn70 = tokens.NewSYN70Token(id, name, symbol, uint8(dec))
			tokenRegistry.Register(syn70)
			fmt.Fprintln(cmd.OutOrStdout(), "syn70 initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().Uint32("decimals", 0, "decimal places")
	initCmd.MarkFlagRequired("name")
	initCmd.MarkFlagRequired("symbol")
	cmd.AddCommand(initCmd)

	registerCmd := &cobra.Command{
		Use:   "register <id> <owner> <name> <game>",
		Short: "Register an in-game asset",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			if err := syn70.RegisterAsset(args[0], args[1], args[2], args[3]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "asset registered")
			return nil
		},
	}
	cmd.AddCommand(registerCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <newOwner>",
		Short: "Transfer asset ownership",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			if err := syn70.TransferAsset(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "asset transferred")
			return nil
		},
	}
	cmd.AddCommand(transferCmd)

	attrCmd := &cobra.Command{
		Use:   "setattr <id> <key> <value>",
		Short: "Set asset attribute",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			if err := syn70.SetAttribute(args[0], args[1], args[2]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "attribute set")
			return nil
		},
	}
	cmd.AddCommand(attrCmd)

	achCmd := &cobra.Command{
		Use:   "achievement <id> <name>",
		Short: "Record achievement",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			if err := syn70.AddAchievement(args[0], args[1]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "achievement recorded")
			return nil
		},
	}
	cmd.AddCommand(achCmd)

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Short: "Show asset info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			a, err := syn70.AssetInfo(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ID:%s Owner:%s Name:%s Game:%s\n", a.ID, a.Owner, a.Name, a.Game)
			if len(a.Attributes) > 0 {
				for k, v := range a.Attributes {
					fmt.Fprintf(cmd.OutOrStdout(), "%s=%s ", k, v)
				}
				fmt.Fprintln(cmd.OutOrStdout())
			}
			if len(a.Achievements) > 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "Achievements:", strings.Join(a.Achievements, ","))
			}
			return nil
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			assets := syn70.ListAssets()
			for _, a := range assets {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s %s %s\n", a.ID, a.Owner, a.Name, a.Game)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	balCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show token balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn70 == nil {
				return fmt.Errorf("token not initialised")
			}
			fmt.Fprintln(cmd.OutOrStdout(), syn70.BalanceOf(args[0]))
			return nil
		},
	}
	cmd.AddCommand(balCmd)

	rootCmd.AddCommand(cmd)
}
