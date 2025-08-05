package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	compCmd := &cobra.Command{Use: "compression", Short: "Compressed ledger snapshots"}

	saveCmd := &cobra.Command{
		Use:   "save [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Write a compressed ledger snapshot.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.SaveCompressedSnapshot(ledger, args[0])
		},
	}

	loadCmd := &cobra.Command{
		Use:   "load [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Load a compressed snapshot and display the height.",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := core.LoadCompressedSnapshot(args[0])
			if err != nil {
				return err
			}
			ledger = l
			h, _ := ledger.Head()
			fmt.Println("height:", h)
			return nil
		},
	}

	compCmd.AddCommand(saveCmd, loadCmd)
	rootCmd.AddCommand(compCmd)
}
