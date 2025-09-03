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
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), core.AddressZero)
			return nil
		},
	}

	isCmd := &cobra.Command{
		Use:   "is [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if address is zero address",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), core.IsZeroAddress(args[0]))
			return nil
		},
	}

	zeroCmd.AddCommand(showCmd, isCmd)
	rootCmd.AddCommand(zeroCmd)
}
