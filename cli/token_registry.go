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
		Run: func(cmd *cobra.Command, args []string) {
			id := tokenRegistry.NextID()
			fmt.Println(id)
		},
	}
	cmd.AddCommand(nextCmd)

	registerBaseCmd := &cobra.Command{
		Use:   "register-base",
		Short: "Register the base token",
		Run: func(cmd *cobra.Command, args []string) {
			if baseToken == nil {
				fmt.Println("base token not initialised")
				return
			}
			tokenRegistry.Register(baseToken)
			fmt.Println("base token registered")
		},
	}
	cmd.AddCommand(registerBaseCmd)

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Short: "Show token info by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			info, ok := tokenRegistry.Info(tokens.TokenID(id))
			if !ok {
				fmt.Println("token not found")
				return
			}
			fmt.Printf("ID:%d Name:%s Symbol:%s Decimals:%d Supply:%d\n", info.ID, info.Name, info.Symbol, info.Decimals, info.TotalSupply)
		},
	}
	cmd.AddCommand(infoCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all registered tokens",
		Run: func(cmd *cobra.Command, args []string) {
			infos := tokenRegistry.List()
			for _, i := range infos {
				fmt.Printf("ID:%d Name:%s Symbol:%s Supply:%d\n", i.ID, i.Name, i.Symbol, i.TotalSupply)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
