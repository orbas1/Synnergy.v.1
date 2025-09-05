package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	// Stage 38 ensures address utilities handle bytes and shortened forms consistently.
	addressCmd := &cobra.Command{Use: "address", Short: "Address utilities"}

	parseCmd := &cobra.Command{
		Use:   "parse [hex]",
		Args:  cobra.ExactArgs(1),
		Short: "Validate and normalise an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), addr.Hex())
			return nil
		},
	}

	bytesCmd := &cobra.Command{
		Use:   "bytes [hex]",
		Args:  cobra.ExactArgs(1),
		Short: "Show raw bytes of an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), hex.EncodeToString(addr.Bytes()))
			return nil
		},
	}

	shortCmd := &cobra.Command{
		Use:   "short [hex]",
		Args:  cobra.ExactArgs(1),
		Short: "Show shortened form of an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), addr.Short())
			return nil
		},
	}

	addressCmd.AddCommand(parseCmd, bytesCmd, shortCmd)
	rootCmd.AddCommand(addressCmd)
}
