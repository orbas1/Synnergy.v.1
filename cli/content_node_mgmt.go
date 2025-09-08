package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
	ilog "synnergy/internal/log"
)

var (
	contentNode = core.NewContentNetworkNode("node1", "addr")
	cnJSON      bool
)

func cnOutput(v interface{}, plain string) {
	if cnJSON {
		b, err := json.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(plain)
	}
}

func init() {
	cnCmd := &cobra.Command{Use: "content_node", Short: "Manage content registry"}
	cnCmd.PersistentFlags().BoolVar(&cnJSON, "json", false, "output json")

	registerCmd := &cobra.Command{Use: "register [id] [name]", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		meta := core.NewContentMeta(args[0], args[1], int64(len(args[1])), args[0])
		if err := contentNode.Register(meta); err != nil {
			ilog.Error("cli_cn_register", "error", err)
			return err
		}
		ilog.Info("cli_cn_register", "id", meta.ID)
		cnOutput(map[string]string{"status": "registered"}, "registered")
		return nil
	}}

	unregisterCmd := &cobra.Command{Use: "unregister [id]", Args: cobra.ExactArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		if err := contentNode.Unregister(args[0]); err != nil {
			ilog.Error("cli_cn_unregister", "error", err)
			return err
		}
		ilog.Info("cli_cn_unregister", "id", args[0])
		cnOutput(map[string]string{"status": "unregistered"}, "unregistered")
		return nil
	}}

	listCmd := &cobra.Command{Use: "list", Run: func(cmd *cobra.Command, args []string) {
		list := contentNode.List()
		ilog.Info("cli_cn_list", "count", len(list))
		if cnJSON {
			b, err := json.Marshal(list)
			if err == nil {
				fmt.Println(string(b))
			}
		} else {
			for _, m := range list {
				fmt.Println(m.ID)
			}
		}
	}}

	cnCmd.AddCommand(registerCmd, unregisterCmd, listCmd)
	rootCmd.AddCommand(cnCmd)
}
