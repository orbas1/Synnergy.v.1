package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var regNode = core.NewRegulatoryNode("regnode1", regManager)

func init() {
	cmd := &cobra.Command{
		Use:   "regnode",
		Short: "Regulatory node operations",
	}

	approveCmd := &cobra.Command{
		Use:   "approve [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Approve or reject a transaction by amount",
		Run: func(cmd *cobra.Command, args []string) {
			amt, _ := strconv.ParseUint(args[0], 10, 64)
			tx := core.Transaction{Amount: amt}
			fmt.Println(regNode.ApproveTransaction(tx))
		},
	}

	flagCmd := &cobra.Command{
		Use:   "flag [addr] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Flag an address for a reason",
		Run:   func(cmd *cobra.Command, args []string) { regNode.FlagEntity(args[0], args[1]) },
	}

	logsCmd := &cobra.Command{
		Use:   "logs [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show logs for an address",
		Run: func(cmd *cobra.Command, args []string) {
			for _, l := range regNode.Logs(args[0]) {
				fmt.Println(l)
			}
		},
	}

	cmd.AddCommand(approveCmd, flagCmd, logsCmd)
	rootCmd.AddCommand(cmd)
}
