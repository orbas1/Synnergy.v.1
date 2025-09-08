package cli

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

// Stage 40 hardens the biometric security node CLI with JSON output and strict
// argument validation so automated agents and web dashboards can consume the
// results deterministically.
var (
	secureNode = core.NewBiometricSecurityNode(currentNode, nil)
	bsnJSON    bool
)

func bsnOutput(v interface{}, plain string) {
	if bsnJSON {
		b, err := json.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(plain)
	}
}

func init() {
	bsnCmd := &cobra.Command{
		Use:   "bsn",
		Short: "Biometric security node operations",
	}
	bsnCmd.PersistentFlags().BoolVar(&bsnJSON, "json", false, "output results in JSON")

	enrollCmd := &cobra.Command{
		Use:   "enroll [addr] [data] [pubKeyHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Enroll biometric data for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			pubBytes, err := hex.DecodeString(args[2])
			if err != nil || len(pubBytes) != ed25519.PublicKeySize {
				return fmt.Errorf("invalid public key: %w", err)
			}
			if err := secureNode.Enroll(args[0], []byte(args[1]), ed25519.PublicKey(pubBytes)); err != nil {
				return err
			}
			bsnOutput(map[string]string{"status": "enrolled"}, "enrolled")
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove biometric data",
		RunE: func(cmd *cobra.Command, args []string) error {
			secureNode.Remove(args[0])
			bsnOutput(map[string]string{"status": "removed"}, "removed")
			return nil
		},
	}

	authCmd := &cobra.Command{
		Use:   "auth [addr] [data] [sigHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Authenticate biometric data",
		RunE: func(cmd *cobra.Command, args []string) error {
			sig, err := hex.DecodeString(args[2])
			if err != nil {
				return fmt.Errorf("invalid signature: %w", err)
			}
			ok := secureNode.Authenticate(args[0], []byte(args[1]), sig)
			bsnOutput(map[string]bool{"authenticated": ok}, fmt.Sprintf("%v", ok))
			return nil
		},
	}

	addTxCmd := &cobra.Command{
		Use:   "addtx [addr] [data] [sigHex] [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(8),
		Short: "Securely add a transaction to the mempool",
		RunE: func(cmd *cobra.Command, args []string) error {
			sig, err := hex.DecodeString(args[2])
			if err != nil {
				return fmt.Errorf("invalid signature: %w", err)
			}
			amt, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			fee, err := strconv.ParseUint(args[6], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee: %w", err)
			}
			nonce, err := strconv.ParseUint(args[7], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid nonce: %w", err)
			}
			tx := core.NewTransaction(args[3], args[4], amt, fee, nonce)
			if err := secureNode.SecureAddTransaction(args[0], []byte(args[1]), sig, tx); err != nil {
				return err
			}
			bsnOutput(map[string]string{"status": "queued"}, "queued")
			return nil
		},
	}

	bsnCmd.AddCommand(enrollCmd, removeCmd, authCmd, addTxCmd)
	rootCmd.AddCommand(bsnCmd)
}
