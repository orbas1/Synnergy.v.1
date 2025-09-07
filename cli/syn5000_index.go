package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var gamblingIndex = map[string]core.GamblingToken{}

func init() {
	cmd := &cobra.Command{
		Use:   "syn5000_index",
		Short: "Gambling token index",
	}

	addCmd := &cobra.Command{
		Use:   "add <symbol>",
		Args:  cobra.ExactArgs(1),
		Short: "Add a new SYN5000 token to index",
		RunE: func(cmd *cobra.Command, args []string) error {
			sym := args[0]
			if sym == "" {
				return fmt.Errorf("symbol required")
			}
			gamblingIndex[sym] = core.NewSYN5000Token(sym, sym, 0)
			cmd.Println("token added")
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tokens in index",
		Run: func(cmd *cobra.Command, args []string) {
			for sym := range gamblingIndex {
				cmd.Println(sym)
			}
		},
	}

	cmd.AddCommand(addCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
