package cli

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanmgrPause")
			loanMgr.Pause()
			printOutput("paused")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "resume",
		Short: "Resume proposal submissions",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanmgrResume")
			loanMgr.Resume()
			printOutput("resumed")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "stats",
		Short: "Display treasury and proposal stats",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LoanmgrStats")
			stats := loanMgr.Stats()
			printOutput(stats)
		},
	})

	rootCmd.AddCommand(cmd)
}
