package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	walletCmd := &cobra.Command{
		Use:   "wallet",
		Short: "Wallet operations",
	}
	var outFile, password string
	newCmd := &cobra.Command{
		Use:   "new",
		Short: "Generate a new wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("WalletNew")
			w, err := core.NewWallet()
			if err != nil {
				return err
			}
			if outFile != "" {
				if err := w.Save(outFile, password); err != nil {
					return err
				}
			}
			printOutput(map[string]string{"address": w.Address, "path": outFile})
			return nil
		},
	}
	newCmd.Flags().StringVar(&outFile, "out", "", "write encrypted wallet to file")
	newCmd.Flags().StringVar(&password, "password", "", "encryption password for wallet file")
	walletCmd.AddCommand(newCmd)
	rootCmd.AddCommand(walletCmd)
}
