package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	synnergy "synnergy"
)

var regNode = synnergy.NewRegulatoryNode("regnode1", regManager)

func init() {
	nodeCmd := &cobra.Command{Use: "regnode", Short: "Regulatory node operations"}

	approveCmd := &cobra.Command{
		Use:   "approve [from] [to] [amount]",
		Args:  cobra.ExactArgs(3),
		Short: "Approve a transaction against regulations",
		Run: func(cmd *cobra.Command, args []string) {
			amount, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			tx := synnergy.Transaction{From: args[0], To: args[1], Amount: amount}
			if regNode.ApproveTransaction(tx) {
				fmt.Println("approved")
			} else {
				fmt.Println("rejected")
			}
		},
	}

	flagCmd := &cobra.Command{
		Use:   "flag [addr] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Flag an address",
		Run: func(cmd *cobra.Command, args []string) {
			regNode.FlagEntity(args[0], args[1])
		},
	}

	logsCmd := &cobra.Command{
		Use:   "logs [addr]",
		Args:  cobra.ExactArgs(1),
		Short: "Show regulatory flags",
		Run: func(cmd *cobra.Command, args []string) {
			logs := regNode.Logs(args[0])
			for _, l := range logs {
				fmt.Println(l)
			}
		},
	}

	nodeCmd.AddCommand(approveCmd, flagCmd, logsCmd)
	rootCmd.AddCommand(nodeCmd)
}
