package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var protocolRegistry = core.NewProtocolRegistry()

func init() {
	cmd := &cobra.Command{
		Use:   "cross_chain_agnostic_protocols",
		Short: "Register cross-chain protocols",
	}

	var listJSON bool
	var getJSON bool

	registerCmd := &cobra.Command{
		Use:   "register <name>",
		Args:  cobra.ExactArgs(1),
		Short: "Register a new protocol definition",
		Run: func(cmd *cobra.Command, args []string) {
			id := protocolRegistry.Register(args[0])
			fmt.Printf("%d gas:%d\n", id, synnergy.GasCost("RegisterProtocol"))
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered protocols",
		Run: func(cmd *cobra.Command, args []string) {
			ps := protocolRegistry.List()
			if listJSON {
				enc, _ := json.Marshal(ps)
				fmt.Println(string(enc))
				return
			}
			for _, p := range ps {
				fmt.Printf("%d: %s\n", p.ID, p.Name)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve a protocol configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			p, ok := protocolRegistry.Get(id)
			if !ok {
				return fmt.Errorf("protocol not found")
			}
			if getJSON {
				enc, _ := json.Marshal(p)
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("%d: %s\n", p.ID, p.Name)
			return nil
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	cmd.AddCommand(registerCmd, listCmd, getCmd)
	rootCmd.AddCommand(cmd)
}
