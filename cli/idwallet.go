package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var idRegistry = core.NewIDRegistry()

func init() {
	idwCmd := &cobra.Command{
		Use:   "idwallet",
		Short: "ID wallet registration",
	}

	registerCmd := &cobra.Command{
		Use:   "register [address] [info]",
		Args:  cobra.ExactArgs(2),
		Short: "Register an ID wallet",
		Run: func(cmd *cobra.Command, args []string) {
			if err := idRegistry.Register(args[0], args[1]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	checkCmd := &cobra.Command{
		Use:   "check [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Check registration info",
		Run: func(cmd *cobra.Command, args []string) {
			info, ok := idRegistry.Info(args[0])
			if !ok {
				fmt.Println("not registered")
				return
			}
			fmt.Println(info)
		},
	}

	idwCmd.AddCommand(registerCmd, checkCmd)
	rootCmd.AddCommand(idwCmd)
}
