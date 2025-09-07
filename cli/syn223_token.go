package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn223 *core.SYN223Token

func parseMap(input string) map[string]uint64 {
	m := make(map[string]uint64)
	if input == "" {
		return m
	}
	parts := strings.Split(input, ",")
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		var amt uint64
		fmt.Sscanf(kv[1], "%d", &amt)
		m[kv[0]] = amt
	}
	return m
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn223",
		Short: "SYN223 token utilities",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the SYN223 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			supply, _ := cmd.Flags().GetUint64("supply")
			if name == "" || symbol == "" || owner == "" || supply == 0 {
				return fmt.Errorf("name, symbol, owner and supply are required")
			}
			syn223 = core.NewSYN223Token(name, symbol, owner, supply)
			fmt.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("name", "", "token name")
	initCmd.Flags().String("symbol", "", "token symbol")
	initCmd.Flags().String("owner", "", "owner address")
	initCmd.Flags().Uint64("supply", 0, "initial supply")
	cmd.AddCommand(initCmd)

	wlCmd := &cobra.Command{
		Use:   "whitelist <addr>",
		Short: "Add address to whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn223 == nil {
				return fmt.Errorf("token not initialised")
			}
			syn223.AddToWhitelist(args[0])
			fmt.Println("whitelisted")
			return nil
		},
	}
	cmd.AddCommand(wlCmd)

	uwlCmd := &cobra.Command{
		Use:   "unwhitelist <addr>",
		Short: "Remove address from whitelist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn223 == nil {
				return fmt.Errorf("token not initialised")
			}
			syn223.RemoveFromWhitelist(args[0])
			fmt.Println("removed from whitelist")
			return nil
		},
	}
	cmd.AddCommand(uwlCmd)

	blCmd := &cobra.Command{
		Use:   "blacklist <addr>",
		Short: "Add address to blacklist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn223 == nil {
				return fmt.Errorf("token not initialised")
			}
			syn223.AddToBlacklist(args[0])
			fmt.Println("blacklisted")
			return nil
		},
	}
	cmd.AddCommand(blCmd)

	ublCmd := &cobra.Command{
		Use:   "unblacklist <addr>",
		Short: "Remove address from blacklist",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn223 == nil {
				return fmt.Errorf("token not initialised")
			}
			syn223.RemoveFromBlacklist(args[0])
			fmt.Println("removed from blacklist")
			return nil
		},
	}
	cmd.AddCommand(ublCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amt>",
		Short: "Transfer tokens",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn223 == nil {
				return fmt.Errorf("token not initialised")
			}
			var amt uint64
			fmt.Sscanf(args[2], "%d", &amt)
			if err := syn223.Transfer(args[0], args[1], amt); err != nil {
				return err
			}
			fmt.Println("transfer complete")
			return nil
		},
	}
	cmd.AddCommand(transferCmd)

	balCmd := &cobra.Command{
		Use:   "balance <addr>",
		Short: "Show balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn223 == nil {
				return fmt.Errorf("token not initialised")
			}
			fmt.Println(syn223.BalanceOf(args[0]))
			return nil
		},
	}
	cmd.AddCommand(balCmd)

	rootCmd.AddCommand(cmd)
}
