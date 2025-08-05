package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var faucet *core.Faucet

func init() {
	cmd := &cobra.Command{
		Use:   "faucet",
		Short: "Interact with the test token faucet",
	}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise faucet balance and parameters",
		Run: func(cmd *cobra.Command, args []string) {
			bal, _ := cmd.Flags().GetUint64("balance")
			amt, _ := cmd.Flags().GetUint64("amount")
			cd, _ := cmd.Flags().GetDuration("cooldown")
			faucet = core.NewFaucet(bal, amt, cd)
			fmt.Println("faucet initialised")
		},
	}
	initCmd.Flags().Uint64("balance", 1000, "initial balance")
	initCmd.Flags().Uint64("amount", 1, "dispense amount")
	initCmd.Flags().Duration("cooldown", time.Minute, "cooldown between requests")
	cmd.AddCommand(initCmd)

	requestCmd := &cobra.Command{
		Use:   "request <addr>",
		Short: "Request funds from the faucet",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if faucet == nil {
				fmt.Println("faucet not initialised")
				return
			}
			amt, err := faucet.Request(args[0])
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("dispensed %d tokens\n", amt)
		},
	}
	cmd.AddCommand(requestCmd)

	balCmd := &cobra.Command{
		Use:   "balance",
		Short: "Show remaining faucet balance",
		Run: func(cmd *cobra.Command, args []string) {
			if faucet == nil {
				fmt.Println("faucet not initialised")
				return
			}
			fmt.Printf("balance: %d\n", faucet.Balance())
		},
	}
	cmd.AddCommand(balCmd)

	cfgCmd := &cobra.Command{
		Use:   "config",
		Short: "Update faucet configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if faucet == nil {
				fmt.Println("faucet not initialised")
				return
			}
			amt, _ := cmd.Flags().GetUint64("amount")
			cd, _ := cmd.Flags().GetDuration("cooldown")
			faucet.UpdateConfig(amt, cd)
			fmt.Println("configuration updated")
		},
	}
	cfgCmd.Flags().Uint64("amount", 1, "dispense amount")
	cfgCmd.Flags().Duration("cooldown", time.Minute, "cooldown between requests")
	cmd.AddCommand(cfgCmd)

	rootCmd.AddCommand(cmd)
}
