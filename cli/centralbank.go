package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var centralBank = core.NewCentralBankingNode("central", "cbnode", ledger, "neutral")

func init() {
	cbCmd := &cobra.Command{Use: "centralbank", Short: "Central bank node operations"}

	infoCmd := &cobra.Command{Use: "info", Short: "Show node info", Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("id: %s address: %s policy: %s\n", centralBank.ID, centralBank.Addr, centralBank.MonetaryPolicy)
	}}

	policyCmd := &cobra.Command{Use: "policy [description]", Args: cobra.ExactArgs(1), Short: "Update monetary policy", Run: func(cmd *cobra.Command, args []string) {
		centralBank.UpdatePolicy(args[0])
	}}

	mintCmd := &cobra.Command{Use: "mint [to] [amount]", Args: cobra.ExactArgs(2), Short: "Mint currency tokens", Run: func(cmd *cobra.Command, args []string) {
		amt, _ := strconv.ParseUint(args[1], 10, 64)
		centralBank.Mint(args[0], amt)
	}}

	cbCmd.AddCommand(infoCmd, policyCmd, mintCmd)
	rootCmd.AddCommand(cbCmd)
}
