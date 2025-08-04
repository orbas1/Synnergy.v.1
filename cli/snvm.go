package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var vm = core.NewSNVM()

func init() {
	snvmCmd := &cobra.Command{
		Use:   "snvm",
		Short: "Interact with the Synnergy VM",
	}
	execCmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute a no-op transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			tx := core.NewTransaction("", "", 0, 0)
			if err := vm.Execute(tx); err != nil {
				return err
			}
			fmt.Println("executed")
			return nil
		},
	}
	snvmCmd.AddCommand(execCmd)
	rootCmd.AddCommand(snvmCmd)
}
