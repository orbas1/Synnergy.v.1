package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var switcher = core.NewConsensusSwitcher(core.ModePoW)

func parseMode(m string) (core.ConsensusMode, error) {
	switch strings.ToLower(m) {
	case "pow":
		return core.ModePoW, nil
	case "pos":
		return core.ModePoS, nil
	case "poh":
		return core.ModePoH, nil
	default:
		return core.ModePoW, fmt.Errorf("unknown mode %s", m)
	}
}

func init() {
	cmd := &cobra.Command{
		Use:   "consensus-mode",
		Short: "Evaluate and view consensus mode",
	}

	evaluateCmd := &cobra.Command{
		Use:   "evaluate",
		Short: "Evaluate mode based on current weights",
		Run: func(cmd *cobra.Command, args []string) {
			mode := switcher.Evaluate(consensus)
			fmt.Println(mode)
		},
	}

	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show last evaluated mode",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(switcher.Mode())
		},
	}

	setCmd := &cobra.Command{
		Use:   "set [mode]",
		Args:  cobra.ExactArgs(1),
		Short: "Set initial mode (pow|pos|poh)",
		RunE: func(cmd *cobra.Command, args []string) error {
			mode, err := parseMode(args[0])
			if err != nil {
				return err
			}
			switcher = core.NewConsensusSwitcher(mode)
			return nil
		},
	}

	cmd.AddCommand(evaluateCmd, showCmd, setCmd)
	rootCmd.AddCommand(cmd)
}
