package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var tokenRegistry = tokens.NewRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Token registry utilities",
	}

	nextCmd := &cobra.Command{
		Use:   "nextid",
		Short: "Generate next token ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := tokenRegistry.NextID()
			fmt.Fprintln(cmd.OutOrStdout(), id)
			return nil
		},
	}
	cmd.AddCommand(nextCmd)

	registerBaseCmd := &cobra.Command{
		Use:   "register-base",
		Short: "Register the base token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if baseToken == nil {
				return fmt.Errorf("base token not initialised")
			}
			tokenRegistry.Register(baseToken)
			fmt.Fprintln(cmd.OutOrStdout(), "base token registered")
			return nil
		},
	}
	cmd.AddCommand(registerBaseCmd)

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Short: "Show token info by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			info, ok := tokenRegistry.Info(tokens.TokenID(id))
			if !ok {
				return fmt.Errorf("token not found")
			}
			fmt.Fprintf(cmd.OutOrStdout(), "ID:%d Name:%s Symbol:%s Decimals:%d Supply:%d\n", info.ID, info.Name, info.Symbol, info.Decimals, info.TotalSupply)
			return nil
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			infos := tokenRegistry.List()
			for _, i := range infos {
				fmt.Fprintf(cmd.OutOrStdout(), "ID:%d Name:%s Symbol:%s Supply:%d\n", i.ID, i.Name, i.Symbol, i.TotalSupply)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
