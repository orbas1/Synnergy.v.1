package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	synnergy "synnergy"
)

var identitySvc = synnergy.NewIdentityService()

func init() {
	identityCmd := &cobra.Command{
		Use:   "identity",
		Short: "Identity verification operations",
	}

	registerCmd := &cobra.Command{
		Use:   "register [addr] [name] [dob] [nationality]",
		Args:  cobra.ExactArgs(4),
		Short: "Register identity information",
		RunE: func(cmd *cobra.Command, args []string) error {
			return identitySvc.Register(args[0], args[1], args[2], args[3])
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [addr] [method]",
		Args:  cobra.ExactArgs(2),
		Short: "Record a verification method",
		RunE: func(cmd *cobra.Command, args []string) error {
			return identitySvc.Verify(args[0], args[1])
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve identity information",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := identitySvc.Info(args[0])
			if !ok {
				fmt.Println("identity not found")
				return
			}
			fmt.Printf("Name: %s\nDOB: %s\nNationality: %s\n", info.Name, info.DateOfBirth, info.Nationality)
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show verification logs",
		Run: func(cmd *cobra.Command, args []string) {
			logs := identitySvc.Logs(args[0])
			for _, l := range logs {
				fmt.Printf("%s - %s\n", l.Timestamp.Format(time.RFC3339), l.Method)
			}
		},
	}

	identityCmd.AddCommand(registerCmd, verifyCmd, infoCmd, logsCmd)
	rootCmd.AddCommand(identityCmd)
}
