package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var plasmaBridge = core.NewPlasmaBridge()

var plasmaCmd = &cobra.Command{
	Use:   "plasma",
	Short: "Interact with the Plasma bridge",
}

func init() {
	rootCmd.AddCommand(plasmaCmd)
}
