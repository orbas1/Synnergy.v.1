package cli

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"synnergy/core"
)

var complianceMgr = core.NewComplianceManager()

func init() {
	mgmtCmd := &cobra.Command{Use: "compliance_management", Short: "Manage address compliance"}

	suspendCmd := &cobra.Command{Use: "suspend [addr]", Args: cobra.ExactArgs(1), Short: "Suspend an address", Run: func(cmd *cobra.Command, args []string) {
		complianceMgr.Suspend(args[0])
	}}

	resumeCmd := &cobra.Command{Use: "resume [addr]", Args: cobra.ExactArgs(1), Short: "Lift a suspension", Run: func(cmd *cobra.Command, args []string) {
		complianceMgr.Resume(args[0])
	}}

	whitelistCmd := &cobra.Command{Use: "whitelist [addr]", Args: cobra.ExactArgs(1), Short: "Add an address to the whitelist", Run: func(cmd *cobra.Command, args []string) {
		complianceMgr.Whitelist(args[0])
	}}

	unwhitelistCmd := &cobra.Command{Use: "unwhitelist [addr]", Args: cobra.ExactArgs(1), Short: "Remove an address from the whitelist", Run: func(cmd *cobra.Command, args []string) {
		complianceMgr.Unwhitelist(args[0])
	}}

	statusCmd := &cobra.Command{Use: "status [addr]", Args: cobra.ExactArgs(1), Short: "Show suspension and whitelist status", Run: func(cmd *cobra.Command, args []string) {
		s, w := complianceMgr.Status(args[0])
		fmt.Printf("suspended: %v whitelisted: %v\n", s, w)
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
			return err
		}
		fmt.Println("ok")
		return nil
	}}

	mgmtCmd.AddCommand(suspendCmd, resumeCmd, whitelistCmd, unwhitelistCmd, statusCmd, reviewCmd)
	rootCmd.AddCommand(mgmtCmd)
}
