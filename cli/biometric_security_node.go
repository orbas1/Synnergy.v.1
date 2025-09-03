package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var secureNode = core.NewBiometricSecurityNode(currentNode, nil)

func init() {
	bsnCmd := &cobra.Command{
		Use:   "bsn",
		Short: "Biometric security node operations",
	}

	enrollCmd := &cobra.Command{
		Use:   "enroll [addr] [data] [pubKeyHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Enroll biometric data for an address",
		Run: func(cmd *cobra.Command, args []string) {
			pub, err := parsePubKey(args[2])
			if err != nil {
				fmt.Println("invalid public key:", err)
				return
			}
			secureNode.Enroll(args[0], []byte(args[1]), pub)
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove biometric data",
		Run: func(cmd *cobra.Command, args []string) {
			secureNode.Remove(args[0])
		},
	}

	authCmd := &cobra.Command{
		Use:   "auth [addr] [data] [sigHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Authenticate biometric data",
		Run: func(cmd *cobra.Command, args []string) {
			sig, err := decodeSig(args[2])
			if err != nil {
				fmt.Println("invalid signature:", err)
				return
			}
			fmt.Println(secureNode.Authenticate(args[0], []byte(args[1]), sig))
		},
	}

	addTxCmd := &cobra.Command{
		Use:   "addtx [addr] [data] [sigHex] [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(8),
		Short: "Securely add a transaction to the mempool",
		Run: func(cmd *cobra.Command, args []string) {
			sig, err := decodeSig(args[2])
			if err != nil {
				fmt.Println("invalid signature:", err)
				return
			}
			amt, _ := strconv.ParseUint(args[5], 10, 64)
			fee, _ := strconv.ParseUint(args[6], 10, 64)
			nonce, _ := strconv.ParseUint(args[7], 10, 64)
			tx := core.NewTransaction(args[3], args[4], amt, fee, nonce)
			if err := secureNode.SecureAddTransaction(args[0], []byte(args[1]), sig, tx); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	bsnCmd.AddCommand(enrollCmd, removeCmd, authCmd, addTxCmd)
	rootCmd.AddCommand(bsnCmd)
}
