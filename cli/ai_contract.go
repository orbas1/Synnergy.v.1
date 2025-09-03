package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	aiVM       = core.NewSimpleVM()
	baseReg    = core.NewContractRegistry(aiVM)
	aiRegistry = core.NewAIContractRegistry(baseReg)
)

func init() {
	aiVM.Start()

	aiCmd := &cobra.Command{Use: "ai_contract", Short: "AI enhanced contract operations"}

	deployCmd := &cobra.Command{
		Use:   "deploy [wasm_file] [model_hash] [manifest] [gas_limit] [owner]",
		Args:  cobra.ExactArgs(5),
		Short: "Deploy an AI enhanced contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			wasm, err := ioutil.ReadFile(args[0])
			if err != nil {
				return err
			}
			gas, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return err
			}
			addr, err := aiRegistry.DeployAIContract(wasm, args[1], args[2], gas, args[4])
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "contract:", addr)
			return nil
		},
	}

	invokeCmd := &cobra.Command{
		Use:   "invoke [addr] [input_hex] [gas_limit]",
		Args:  cobra.ExactArgs(3),
		Short: "Invoke infer method on AI contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := hex.DecodeString(args[1])
			if err != nil {
				return err
			}
			gas, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}
			out, used, err := aiRegistry.InvokeAIContract(args[0], input, gas)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "output=%s gas=%d\n", hex.EncodeToString(out), used)
			return nil
		},
	}

	modelCmd := &cobra.Command{
		Use:   "model [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show model hash for a contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			h, ok := aiRegistry.ModelHash(args[0])
			if !ok {
				return fmt.Errorf("not found")
			}
			fmt.Fprintln(cmd.OutOrStdout(), h)
			return nil
		},
	}

	var listJSON bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List deployed AI contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if listJSON {
				type item struct {
					Address   string `json:"address"`
					ModelHash string `json:"model_hash"`
				}
				var out []item
				for _, c := range baseReg.List() {
					h, _ := aiRegistry.ModelHash(c.Address)
					out = append(out, item{c.Address, h})
				}
				b, err := json.MarshalIndent(out, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(b))
				return nil
			}
			for _, c := range baseReg.List() {
				h, _ := aiRegistry.ModelHash(c.Address)
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", c.Address, h)
			}
			return nil
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	aiCmd.AddCommand(deployCmd, invokeCmd, modelCmd, listCmd)
	rootCmd.AddCommand(aiCmd)
}
