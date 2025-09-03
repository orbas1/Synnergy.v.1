package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var genesisWallets = core.DefaultGenesisWallets()

func init() {
	genesisCmd := &cobra.Command{
		Use:   "genesis",
		Short: "Genesis utilities",
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show default genesis wallet addresses",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Genesis: %s\n", genesisWallets.Genesis)
			fmt.Printf("InternalDevelopment: %s\n", genesisWallets.InternalDevelopment)
			fmt.Printf("InternalCharity: %s\n", genesisWallets.InternalCharity)
			fmt.Printf("ExternalCharity: %s\n", genesisWallets.ExternalCharity)
			fmt.Printf("LoanPool: %s\n", genesisWallets.LoanPool)
			fmt.Printf("PassiveIncome: %s\n", genesisWallets.PassiveIncome)
			fmt.Printf("ValidatorsMiners: %s\n", genesisWallets.ValidatorsMiners)
			fmt.Printf("AuthorityNodes: %s\n", genesisWallets.AuthorityNodes)
			fmt.Printf("NodeHosts: %s\n", genesisWallets.NodeHosts)
			fmt.Printf("CreatorWallet: %s\n", genesisWallets.CreatorWallet)
		},
	}

	allocateCmd := &cobra.Command{
		Use:   "allocate [total]",
		Args:  cobra.ExactArgs(1),
		Short: "Allocate fees to genesis wallets",
		Run: func(cmd *cobra.Command, args []string) {
			total, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid total:", err)
				return
			}
			dist := core.AllocateToGenesisWallets(total, genesisWallets)
			for addr, amt := range dist {
				fmt.Printf("%s: %d\n", addr, amt)
			}
		},
	}


	initBlockCmd := &cobra.Command{
		Use:   "init-block",
		Short: "Initialise the chain's genesis block",

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise the genesis block",

		RunE: func(cmd *cobra.Command, args []string) error {
			stats, _, err := currentNode.InitGenesis(genesisWallets)
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "hash: %s\ncirculating: %d\nremaining: %d\n", stats.Hash, stats.Circulating, stats.Remaining)

			fmt.Fprintf(cmd.OutOrStdout(), "genesis block %s height %d\n", stats.Hash, stats.Height)

			return nil
		},
	}


	genesisCmd.AddCommand(showCmd, allocateCmd, initBlockCmd)

	genesisCmd.AddCommand(showCmd, allocateCmd, initCmd)

	rootCmd.AddCommand(genesisCmd)
}
