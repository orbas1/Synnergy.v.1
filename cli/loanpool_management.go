package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var loanMgr = core.NewLoanPoolManager(loanPool)

func init() {
	cmd := &cobra.Command{
		Use:   "loanmgr",
		Short: "Administrative loan pool controls",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "pause",
		Short: "Pause new proposals",
		Run:   func(cmd *cobra.Command, args []string) { loanMgr.Pause() },
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "resume",
		Short: "Resume proposal submissions",
		Run:   func(cmd *cobra.Command, args []string) { loanMgr.Resume() },
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "stats",
		Short: "Display treasury and proposal stats",
		Run: func(cmd *cobra.Command, args []string) {
			stats := loanMgr.Stats()
			b, _ := json.MarshalIndent(stats, "", "  ")
			fmt.Println(string(b))
		},
	})

	rootCmd.AddCommand(cmd)
}
