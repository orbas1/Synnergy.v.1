package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var daoTokenLedger = core.NewDAOTokenLedger()

func init() {
	tokenCmd := &cobra.Command{
		Use:   "dao-token",
		Short: "DAO token ledger operations",
	}

	var mintJSON bool
	var mintPub, mintMsg, mintSig string
	mintCmd := &cobra.Command{
		Use:   "mint <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Mint tokens to an address",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Mint")
			ok, err := VerifySignature(mintPub, mintMsg, mintSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid amount")
				return
			}
			daoTokenLedger.Mint(args[0], amt)
			if mintJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "minted"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "minted")
		},
	}
	mintCmd.Flags().BoolVar(&mintJSON, "json", false, "output as JSON")
	mintCmd.Flags().StringVar(&mintPub, "pub", "", "hex encoded public key")
	mintCmd.Flags().StringVar(&mintMsg, "msg", "", "hex encoded message")
	mintCmd.Flags().StringVar(&mintSig, "sig", "", "hex encoded signature")

	var transferJSON bool
	var transferPub, transferMsg, transferSig string
	transferCmd := &cobra.Command{
		Use:   "transfer <from> <to> <amount>",
		Args:  cobra.ExactArgs(3),
		Short: "Transfer tokens",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Transfer")
			ok, err := VerifySignature(transferPub, transferMsg, transferSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			amt, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid amount")
				return
			}
			if err := daoTokenLedger.Transfer(args[0], args[1], amt); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if transferJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "transferred"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "transferred")
		},
	}
	transferCmd.Flags().BoolVar(&transferJSON, "json", false, "output as JSON")
	transferCmd.Flags().StringVar(&transferPub, "pub", "", "hex encoded public key")
	transferCmd.Flags().StringVar(&transferMsg, "msg", "", "hex encoded message")
	transferCmd.Flags().StringVar(&transferSig, "sig", "", "hex encoded signature")

	var balanceJSON bool
	balanceCmd := &cobra.Command{
		Use:   "balance <addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Get token balance",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Balance")
			bal := daoTokenLedger.Balance(args[0])
			if balanceJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]uint64{"balance": bal})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), bal)
		},
	}
	balanceCmd.Flags().BoolVar(&balanceJSON, "json", false, "output as JSON")

	var burnJSON bool
	var burnPub, burnMsg, burnSig string
	burnCmd := &cobra.Command{
		Use:   "burn <addr> <amount>",
		Args:  cobra.ExactArgs(2),
		Short: "Burn tokens from an address",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Burn")
			ok, err := VerifySignature(burnPub, burnMsg, burnSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid amount")
				return
			}
			if err := daoTokenLedger.Burn(args[0], amt); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if burnJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "burned"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "burned")
		},
	}
	burnCmd.Flags().BoolVar(&burnJSON, "json", false, "output as JSON")
	burnCmd.Flags().StringVar(&burnPub, "pub", "", "hex encoded public key")
	burnCmd.Flags().StringVar(&burnMsg, "msg", "", "hex encoded message")
	burnCmd.Flags().StringVar(&burnSig, "sig", "", "hex encoded signature")

	tokenCmd.AddCommand(mintCmd, transferCmd, balanceCmd, burnCmd)
	rootCmd.AddCommand(tokenCmd)
}
