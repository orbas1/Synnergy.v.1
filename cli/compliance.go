package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
	ilog "synnergy/internal/log"
)

var (
	compliance = core.NewComplianceService()
	compJSON   bool
)

func compOutput(v interface{}, plain string) {
	if compJSON {
		b, err := json.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(plain)
	}
}

func init() {
	compCmd := &cobra.Command{Use: "compliance", Short: "Compliance operations"}
	compCmd.PersistentFlags().BoolVar(&compJSON, "json", false, "output results in JSON")

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
			ilog.Info("cli_validate_kyc", "address", input.Address)
			compOutput(map[string]string{"commitment": commit}, fmt.Sprintf("commitment: %s", commit))
			return nil
		},
	}

	eraseCmd := &cobra.Command{
		Use:   "erase [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a user's KYC data.",
		RunE: func(cmd *cobra.Command, args []string) error {
			compliance.EraseKYC(args[0])
			ilog.Info("cli_erase_kyc", "address", args[0])
			compOutput(map[string]string{"status": "erased"}, "erased")
			return nil
		},
	}

	fraudCmd := &cobra.Command{
		Use:   "fraud [address] [severity]",
		Args:  cobra.ExactArgs(2),
		Short: "Record a fraud signal.",
		RunE: func(cmd *cobra.Command, args []string) error {
			sev, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid severity: %w", err)
			}
			compliance.RecordFraud(args[0], sev)
			ilog.Info("cli_fraud", "address", args[0], "severity", sev)
			compOutput(map[string]string{"status": "recorded"}, "recorded")
			return nil
		},
	}

	riskCmd := &cobra.Command{
		Use:   "risk [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Retrieve accumulated fraud risk score.",
		RunE: func(cmd *cobra.Command, args []string) error {
			score := compliance.RiskScore(args[0])
			ilog.Info("cli_risk", "address", args[0], "score", score)
			compOutput(map[string]int{"risk": score}, fmt.Sprintf("%d", score))
			return nil
		},
	}

	auditCmd := &cobra.Command{
		Use:   "audit [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Display the audit trail for an address.",
		RunE: func(cmd *cobra.Command, args []string) error {
			events := compliance.AuditTrail(args[0])
			ilog.Info("cli_audit", "address", args[0])
			if compJSON {
				b, err := json.Marshal(events)
				if err == nil {
					fmt.Println(string(b))
				}
				return nil
			}
			for _, e := range events {
				b, _ := json.Marshal(e)
				fmt.Println(string(b))
			}
			return nil
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
			thr, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				return fmt.Errorf("invalid threshold: %w", err)
			}
			anomaly := compliance.MonitorTransaction(tx, thr)
			ilog.Info("cli_monitor", "id", tx.ID, "anomaly", anomaly)
			if anomaly {
				compOutput(map[string]string{"status": "anomaly"}, "anomaly detected")
			} else {
				compOutput(map[string]string{"status": "ok"}, "ok")
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
			ilog.Info("cli_verify_zkp", "result", ok)
			compOutput(map[string]bool{"valid": ok}, fmt.Sprintf("%v", ok))
			return nil
		},
	}

	compCmd.AddCommand(validateCmd, eraseCmd, fraudCmd, riskCmd, auditCmd, monitorCmd, zkpCmd)
	rootCmd.AddCommand(compCmd)
}
