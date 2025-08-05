package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var itemRegistry = tokens.NewItemRegistry()

func parseAttrs(input string) map[string]string {
	m := make(map[string]string)
	if input == "" {
		return m
	}
	for _, kv := range strings.Split(input, ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn2369",
		Short: "Virtual item registry",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a virtual item",
		Run: func(cmd *cobra.Command, args []string) {
			owner, _ := cmd.Flags().GetString("owner")
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("desc")
			attrs, _ := cmd.Flags().GetString("attrs")
			it := itemRegistry.CreateItem(owner, name, desc, parseAttrs(attrs))
			fmt.Println(it.ItemID)
		},
	}
	createCmd.Flags().String("owner", "", "owner address")
	createCmd.Flags().String("name", "", "item name")
	createCmd.Flags().String("desc", "", "description")
	createCmd.Flags().String("attrs", "", "key=value attributes")
	cmd.AddCommand(createCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <id> <newOwner>",
		Short: "Transfer item ownership",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := itemRegistry.TransferItem(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(transferCmd)

	updateCmd := &cobra.Command{
		Use:   "update-attrs <id> <attrs>",
		Short: "Update item attributes",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := itemRegistry.UpdateAttributes(args[0], parseAttrs(args[1])); err != nil {
				fmt.Printf("error: %v\n", err)
			}
		},
	}
	cmd.AddCommand(updateCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get item info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			it, ok := itemRegistry.GetItem(args[0])
			if !ok {
				fmt.Println("item not found")
				return
			}
			fmt.Printf("ID:%s Owner:%s Name:%s Desc:%s\n", it.ItemID, it.Owner, it.Name, it.Description)
			if len(it.Attributes) > 0 {
				for k, v := range it.Attributes {
					fmt.Printf("%s=%s ", k, v)
				}
				fmt.Println()
			}
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List virtual items",
		Run: func(cmd *cobra.Command, args []string) {
			items := itemRegistry.ListItems()
			for _, it := range items {
				fmt.Printf("%s %s %s\n", it.ItemID, it.Owner, it.Name)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
