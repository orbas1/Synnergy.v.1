package cli

import (
	"context"
	"os"

	synn "synnergy"
	"synnergy/core"

	"github.com/spf13/cobra"
)

var marketplace = core.NewSmartContractMarketplace(core.NewSimpleVM())

func init() {
	cmd := &cobra.Command{
		Use:   "marketplace",
		Short: "Deploy and trade smart contracts",
	}

	deployCmd := &cobra.Command{
		Use:   "deploy [wasm] [owner]",
		Args:  cobra.ExactArgs(2),
		Short: "Deploy a WASM contract to the marketplace",
		RunE: func(cmd *cobra.Command, args []string) error {
			wasm, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			gas := synn.GasCost("DeploySmartContract")
			addr, err := marketplace.DeployContract(cmd.Context(), wasm, "", gas, args[1])
			if err != nil {
				return err
			}
			printOutput(map[string]string{"address": addr})
			return nil
		},
	}

	tradeCmd := &cobra.Command{
		Use:   "trade [addr] [newOwner]",
		Args:  cobra.ExactArgs(2),
		Short: "Transfer ownership of a deployed contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := marketplace.TradeContract(context.Background(), args[0], args[1]); err != nil {
				return err
			}
			printOutput(map[string]string{"status": "traded"})
			return nil
		},
	}

	cmd.AddCommand(deployCmd, tradeCmd)
	rootCmd.AddCommand(cmd)
}
