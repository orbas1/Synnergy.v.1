package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/core"
)

var compressJSON bool

func compOut(v interface{}, plain string) {
	if compressJSON {
		b, err := json.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(plain)
	}
}

func init() {
	compCmd := &cobra.Command{Use: "compression", Short: "Compressed ledger snapshots"}
	compCmd.PersistentFlags().BoolVar(&compressJSON, "json", false, "output results in JSON")

	saveCmd := &cobra.Command{
		Use:   "save [file]",
		Args:  cobra.ExactArgs(1),
		Short: "Write a compressed ledger snapshot.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := core.SaveCompressedSnapshot(ledger, args[0]); err != nil {
				return err
			}
			compOut(map[string]string{"status": "saved"}, "saved")
			return nil
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
			compOut(map[string]uint64{"height": uint64(h)}, fmt.Sprintf("height: %d", h))
			return nil
		},
	}

	compCmd.AddCommand(saveCmd, loadCmd)
	rootCmd.AddCommand(compCmd)
}
