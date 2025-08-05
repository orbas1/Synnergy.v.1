package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var biomAuth = core.NewBiometricsAuth()

func init() {
	authCmd := &cobra.Command{
		Use:   "bioauth",
		Short: "Manage biometric templates",
	}

	enrollCmd := &cobra.Command{
		Use:   "enroll [addr] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Enroll biometric data for an address",
		Run: func(cmd *cobra.Command, args []string) {
			biomAuth.Enroll(args[0], []byte(args[1]))
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [addr] [data]",
		Args:  cobra.ExactArgs(2),
		Short: "Verify biometric data for an address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(biomAuth.Verify(args[0], []byte(args[1])))
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove biometric data for an address",
		Run: func(cmd *cobra.Command, args []string) {
			biomAuth.Remove(args[0])
		},
	}

	enrolledCmd := &cobra.Command{
		Use:   "enrolled [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if an address has enrolled biometrics",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(biomAuth.Enrolled(args[0]))
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all enrolled addresses",
		Run: func(cmd *cobra.Command, args []string) {
			for _, addr := range biomAuth.List() {
				fmt.Println(addr)
			}
		},
	}

	authCmd.AddCommand(enrollCmd, verifyCmd, removeCmd, enrolledCmd, listCmd)
	rootCmd.AddCommand(authCmd)
}
