package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var authorityIndex = core.NewAuthorityNodeIndex()

func init() {
	idxCmd := &cobra.Command{Use: "authority_index", Short: "Authority node index"}

	addCmd := &cobra.Command{
		Use:   "add [address] [role]",
		Args:  cobra.ExactArgs(2),
		Short: "Add authority node to index",
		Run: func(cmd *cobra.Command, args []string) {
			node := &core.AuthorityNode{Address: args[0], Role: args[1], Votes: make(map[string]bool)}
			authorityIndex.Add(node)
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Get authority node details",
		Run: func(cmd *cobra.Command, args []string) {
			n, ok := authorityIndex.Get(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%s role:%s votes:%d\n", n.Address, n.Role, len(n.Votes))
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove authority node from index",
		Run: func(cmd *cobra.Command, args []string) {
			authorityIndex.Remove(args[0])
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List authority nodes",
		Run: func(cmd *cobra.Command, args []string) {
			for _, n := range authorityIndex.List() {
				fmt.Printf("%s role:%s votes:%d\n", n.Address, n.Role, len(n.Votes))
			}
		},
	}

	idxCmd.AddCommand(addCmd, getCmd, removeCmd, listCmd)
	rootCmd.AddCommand(idxCmd)
}
