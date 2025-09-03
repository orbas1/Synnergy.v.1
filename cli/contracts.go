package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	contractVM       = core.NewSimpleVM()
	contractRegistry = core.NewContractRegistry(contractVM)
)

func init() {
	// start VM to allow contract execution
	_ = contractVM.Start()

	contractsCmd := &cobra.Command{
		Use:   "contracts",
		Short: "Compile, deploy and invoke smart contracts",
	}

	compileCmd := &cobra.Command{
		Use:   "compile [src.wat|src.wasm]",
		Args:  cobra.ExactArgs(1),
		Short: "Compile WAT or WASM to deterministic bytecode",
		RunE: func(cmd *cobra.Command, args []string) error {
			src, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			_, hash, err := core.CompileWASM(src)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), hash)
			return nil
		},
	}

	var wasmPath, manifestPath, owner string
	var gasLimit uint64

	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy compiled WASM",
		RunE: func(cmd *cobra.Command, args []string) error {
			if wasmPath == "" {
				return fmt.Errorf("--wasm required")
			}
			wasm, err := os.ReadFile(wasmPath)
			if err != nil {
				return err
			}
			var manifest string
			if manifestPath != "" {
				m, err := os.ReadFile(manifestPath)
				if err != nil {
					return err
				}
				manifest = string(m)
			}
			addr, err := contractRegistry.Deploy(wasm, manifest, gasLimit, owner)
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), addr)
			return nil
		},
	}
	deployCmd.Flags().StringVar(&wasmPath, "wasm", "", "Path to compiled WASM")
	deployCmd.Flags().StringVar(&manifestPath, "ric", "", "Path to Ricardian manifest")
	deployCmd.Flags().Uint64Var(&gasLimit, "gas", 100000, "Gas limit")
	deployCmd.Flags().StringVar(&owner, "owner", "", "Owner address")

	var invokeMethod, invokeArgs string
	var invokeGas uint64
	invokeCmd := &cobra.Command{
		Use:   "invoke <address>",
		Args:  cobra.ExactArgs(1),
		Short: "Invoke a contract method",
		RunE: func(cmd *cobra.Command, args []string) error {
			addr := args[0]
			out, gas, err := contractRegistry.Invoke(addr, invokeMethod, []byte(invokeArgs), invokeGas)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "output: %s\ngas: %d\n", string(out), gas)
			return nil
		},
	}
	invokeCmd.Flags().StringVar(&invokeMethod, "method", "", "Contract method to call")
	invokeCmd.Flags().StringVar(&invokeArgs, "args", "", "Arguments as raw bytes")
	invokeCmd.Flags().Uint64Var(&invokeGas, "gas", 0, "Gas limit (0 for default)")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List deployed contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, c := range contractRegistry.List() {
				fmt.Fprintf(cmd.OutOrStdout(), "%s owner=%s gas=%d paused=%v\n", c.Address, c.Owner, c.GasLimit, c.Paused)
			}
			return nil
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <address>",
		Args:  cobra.ExactArgs(1),
		Short: "Show contract manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, ok := contractRegistry.Get(args[0])
			if !ok {
				return fmt.Errorf("contract not found")
			}
			fmt.Fprintln(cmd.OutOrStdout(), c.Manifest)
			return nil
		},
	}

	contractsCmd.AddCommand(compileCmd, deployCmd, invokeCmd, listCmd, infoCmd)
	rootCmd.AddCommand(contractsCmd)
}
