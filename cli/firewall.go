package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var firewall = core.NewFirewall()

func init() {
	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "Manage firewall rules",
	}

	allowCmd := &cobra.Command{
		Use:   "allow <ip>",
		Short: "Allow connections from an IP",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("FirewallAllow")
			firewall.Allow(args[0])
			printOutput("ip allowed")
		},
	}
	cmd.AddCommand(allowCmd)

	blockCmd := &cobra.Command{
		Use:   "block <ip>",
		Short: "Block connections from an IP",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("FirewallBlock")
			firewall.Block(args[0])
			printOutput("ip blocked")
		},
	}
	cmd.AddCommand(blockCmd)

	checkCmd := &cobra.Command{
		Use:   "check <ip>",
		Short: "Check if an IP is allowed",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("FirewallCheck")
			printOutput(firewall.IsAllowed(args[0]))
		},
	}
	cmd.AddCommand(checkCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List firewall rules",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("FirewallList")
			allowed, blocked := firewall.Rules()
			printOutput(map[string][]string{"allowed": allowed, "blocked": blocked})
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
