package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	feeCmd := &cobra.Command{Use: "fees", Short: "Fee distribution utilities"}

	distributeCmd := &cobra.Command{
		Use:   "distribute [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Distribute fees across genesis wallets",
		Run: func(cmd *cobra.Command, args []string) {
			total, _ := strconv.ParseUint(args[0], 10, 64)
			wallets := core.DefaultGenesisWallets()
			alloc := core.AllocateToGenesisWallets(total, wallets)
			fmt.Printf("internal development: %d\n", alloc[wallets.InternalDevelopment])
			fmt.Printf("internal charity: %d\n", alloc[wallets.InternalCharity])
			fmt.Printf("external charity: %d\n", alloc[wallets.ExternalCharity])
			fmt.Printf("loan pool: %d\n", alloc[wallets.LoanPool])
			fmt.Printf("passive income: %d\n", alloc[wallets.PassiveIncome])
			fmt.Printf("validators/miners: %d\n", alloc[wallets.ValidatorsMiners])
			fmt.Printf("node hosts: %d\n", alloc[wallets.NodeHosts])
			fmt.Printf("creator wallet: %d\n", alloc[wallets.CreatorWallet])
		},
	}

	walletsCmd := &cobra.Command{
		Use:   "wallets",
		Short: "Display default genesis wallet addresses",
		Run: func(cmd *cobra.Command, args []string) {
			w := core.DefaultGenesisWallets()
			fmt.Printf("genesis: %s\n", w.Genesis)
			fmt.Printf("internal development: %s\n", w.InternalDevelopment)
			fmt.Printf("internal charity: %s\n", w.InternalCharity)
			fmt.Printf("external charity: %s\n", w.ExternalCharity)
			fmt.Printf("loan pool: %s\n", w.LoanPool)
			fmt.Printf("passive income: %s\n", w.PassiveIncome)
			fmt.Printf("validators/miners: %s\n", w.ValidatorsMiners)
			fmt.Printf("node hosts: %s\n", w.NodeHosts)
			fmt.Printf("creator wallet: %s\n", w.CreatorWallet)
		},
	}

	feeCmd.AddCommand(distributeCmd, walletsCmd)
	rootCmd.AddCommand(feeCmd)
}
