package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	addressCmd := &cobra.Command{Use: "address", Short: "Address utilities"}

	parseCmd := &cobra.Command{
		Use:   "parse [hex]",
		Args:  cobra.ExactArgs(1),
		Short: "Validate and normalise an address",
		Run: func(cmd *cobra.Command, args []string) {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(addr.Hex())
		},
	}

	bytesCmd := &cobra.Command{
		Use:   "bytes [hex]",
		Args:  cobra.ExactArgs(1),
		Short: "Show raw bytes of an address",
		Run: func(cmd *cobra.Command, args []string) {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(hex.EncodeToString(addr.Bytes()))
		},
	}

	shortCmd := &cobra.Command{
		Use:   "short [hex]",
		Args:  cobra.ExactArgs(1),
		Short: "Show shortened form of an address",
		Run: func(cmd *cobra.Command, args []string) {
			addr, err := core.StringToAddress(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(addr.Short())
		},
	}

	addressCmd.AddCommand(parseCmd, bytesCmd, shortCmd)
	rootCmd.AddCommand(addressCmd)
}
