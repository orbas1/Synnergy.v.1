package cli

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

// Stage 40 adds JSON output and error reporting to the biometric template
// management CLI so integrations can process responses reliably.
var (
	biomAuth = core.NewBiometricsAuth()
	bioJSON  bool
)

func bioOutput(v interface{}, plain string) {
	if bioJSON {
		b, err := json.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(plain)
	}
}

func init() {
	authCmd := &cobra.Command{
		Use:   "bioauth",
		Short: "Manage biometric templates",
	}
	authCmd.PersistentFlags().BoolVar(&bioJSON, "json", false, "output results in JSON")

	enrollCmd := &cobra.Command{
		Use:   "enroll [addr] [data] [pubKeyHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Enroll biometric data for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			pubBytes, err := hex.DecodeString(args[2])
			if err != nil || len(pubBytes) != ed25519.PublicKeySize {
				return fmt.Errorf("invalid public key: %w", err)
			}
			biomAuth.Enroll(args[0], []byte(args[1]), ed25519.PublicKey(pubBytes))
			bioOutput(map[string]string{"status": "enrolled"}, "enrolled")
			return nil
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [addr] [data] [sigHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Verify biometric data for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			sig, err := hex.DecodeString(args[2])
			if err != nil {
				return fmt.Errorf("invalid signature: %w", err)
			}
			ok := biomAuth.Verify(args[0], []byte(args[1]), sig)
			bioOutput(map[string]bool{"verified": ok}, fmt.Sprintf("%v", ok))
			return nil
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove biometric data for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			biomAuth.Remove(args[0])
			bioOutput(map[string]string{"status": "removed"}, "removed")
			return nil
		},
	}

	enrolledCmd := &cobra.Command{
		Use:   "enrolled [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if an address has enrolled biometrics",
		RunE: func(cmd *cobra.Command, args []string) error {
			ok := biomAuth.Enrolled(args[0])
			bioOutput(map[string]bool{"enrolled": ok}, fmt.Sprintf("%v", ok))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all enrolled addresses",
		RunE: func(cmd *cobra.Command, args []string) error {
			addrs := biomAuth.List()
			if bioJSON {
				b, err := json.Marshal(map[string][]string{"addresses": addrs})
				if err == nil {
					fmt.Println(string(b))
				}
				return nil
			}
			for _, addr := range addrs {
				fmt.Println(addr)
			}
			return nil
		},
	}

	authCmd.AddCommand(enrollCmd, verifyCmd, removeCmd, enrolledCmd, listCmd)
	rootCmd.AddCommand(authCmd)
}
