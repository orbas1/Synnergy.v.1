package cli

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var electedNode *core.ElectedAuthorityNode

func init() {
	cmd := &cobra.Command{
		Use:   "elected-node",
		Short: "Manage elected authority nodes",
	}

	var createJSON bool
	var createPub, createMsg, createSig string
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an elected authority node",
		Run: func(cmd *cobra.Command, args []string) {
			addr, _ := cmd.Flags().GetString("addr")
			role, _ := cmd.Flags().GetString("role")
			termDur, _ := cmd.Flags().GetDuration("term")
			gasPrint("CreateElectedNode")
			ok, err := VerifySignature(createPub, createMsg, createSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			electedNode = core.NewElectedAuthorityNode(addr, role, termDur)
			if createJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "created"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "elected authority node created")
		},
	}
	createCmd.Flags().String("addr", "", "node address")
	createCmd.Flags().String("role", "validator", "node role")
	createCmd.Flags().Duration("term", time.Hour, "term duration")
	createCmd.Flags().BoolVar(&createJSON, "json", false, "output as JSON")
	createCmd.Flags().StringVar(&createPub, "pub", "", "hex encoded public key")
	createCmd.Flags().StringVar(&createMsg, "msg", "", "hex encoded message")
	createCmd.Flags().StringVar(&createSig, "sig", "", "hex encoded signature")
	cmd.AddCommand(createCmd)

	var statusJSON bool
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show node status",
		Run: func(cmd *cobra.Command, args []string) {
			if electedNode == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "node not initialised")
				return
			}
			active := electedNode.IsActive(time.Now())
			if statusJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]any{
					"address":  electedNode.Address,
					"role":     electedNode.Role,
					"term_end": electedNode.TermEnd.Format(time.RFC3339),
					"active":   active,
				})
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "address: %s role: %s term_end: %s active: %v\n",
				electedNode.Address,
				electedNode.Role,
				electedNode.TermEnd.Format(time.RFC3339),
				active)
		},
	}
	statusCmd.Flags().BoolVar(&statusJSON, "json", false, "output as JSON")
	cmd.AddCommand(statusCmd)

	rootCmd.AddCommand(cmd)
}
