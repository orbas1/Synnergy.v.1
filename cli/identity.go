package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var identitySvc = core.NewIdentityService()

func init() {
	idCmd := &cobra.Command{
		Use:   "identity",
		Short: "Identity verification service",
	}

	registerCmd := &cobra.Command{
		Use:   "register [addr] [name] [dob] [nationality]",
		Args:  cobra.ExactArgs(4),
		Short: "Register identity information",
		Run: func(cmd *cobra.Command, args []string) {
			if err := identitySvc.Register(args[0], args[1], args[2], args[3]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	verifyCmd := &cobra.Command{
		Use:   "verify [addr] [method]",
		Args:  cobra.ExactArgs(2),
		Short: "Record a verification method",
		Run: func(cmd *cobra.Command, args []string) {
			if err := identitySvc.Verify(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show identity information",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := identitySvc.Info(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("Name: %s DOB: %s Nationality: %s\n", info.Name, info.DateOfBirth, info.Nationality)
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show verification logs",
		Run: func(cmd *cobra.Command, args []string) {
			logs := identitySvc.Logs(args[0])
			for _, l := range logs {
				fmt.Printf("%s %s\n", l.Method, l.Timestamp.Format("2006-01-02T15:04:05"))
			}
		},
	}

	idCmd.AddCommand(registerCmd, verifyCmd, infoCmd, logsCmd)
	rootCmd.AddCommand(idCmd)
}
