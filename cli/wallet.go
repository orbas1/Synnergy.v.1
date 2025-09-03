package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	walletCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Wallet operations",
	}
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Generate a new wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := core.NewWallet()
			if err != nil {
				return err
			}
			fmt.Printf("address: %s\n", w.Address)
			return nil
		},
	}
	walletCmd.AddCommand(newCmd)
	rootCmd.AddCommand(walletCmd)
}
