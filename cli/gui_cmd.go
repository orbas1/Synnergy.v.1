package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// guiCmd launches the desktop GUI shell when available.
var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch the desktop GUI shell",
	RunE: func(cmd *cobra.Command, args []string) error {
		shellDir := filepath.Join("GUI", "desktop-shell")
		if _, err := os.Stat(shellDir); err != nil {
			if os.IsNotExist(err) {
				abs, _ := filepath.Abs(shellDir)
				return fmt.Errorf("desktop GUI shell not found at %s", abs)
			}
			return fmt.Errorf("unable to access GUI directory: %w", err)
		}
		c := exec.Command("npm", "start")
		c.Dir = shellDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Env = os.Environ() // preserve env for auth
		return c.Run()
	},
}
