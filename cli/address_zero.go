package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	zeroCmd := &cobra.Command{Use: "addrzero", Short: "Zero address utilities"}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Display the zero address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(core.AddressZero)
		},
	}

	isCmd := &cobra.Command{
		Use:   "is [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if address is zero address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(core.IsZeroAddress(args[0]))
		},
	}

	zeroCmd.AddCommand(showCmd, isCmd)
	rootCmd.AddCommand(zeroCmd)
}
