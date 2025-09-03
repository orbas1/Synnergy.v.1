package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	plasmaBridge = core.NewPlasmaBridge()
	plasmaJSON   bool
)

var plasmaCmd = &cobra.Command{
	Use:   "plasma",
	Short: "Interact with the Plasma bridge",
}

func init() {
	plasmaCmd.PersistentFlags().BoolVar(&plasmaJSON, "json", false, "output as JSON")
	rootCmd.AddCommand(plasmaCmd)
}
