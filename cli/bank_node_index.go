package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var bankIndex = core.NewBankNodeIndex()

func init() {
	idxCmd := &cobra.Command{Use: "bank_index", Short: "Bank node index"}

	addCmd := &cobra.Command{
		Use:   "add [id] [type]",
		Args:  cobra.ExactArgs(2),
		Short: "Add bank node record to index",
		Run: func(cmd *cobra.Command, args []string) {
			bankIndex.Add(&core.BankNodeRecord{ID: args[0], Type: args[1]})
		},
	}

	var getJSON bool
	getCmd := &cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get bank node details",
		Run: func(cmd *cobra.Command, args []string) {
			rec, ok := bankIndex.Get(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			if getJSON {
				enc, _ := json.Marshal(rec)
				fmt.Println(string(enc))
				return
			}
			fmt.Printf("%s type:%s\n", rec.ID, rec.Type)
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	removeCmd := &cobra.Command{
		Use:   "remove [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove bank node record",
		Run: func(cmd *cobra.Command, args []string) {
			bankIndex.Remove(args[0])
		},
	}

	var listJSON bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List bank node records",
		Run: func(cmd *cobra.Command, args []string) {
			recs := bankIndex.List()
			if listJSON {
				enc, _ := json.Marshal(recs)
				fmt.Println(string(enc))
				return
			}
			for _, r := range recs {
				fmt.Printf("%s type:%s\n", r.ID, r.Type)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	idxCmd.AddCommand(addCmd, getCmd, removeCmd, listCmd)
	rootCmd.AddCommand(idxCmd)
}
