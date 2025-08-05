package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var compliance = core.NewComplianceService()

func init() {
	compCmd := &cobra.Command{Use: "compliance", Short: "Compliance operations"}

	validateCmd := &cobra.Command{
		Use:   "validate [kyc.json]",
		Args:  cobra.ExactArgs(1),
		Short: "Validate and store a KYC document commitment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var input struct {
				Address string `json:"address"`
				Data    string `json:"data"`
			}
			b, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			if err := json.Unmarshal(b, &input); err != nil {
				return err
			}
			commit, err := compliance.ValidateKYC(input.Address, []byte(input.Data))
			if err != nil {
				return err
			}
			fmt.Println("commitment:", commit)
			return nil
		},
	}

	eraseCmd := &cobra.Command{
		Use:   "erase [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a user's KYC data.",
		Run: func(cmd *cobra.Command, args []string) {
			compliance.EraseKYC(args[0])
		},
	}

	fraudCmd := &cobra.Command{
		Use:   "fraud [address] [severity]",
		Args:  cobra.ExactArgs(2),
		Short: "Record a fraud signal.",
		Run: func(cmd *cobra.Command, args []string) {
			sev, _ := strconv.Atoi(args[1])
			compliance.RecordFraud(args[0], sev)
		},
	}

	riskCmd := &cobra.Command{
		Use:   "risk [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve accumulated fraud risk score.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(compliance.RiskScore(args[0]))
		},
	}

	auditCmd := &cobra.Command{
		Use:   "audit [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Display the audit trail for an address.",
		Run: func(cmd *cobra.Command, args []string) {
			for _, e := range compliance.AuditTrail(args[0]) {
				b, _ := json.Marshal(e)
				fmt.Println(string(b))
			}
		},
	}

	monitorCmd := &cobra.Command{
		Use:   "monitor [tx.json] [threshold]",
		Args:  cobra.ExactArgs(2),
		Short: "Run anomaly detection on a transaction.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var tx core.ComplianceTransaction
			b, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			if err := json.Unmarshal(b, &tx); err != nil {
				return err
			}
			thr, _ := strconv.ParseFloat(args[1], 64)
			if compliance.MonitorTransaction(tx, thr) {
				fmt.Println("anomaly detected")
			} else {
				fmt.Println("ok")
			}
			return nil
		},
	}

	zkpCmd := &cobra.Command{
		Use:   "verifyzkp [blob.bin] [commitmentHex] [proofHex]",
		Args:  cobra.ExactArgs(3),
		Short: "Verify a zero-knowledge proof.",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}
			ok := compliance.VerifyZKP(b, args[1], args[2])
			fmt.Println(ok)
			return nil
		},
	}

	compCmd.AddCommand(validateCmd, eraseCmd, fraudCmd, riskCmd, auditCmd, monitorCmd, zkpCmd)
	rootCmd.AddCommand(compCmd)
}
