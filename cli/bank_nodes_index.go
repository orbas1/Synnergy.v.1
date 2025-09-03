package cli

import (
        "encoding/json"
        "fmt"

        "github.com/spf13/cobra"
        "synnergy/core"
)

func init() {
	bankNodesCmd := &cobra.Command{
		Use:   "banknodes",
		Short: "Bank node type utilities",
	}

        var jsonOut bool
        typesCmd := &cobra.Command{
                Use:   "types",
                Short: "List supported bank node types",
                Run: func(cmd *cobra.Command, args []string) {
                        if jsonOut {
                                enc, _ := json.Marshal(core.BankNodeTypes)
                                fmt.Println(string(enc))
                                return
                        }
                        for _, t := range core.BankNodeTypes {
                                fmt.Println(t)
                        }
                },
        }
        typesCmd.Flags().BoolVar(&jsonOut, "json", false, "output as JSON")

	bankNodesCmd.AddCommand(typesCmd)
	rootCmd.AddCommand(bankNodesCmd)
}
