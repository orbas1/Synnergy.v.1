package cli

import (
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
		Use:   "approve [from] [amount]",
		Args:  cobra.ExactArgs(2),
		Short: "Approve or reject a transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegNodeApprove")
			amt, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			tx := core.Transaction{From: args[0], Amount: amt}
			if regNode.ApproveTransaction(tx) {
				printOutput(map[string]string{"status": "approved"})
			} else {
				printOutput(map[string]any{"status": "rejected", "logs": regNode.Logs(args[0])})
			}
			return nil
		},
	}

	flagCmd := &cobra.Command{
		Use:   "flag [addr] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Flag an address for a reason",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegNodeFlag")
			regNode.FlagEntity(args[0], args[1])
			printOutput(map[string]any{"status": "flagged", "address": args[0]})
			return nil
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show logs for an address",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegNodeLogs")
			printOutput(regNode.Logs(args[0]))
			return nil
		},
	}

	cmd.AddCommand(approveCmd, flagCmd, logsCmd)
	rootCmd.AddCommand(cmd)
}
