package cli

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	opCmd := &cobra.Command{
		Use:   "opcodes",
		Short: "Inspect VM opcodes",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all opcode mappings",
		Run: func(cmd *cobra.Command, args []string) {
			for _, line := range core.DebugDump() {
				fmt.Println(line)
			}
		},
	}

	hexCmd := &cobra.Command{
		Use:   "hex [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Show opcode hex for a function name",
		Run: func(cmd *cobra.Command, args []string) {
			h, err := core.HexDump(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(h)
		},
	}

	bytesCmd := &cobra.Command{
		Use:   "bytes [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Show raw opcode bytes",
		Run: func(cmd *cobra.Command, args []string) {
			b, err := core.ToBytecode(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(hex.EncodeToString(b))
		},
	}

	opCmd.AddCommand(listCmd, hexCmd, bytesCmd)
	rootCmd.AddCommand(opCmd)
}
