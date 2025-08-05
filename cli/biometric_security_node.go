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
		Use:   "enroll [addr] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Enroll biometric data for an address",
		Run: func(cmd *cobra.Command, args []string) {
			secureNode.Enroll(args[0], []byte(args[1]))
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
		Use:   "auth [addr] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Authenticate biometric data",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(secureNode.Authenticate(args[0], []byte(args[1])))
		},
	}

	addTxCmd := &cobra.Command{
		Use:   "addtx [addr] [data] [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(7),
		Short: "Securely add a transaction to the mempool",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[4], 10, 64)
			fee, _ := strconv.ParseUint(args[5], 10, 64)
			nonce, _ := strconv.ParseUint(args[6], 10, 64)
			tx := core.NewTransaction(args[2], args[3], amt, fee, nonce)
			if err := secureNode.SecureAddTransaction(args[0], []byte(args[1]), tx); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	bsnCmd.AddCommand(enrollCmd, removeCmd, authCmd, addTxCmd)
	rootCmd.AddCommand(bsnCmd)
}
