package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
	"synnergy/internal/tokens"
)

var centralBankToken = tokens.NewSYN10Token(1, "CBDC", "cSYN", "central", 1, 2)
var centralBank = core.NewCentralBankingNode("central", "cbnode", ledger, "neutral", centralBankToken)
var centralBankJSON bool

func init() {
	cbCmd := &cobra.Command{Use: "centralbank", Short: "Central bank node operations"}
	cbCmd.PersistentFlags().BoolVar(&centralBankJSON, "json", false, "output as JSON")

	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "Show node info",
		RunE: func(cmd *cobra.Command, args []string) error {
			if centralBankJSON {
				enc := json.NewEncoder(os.Stdout)
				return enc.Encode(map[string]string{"id": centralBank.ID, "address": centralBank.Addr, "policy": centralBank.MonetaryPolicy})
			}
			fmt.Fprintf(os.Stdout, "id: %s address: %s policy: %s\n", centralBank.ID, centralBank.Addr, centralBank.MonetaryPolicy)
			return nil
		},
	}

	policyCmd := &cobra.Command{
		Use:   "policy [description]",
		Args:  cobra.ExactArgs(1),
		Short: "Update monetary policy",
		RunE: func(cmd *cobra.Command, args []string) error {
			centralBank.UpdatePolicy(args[0])
			return nil
		},
	}

	mintCmd := &cobra.Command{Use: "mint [to] [amount]", Args: cobra.ExactArgs(2), Short: "Mint CBDC tokens", RunE: func(cmd *cobra.Command, args []string) error {
		amt, _ := strconv.ParseUint(args[1], 10, 64)
		return centralBank.MintCBDC(args[0], amt)
	}}

	cbCmd.AddCommand(infoCmd, policyCmd, mintCmd)
	rootCmd.AddCommand(cbCmd)
}
