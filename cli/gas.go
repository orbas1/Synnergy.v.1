package cli

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	synn "synnergy"
)

var (
	gasCmd = &cobra.Command{
		Use:   "gas",
		Short: "Interact with gas table",
	}
)

func init() {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List gas costs",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("GasList")
			entries := synn.GasCatalogue()
			if jsonOutput {
				printOutput(entries)
				return
			}
			buf := &bytes.Buffer{}
			tw := tabwriter.NewWriter(buf, 0, 0, 2, ' ', 0)
			fmt.Fprintln(tw, "OPCODE\tCOST\tCATEGORY\tDESCRIPTION")
			for _, entry := range entries {
				desc := entry.Description
				if desc == "" {
					desc = "-"
				}
				fmt.Fprintf(tw, "%s\t%d\t%s\t%s\n", entry.Name, entry.Cost, entry.Category, desc)
			}
			tw.Flush()
			fmt.Print(buf.String())
		},
	}
	gasCmd.AddCommand(listCmd)
	rootCmd.AddCommand(gasCmd)
}
