package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var assetRegistry = core.NewAssetRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn800_token",
		Short: "Manage SYN800 asset tokens",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <desc> <valuation> <loc> <type> <cert>",
		Args:  cobra.ExactArgs(6),
		Short: "Register an asset",
		Run: func(cmd *cobra.Command, args []string) {
			val, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid valuation")
				return
			}
			if _, err := assetRegistry.Register(args[0], args[1], val, args[3], args[4], args[5]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	updateCmd := &cobra.Command{
		Use:   "update <id> <valuation>",
		Args:  cobra.ExactArgs(2),
		Short: "Update asset valuation",
		Run: func(cmd *cobra.Command, args []string) {
			val, _ := strconv.ParseUint(args[1], 10, 64)
			if err := assetRegistry.UpdateValuation(args[0], val); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Display asset information",
		Run: func(cmd *cobra.Command, args []string) {
			if a, ok := assetRegistry.Get(args[0]); ok {
				b, _ := json.MarshalIndent(a, "", "  ")
				fmt.Println(string(b))
			} else {
				fmt.Println("not found")
			}
		},
	}

	cmd.AddCommand(registerCmd, updateCmd, infoCmd)
	rootCmd.AddCommand(cmd)
}
