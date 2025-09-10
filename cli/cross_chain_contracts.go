package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	synnergy "synnergy"
	"synnergy/core"
)

var crossContractRegistry = core.NewCrossChainRegistry()

func init() {
	crossContractRegistry.AuthorizeRelayer("cli_relayer")
	cmd := &cobra.Command{
		Use:   "xcontract",
		Short: "Register cross-chain contract mappings",
	}

	var listJSON bool
	var getJSON bool
	var registerJSON bool
	var removeJSON bool

	registerCmd := &cobra.Command{
		Use:   "register <local_addr> <remote_chain> <remote_addr>",
		Args:  cobra.ExactArgs(3),
		Short: "Register a contract mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := crossContractRegistry.RegisterMapping("cli_relayer", args[0], args[1], args[2]); err != nil {
				return err
			}
			gas := synnergy.GasCost("RegisterXContract")
			if registerJSON {
				enc, _ := json.Marshal(map[string]interface{}{"gas": gas})
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("gas:%d\n", gas)
			return nil
		},
	}
	registerCmd.Flags().BoolVar(&registerJSON, "json", false, "output as JSON")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered mappings",
		Run: func(cmd *cobra.Command, args []string) {
			maps := crossContractRegistry.ListMappings()
			if listJSON {
				enc, _ := json.Marshal(maps)
				fmt.Println(string(enc))
				return
			}
			for _, m := range maps {
				fmt.Printf("%s -> %s:%s\n", m.LocalAddress, m.RemoteChain, m.RemoteAddress)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	getCmd := &cobra.Command{
		Use:   "get <local_addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve mapping info",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, ok := crossContractRegistry.GetMapping(args[0])
			if !ok {
				return fmt.Errorf("mapping not found")
			}
			if getJSON {
				enc, _ := json.Marshal(m)
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("%s -> %s:%s\n", m.LocalAddress, m.RemoteChain, m.RemoteAddress)
			return nil
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	removeCmd := &cobra.Command{
		Use:   "remove <local_addr>",
		Args:  cobra.ExactArgs(1),
		Short: "Delete a mapping",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := crossContractRegistry.RemoveMapping("cli_relayer", args[0]); err != nil {
				return err
			}
			gas := synnergy.GasCost("RemoveXContract")
			if removeJSON {
				enc, _ := json.Marshal(map[string]interface{}{"gas": gas})
				fmt.Println(string(enc))
				return nil
			}
			fmt.Printf("gas:%d\n", gas)
			return nil
		},
	}
	removeCmd.Flags().BoolVar(&removeJSON, "json", false, "output as JSON")

	cmd.AddCommand(registerCmd, listCmd, getCmd, removeCmd)
	rootCmd.AddCommand(cmd)
}
