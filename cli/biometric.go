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
		Use:   "enroll [userID] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Enroll biometric data for a user",
		Run: func(cmd *cobra.Command, args []string) {
			biometricSvc.Enroll(args[0], []byte(args[1]))
			fmt.Println("biometric enrolled")
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [userID] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Verify biometric data for a user",
		Run: func(cmd *cobra.Command, args []string) {
			ok := biometricSvc.Verify(args[0], []byte(args[1]))
			fmt.Println(ok)
		},
	}

	biometricCmd.AddCommand(enrollCmd, verifyCmd)
	rootCmd.AddCommand(biometricCmd)
}
