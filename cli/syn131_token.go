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
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			owner, _ := cmd.Flags().GetString("owner")
			val, _ := cmd.Flags().GetUint64("valuation")
			if id == "" || name == "" || symbol == "" || owner == "" {
				return fmt.Errorf("id, name, symbol and owner must be provided")
			}
			if _, err := syn131.Create(id, name, symbol, owner, val); err != nil {
				return err
			}
			cmd.Println("token created")
			return nil
		},
	}
	createCmd.Flags().String("id", "", "token id")
	createCmd.Flags().String("name", "", "name")
	createCmd.Flags().String("symbol", "", "symbol")
	createCmd.Flags().String("owner", "", "owner")
	createCmd.Flags().Uint64("valuation", 0, "valuation")
	createCmd.MarkFlagRequired("id")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")
	createCmd.MarkFlagRequired("owner")
	cmd.AddCommand(createCmd)

	valCmd := &cobra.Command{
		Use:   "value <id> <valuation>",
		Args:  cobra.ExactArgs(2),
		Short: "Update valuation",
		RunE: func(cmd *cobra.Command, args []string) error {
			val, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid valuation")
			}
			if err := syn131.UpdateValuation(args[0], val); err != nil {
				return err
			}
			cmd.Println("valuation updated")
			return nil
		},
	}
	cmd.AddCommand(valCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get token info",
		RunE: func(cmd *cobra.Command, args []string) error {
			tok, ok := syn131.Get(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			cmd.Printf("%s %s %s owner:%s val:%d\n", tok.ID, tok.Name, tok.Symbol, tok.Owner, tok.Valuation)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
