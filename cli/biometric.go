package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

// biometricSvc is shared across CLI commands to provide enrollment and
// verification capabilities.
var biometricSvc = core.NewBiometricService()

func init() {
	biometricCmd := &cobra.Command{
		Use:   "biometric",
		Short: "Biometric authentication operations",
	}

	enrollCmd := &cobra.Command{
		Use:   "enroll [userID] [data] [pubKeyHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Enroll biometric data for a user",
		Run: func(cmd *cobra.Command, args []string) {
			pub, err := parsePubKey(args[2])
			if err != nil {
				fmt.Println("invalid public key:", err)
				return
			}
			biometricSvc.Enroll(args[0], []byte(args[1]), pub)
			fmt.Println("biometric enrolled")
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [userID] [data] [sigHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Verify biometric data for a user",
		Run: func(cmd *cobra.Command, args []string) {
			sig, err := decodeSig(args[2])
			if err != nil {
				fmt.Println("invalid signature:", err)
				return
			}
			ok := biometricSvc.Verify(args[0], []byte(args[1]), sig)
			fmt.Println(ok)
		},
	}

	biometricCmd.AddCommand(enrollCmd, verifyCmd)
	rootCmd.AddCommand(biometricCmd)
}
