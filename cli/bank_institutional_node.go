package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var bankInstNode = core.NewBankInstitutionalNode("bank1", "addr1", ledger)

func init() {
	bankCmd := &cobra.Command{
		Use:   "bankinst",
		Short: "Bank institutional node operations",
	}

	regCmd := &cobra.Command{
		Use:   "register [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Register a participating institution",
		Run: func(cmd *cobra.Command, args []string) {
			bankInstNode.RegisterInstitution(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered institutions",
		Run: func(cmd *cobra.Command, args []string) {
			for name := range bankInstNode.Institutions {
				fmt.Println(name)
			}
		},
	}

	isCmd := &cobra.Command{
		Use:   "is [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if an institution is registered",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(bankInstNode.IsRegistered(args[0]))
		},
	}

	bankCmd.AddCommand(regCmd, listCmd, isCmd)
	rootCmd.AddCommand(bankCmd)
}
