package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	synnergy "synnergy"
)

var idRegistry = synnergy.NewIDRegistry()

func init() {
	idwCmd := &cobra.Command{
		Use:   "idwallet",
		Short: "ID wallet registry operations",
	}

	registerCmd := &cobra.Command{
		Use:   "register [addr] [info]",
		Args:  cobra.ExactArgs(2),
		Short: "Register a wallet",
		RunE: func(cmd *cobra.Command, args []string) error {
			return idRegistry.Register(args[0], args[1])
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show registration info",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := idRegistry.Info(args[0])
			if !ok {
				fmt.Println("not registered")
				return
			}
			fmt.Println(info)
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if a wallet is registered",
		Run: func(cmd *cobra.Command, args []string) {
			if idRegistry.IsRegistered(args[0]) {
				fmt.Println("registered")
			} else {
				fmt.Println("not registered")
			}
		},
	}

	idwCmd.AddCommand(registerCmd, infoCmd, checkCmd)
	rootCmd.AddCommand(idwCmd)
}
