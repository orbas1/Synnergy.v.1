package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"synnergy/core"
	ilog "synnergy/internal/log"
)

var (
	complianceMgr = core.NewComplianceManager()
	cmJSON        bool
)

func cmOutput(v interface{}, plain string) {
	if cmJSON {
		b, err := json.Marshal(v)
		if err == nil {
			fmt.Println(string(b))
		}
	} else {
		fmt.Println(plain)
	}
}

func init() {
	mgmtCmd := &cobra.Command{Use: "compliance_management", Short: "Manage address compliance"}
	mgmtCmd.PersistentFlags().BoolVar(&cmJSON, "json", false, "output results in JSON")

	suspendCmd := &cobra.Command{Use: "suspend [addr]", Args: cobra.ExactArgs(1), Short: "Suspend an address", RunE: func(cmd *cobra.Command, args []string) error {
		if err := complianceMgr.Suspend(args[0]); err != nil {
			ilog.Error("cli_suspend", "error", err)
			return err
		}
		ilog.Info("cli_suspend", "address", args[0])
		cmOutput(map[string]string{"status": "suspended"}, "suspended")
		return nil
	}}

	resumeCmd := &cobra.Command{Use: "resume [addr]", Args: cobra.ExactArgs(1), Short: "Lift a suspension", RunE: func(cmd *cobra.Command, args []string) error {
		if err := complianceMgr.Resume(args[0]); err != nil {
			ilog.Error("cli_resume", "error", err)
			return err
		}
		ilog.Info("cli_resume", "address", args[0])
		cmOutput(map[string]string{"status": "resumed"}, "resumed")
		return nil
	}}

	whitelistCmd := &cobra.Command{Use: "whitelist [addr]", Args: cobra.ExactArgs(1), Short: "Add an address to the whitelist", RunE: func(cmd *cobra.Command, args []string) error {
		if err := complianceMgr.Whitelist(args[0]); err != nil {
			ilog.Error("cli_whitelist", "error", err)
			return err
		}
		ilog.Info("cli_whitelist", "address", args[0])
		cmOutput(map[string]string{"status": "whitelisted"}, "whitelisted")
		return nil
	}}

	unwhitelistCmd := &cobra.Command{Use: "unwhitelist [addr]", Args: cobra.ExactArgs(1), Short: "Remove an address from the whitelist", RunE: func(cmd *cobra.Command, args []string) error {
		if err := complianceMgr.Unwhitelist(args[0]); err != nil {
			ilog.Error("cli_unwhitelist", "error", err)
			return err
		}
		ilog.Info("cli_unwhitelist", "address", args[0])
		cmOutput(map[string]string{"status": "unwhitelisted"}, "unwhitelisted")
		return nil
	}}

	statusCmd := &cobra.Command{Use: "status [addr]", Args: cobra.ExactArgs(1), Short: "Show suspension and whitelist status", RunE: func(cmd *cobra.Command, args []string) error {
		s, w := complianceMgr.Status(args[0])
		ilog.Info("cli_status", "address", args[0], "suspended", s, "whitelisted", w)
		cmOutput(map[string]bool{"suspended": s, "whitelisted": w}, fmt.Sprintf("suspended: %v whitelisted: %v", s, w))
		return nil
	}}

	reviewCmd := &cobra.Command{Use: "review [tx.json]", Args: cobra.ExactArgs(1), Short: "Check a transaction before broadcast", RunE: func(cmd *cobra.Command, args []string) error {
		var tx core.Transaction
		b, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		if err := json.Unmarshal(b, &tx); err != nil {
			return err
		}
		if err := complianceMgr.ReviewTransaction(tx); err != nil {
			ilog.Error("cli_review", "error", err)
			return err
		}
		ilog.Info("cli_review", "from", tx.From, "to", tx.To)
		cmOutput(map[string]string{"status": "ok"}, "ok")
		return nil
	}}

	mgmtCmd.AddCommand(suspendCmd, resumeCmd, whitelistCmd, unwhitelistCmd, statusCmd, reviewCmd)
	rootCmd.AddCommand(mgmtCmd)
}
