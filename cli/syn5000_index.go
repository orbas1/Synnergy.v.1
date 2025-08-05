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
		Run: func(cmd *cobra.Command, args []string) {
			gamblingIndex[args[0]] = core.NewSYN5000Token(args[0], args[0], 0)
			fmt.Println("token added")
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List tokens in index",
		Run: func(cmd *cobra.Command, args []string) {
			for sym := range gamblingIndex {
				fmt.Println(sym)
			}
		},
	}

	cmd.AddCommand(addCmd, listCmd)
	rootCmd.AddCommand(cmd)
}
