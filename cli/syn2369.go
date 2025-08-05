package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var itemReg = tokens.NewItemRegistry()

func parseAttrMap(input string) map[string]string {
	m := make(map[string]string)
	if input == "" {
		return m
	}
	for _, part := range strings.Split(input, ",") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
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
		Short: "Create virtual item",
		Run: func(cmd *cobra.Command, args []string) {
			owner, _ := cmd.Flags().GetString("owner")
			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("desc")
			attrs, _ := cmd.Flags().GetString("attrs")
			it := itemReg.CreateItem(owner, name, desc, parseAttrMap(attrs))
			fmt.Println(it.ItemID)
		},
	}
	createCmd.Flags().String("owner", "", "item owner")
	createCmd.Flags().String("name", "", "item name")
	createCmd.Flags().String("desc", "", "description")
	createCmd.Flags().String("attrs", "", "attributes key=value pairs")
	cmd.AddCommand(createCmd)

	transferCmd := &cobra.Command{
		Use:   "transfer <itemID> <newOwner>",
		Short: "Transfer item",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := itemReg.TransferItem(args[0], args[1]); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("transferred")
			}
		},
	}
	cmd.AddCommand(transferCmd)

	attrsCmd := &cobra.Command{
		Use:   "attrs <itemID> <key=value,...>",
		Short: "Update item attributes",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := itemReg.UpdateAttributes(args[0], parseAttrMap(args[1])); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("updated")
			}
		},
	}
	cmd.AddCommand(attrsCmd)

	getCmd := &cobra.Command{
		Use:   "get <itemID>",
		Short: "Get item info",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			it, ok := itemReg.GetItem(args[0])
			if !ok {
				fmt.Println("item not found")
				return
			}
			fmt.Printf("%+v\n", *it)
		},
	}
	cmd.AddCommand(getCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List items",
		Run: func(cmd *cobra.Command, args []string) {
			for _, it := range itemReg.ListItems() {
				fmt.Printf("%+v\n", *it)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
