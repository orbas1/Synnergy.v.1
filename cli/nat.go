package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var natMgr = core.NewNATManager()

func init() {
	natCmd := &cobra.Command{
		Use:   "nat",
		Short: "Manage NAT port mappings",
	}

	mapCmd := &cobra.Command{
		Use:   "map [port]",
		Args:  cobra.ExactArgs(1),
		Short: "Map a local port",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := strconv.Atoi(args[0])
			if err != nil || p <= 0 || p > 65535 {
				return fmt.Errorf("invalid port: %s", args[0])
			}
			id, _ := cmd.Flags().GetString("id")
			natMgr.MapPort(id, p)
			fmt.Fprintf(cmd.OutOrStdout(), "mapped %d\n", p)
			return nil
		},
	}
	mapCmd.Flags().String("id", "self", "mapping identifier")

	unmapCmd := &cobra.Command{
		Use:   "unmap",
		Short: "Remove current port mapping",
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			natMgr.RemoveMapping(id)
			fmt.Fprintln(cmd.OutOrStdout(), "unmapped")
		},
	}
	unmapCmd.Flags().String("id", "self", "mapping identifier")

	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "Show external IP",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), natMgr.ExternalIP())
		},
	}

	natCmd.AddCommand(mapCmd, unmapCmd, ipCmd)
	rootCmd.AddCommand(natCmd)
}
