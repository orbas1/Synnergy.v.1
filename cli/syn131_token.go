package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn131 = core.NewSYN131Registry()

func init() {
	cmd := &cobra.Command{
		Use:   "syn131",
		Short: "SYN131 intangible asset tokens",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new token",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			val, _ := cmd.Flags().GetUint64("valuation")
			if _, err := syn131.Create(id, name, symbol, owner, val); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("token created")
		},
	}
	createCmd.Flags().String("id", "", "token id")
	createCmd.Flags().String("name", "", "name")
	createCmd.Flags().String("symbol", "", "symbol")
	createCmd.Flags().String("owner", "", "owner")
	createCmd.Flags().Uint64("valuation", 0, "valuation")
	cmd.AddCommand(createCmd)

	valCmd := &cobra.Command{
		Use:   "value <id> <valuation>",
		Args:  cobra.ExactArgs(2),
		Short: "Update valuation",
		Run: func(cmd *cobra.Command, args []string) {
			val, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				fmt.Println("invalid valuation")
				return
			}
			if err := syn131.UpdateValuation(args[0], val); err != nil {
				fmt.Println(err)
			}
		},
	}
	cmd.AddCommand(valCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get token info",
		Run: func(cmd *cobra.Command, args []string) {
			tok, ok := syn131.Get(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%s %s %s owner:%s val:%d\n", tok.ID, tok.Name, tok.Symbol, tok.Owner, tok.Valuation)
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
