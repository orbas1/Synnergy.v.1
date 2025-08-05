package cli

import (
	"encoding/hex"
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
		Run: func(cmd *cobra.Command, args []string) {
			wasm, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Println("read error:", err)
				return
			}
			gas, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				fmt.Println("gas error:", err)
				return
			}
			addr, err := aiRegistry.DeployAIContract(wasm, args[1], args[2], gas, args[4])
			if err != nil {
				fmt.Println("deploy error:", err)
				return
			}
			fmt.Println("contract:", addr)
		},
	}

	invokeCmd := &cobra.Command{
		Use:   "invoke [addr] [input_hex] [gas_limit]",
		Args:  cobra.ExactArgs(3),
		Short: "Invoke infer method on AI contract",
		Run: func(cmd *cobra.Command, args []string) {
			input, err := hex.DecodeString(args[1])
			if err != nil {
				fmt.Println("input error:", err)
				return
			}
			gas, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("gas error:", err)
				return
			}
			out, used, err := aiRegistry.InvokeAIContract(args[0], input, gas)
			if err != nil {
				fmt.Println("invoke error:", err)
				return
			}
			fmt.Printf("output=%s gas=%d\n", hex.EncodeToString(out), used)
		},
	}

	modelCmd := &cobra.Command{
		Use:   "model [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show model hash for a contract",
		Run: func(cmd *cobra.Command, args []string) {
			h, ok := aiRegistry.ModelHash(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Println(h)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List deployed AI contracts",
		Run: func(cmd *cobra.Command, args []string) {
			for _, c := range baseReg.List() {
				h, _ := aiRegistry.ModelHash(c.Address)
				fmt.Printf("%s %s\n", c.Address, h)
			}
		},
	}

	aiCmd.AddCommand(deployCmd, invokeCmd, modelCmd, listCmd)
	rootCmd.AddCommand(aiCmd)
}
