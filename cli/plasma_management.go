package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	synnergy "synnergy"
)

func init() {
	mgmtCmd := &cobra.Command{
		Use:   "plasma-mgmt",
		Short: "Manage Plasma bridge operations",
	}

	pauseCmd := &cobra.Command{
		Use:   "pause",
		Short: "Pause Plasma bridge operations",
		Run: func(cmd *cobra.Command, args []string) {
			plasmaBridge.Pause()
			fmt.Printf("gas:%d\n", synnergy.GasCost("PlasmaPause"))
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume Plasma bridge operations",
		Run: func(cmd *cobra.Command, args []string) {
			plasmaBridge.Resume()
			fmt.Printf("gas:%d\n", synnergy.GasCost("PlasmaResume"))
		},
	}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show whether Plasma bridge is paused",
		Run: func(cmd *cobra.Command, args []string) {
			status := plasmaBridge.Status()
			if plasmaJSON {
				enc, _ := json.Marshal(map[string]bool{"paused": status})
				fmt.Println(string(enc))
				return
			}
			fmt.Println(status)
		},
	}

	mgmtCmd.AddCommand(pauseCmd, resumeCmd, statusCmd)
	plasmaCmd.AddCommand(mgmtCmd)
}
