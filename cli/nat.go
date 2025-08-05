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
		Run: func(cmd *cobra.Command, args []string) {
			p, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("invalid port")
				return
			}
			natMgr.Map(p)
		},
	}

	unmapCmd := &cobra.Command{
		Use:   "unmap",
		Short: "Remove current port mapping",
		Run:   func(cmd *cobra.Command, args []string) { natMgr.Unmap() },
	}

	ipCmd := &cobra.Command{
		Use:   "ip",
		Short: "Show external IP",
		Run:   func(cmd *cobra.Command, args []string) { fmt.Println(natMgr.ExternalIP()) },
	}

	natCmd.AddCommand(mapCmd, unmapCmd, ipCmd)
	rootCmd.AddCommand(natCmd)
}
