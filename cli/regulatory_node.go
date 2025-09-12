package cli

import (
	"errors"
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
			walletPath, _ := cmd.Flags().GetString("wallet")
			password, _ := cmd.Flags().GetString("password")
			if walletPath != "" {
				w, err := core.LoadWallet(walletPath, password)
				if err != nil {
					return err
				}
				if w.Address != args[0] {
					return errors.New("wallet address mismatch")
				}
				regNode.RegisterWallet(w)
				if _, err := w.Sign(&tx); err != nil {
					return err
				}
			}
			if err := regNode.ApproveTransaction(tx); err != nil {
				printOutput(map[string]any{"status": "rejected", "reason": err.Error(), "logs": regNode.Logs(args[0])})
				return nil
			}
			printOutput(map[string]string{"status": "approved"})
			return nil
		},
	}
	approveCmd.Flags().String("wallet", "", "wallet file for signing")
	approveCmd.Flags().String("password", "", "wallet password")

	flagCmd := &cobra.Command{
		Use:   "flag [addr] [reason]",
		Args:  cobra.ExactArgs(2),
		Short: "Flag an address for a reason",
		RunE: func(cmd *cobra.Command, args []string) error {
			gasPrint("RegNodeFlag")
			if err := regNode.FlagEntity(args[0], args[1]); err != nil {
				return err
			}
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

        auditCmd := &cobra.Command{
                Use:   "audit [addr]",
                Args:  cobra.ExactArgs(1),
                Short: "Audit an address for regulatory flags",
                RunE: func(cmd *cobra.Command, args []string) error {
                        gasPrint("RegNodeAudit")
                        logs, flagged := regNode.Audit(args[0])
                        if !flagged {
                                printOutput(map[string]string{"status": "clean"})
                                return nil
                        }
                        printOutput(map[string]any{"status": "flagged", "logs": logs})
                        return nil
                },
        }

        cmd.AddCommand(approveCmd, flagCmd, logsCmd, auditCmd)
	rootCmd.AddCommand(cmd)
}
