package cli

import (
	"encoding/hex"
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
			gasPrint("OpcodesList")
			printOutput(core.Catalogue())
		},
	}

	hexCmd := &cobra.Command{
		Use:   "hex [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Show opcode hex for a function name",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("OpcodesHex")
			h, err := core.HexDump(args[0])
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]string{"hex": h})
		},
	}

	bytesCmd := &cobra.Command{
		Use:   "bytes [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Show raw opcode bytes",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("OpcodesBytes")
			b, err := core.ToBytecode(args[0])
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			printOutput(map[string]string{"bytes": hex.EncodeToString(b)})
		},
	}

	opCmd.AddCommand(listCmd, hexCmd, bytesCmd)
	rootCmd.AddCommand(opCmd)
}
