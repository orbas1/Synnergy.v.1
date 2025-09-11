package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	creatorCmd := &cobra.Command{
		Use:   "creator",
		Short: "Creator wallet controls",
	}

	var walletPath, password string
	disableCmd := &cobra.Command{
		Use:   "disable-distribution",
		Short: "Redirect creator fee share to node hosts",
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := core.LoadWallet(walletPath, password)
			if err != nil {
				return err
			}
			if w.Address != genesisWallets.CreatorWallet {
				return fmt.Errorf("wallet does not match creator wallet")
			}
			core.SetCreatorDistribution(false)
			fmt.Fprintln(cmd.OutOrStdout(), "creator distribution disabled")
			return nil
		},
	}
	disableCmd.Flags().StringVar(&walletPath, "wallet", "", "path to creator wallet file")
	disableCmd.Flags().StringVar(&password, "password", "", "wallet password")
	disableCmd.MarkFlagRequired("wallet")
	disableCmd.MarkFlagRequired("password")

	creatorCmd.AddCommand(disableCmd)
	rootCmd.AddCommand(creatorCmd)
}
