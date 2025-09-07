package cli

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, _ := cmd.Flags().GetString("owner")
			name, _ := cmd.Flags().GetString("name")
			if owner == "" || name == "" {
				return errors.New("owner and name are required")
			}
			desc, _ := cmd.Flags().GetString("desc")
			attrs, _ := cmd.Flags().GetString("attrs")
			it := itemRegistry.CreateItem(owner, name, desc, parseAttrs(attrs))
			fmt.Println(it.ItemID)
			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			return itemRegistry.TransferItem(args[0], args[1])
		},
	}
	cmd.AddCommand(transferCmd)

	updateCmd := &cobra.Command{
		Use:   "update-attrs <id> <attrs>",
		Short: "Update item attributes",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return itemRegistry.UpdateAttributes(args[0], parseAttrs(args[1]))
		},
	}
	cmd.AddCommand(updateCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get item info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			it, ok := itemRegistry.GetItem(args[0])
			if !ok {
				return errors.New("item not found")
			}
			fmt.Printf("ID:%s Owner:%s Name:%s Desc:%s\n", it.ItemID, it.Owner, it.Name, it.Description)
			if len(it.Attributes) > 0 {
				for k, v := range it.Attributes {
					fmt.Printf("%s=%s ", k, v)
				}
				fmt.Println()
			}
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List virtual items",
		RunE: func(cmd *cobra.Command, args []string) error {
			items := itemRegistry.ListItems()
			for _, it := range items {
				fmt.Printf("%s %s %s\n", it.ItemID, it.Owner, it.Name)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
