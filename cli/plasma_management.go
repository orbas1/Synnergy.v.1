package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	mgmtCmd := &cobra.Command{
		Use:   "plasma-mgmt",
		Short: "Manage Plasma bridge operations",
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause Plasma bridge operations",
		Run:   func(cmd *cobra.Command, args []string) { plasmaBridge.Pause() },
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume Plasma bridge operations",
		Run:   func(cmd *cobra.Command, args []string) { plasmaBridge.Resume() },
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show whether Plasma bridge is paused",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(plasmaBridge.Status())
		},
	}

	mgmtCmd.AddCommand(pauseCmd, resumeCmd, statusCmd)
	plasmaCmd.AddCommand(mgmtCmd)
}
