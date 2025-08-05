package cli

import (
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

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an elected authority node",
		Run: func(cmd *cobra.Command, args []string) {
			addr, _ := cmd.Flags().GetString("addr")
			role, _ := cmd.Flags().GetString("role")
			termDur, _ := cmd.Flags().GetDuration("term")
			electedNode = core.NewElectedAuthorityNode(addr, role, termDur)
			fmt.Println("elected authority node created")
		},
	}
	createCmd.Flags().String("addr", "", "node address")
	createCmd.Flags().String("role", "validator", "node role")
	createCmd.Flags().Duration("term", time.Hour, "term duration")
	cmd.AddCommand(createCmd)

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show node status",
		Run: func(cmd *cobra.Command, args []string) {
			if electedNode == nil {
				fmt.Println("node not initialised")
				return
			}
			fmt.Printf("address: %s role: %s term_end: %s active: %v\n",
				electedNode.Address,
				electedNode.Role,
				electedNode.TermEnd.Format(time.RFC3339),
				electedNode.IsActive(time.Now()))
		},
	}
	cmd.AddCommand(statusCmd)

	rootCmd.AddCommand(cmd)
}
