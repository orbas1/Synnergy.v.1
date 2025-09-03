package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

// difficultyMgr adjusts PoW difficulty from block time samples.
var difficultyMgr = core.NewDifficultyManager(consensus, 10, 1, 10)

func init() {
	cmd := &cobra.Command{
		Use:   "consensus-difficulty",
		Short: "Manage PoW difficulty",
	}

	sampleCmd := &cobra.Command{
		Use:   "sample [seconds]",
		Args:  cobra.ExactArgs(1),
		Short: "Add block time sample and show new difficulty",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("AddSample")
			d, _ := strconv.ParseFloat(args[0], 64)
			fmt.Println(difficultyMgr.AddSample(d))
		},
	}

	valueCmd := &cobra.Command{
		Use:   "value",
		Short: "Show current difficulty",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Difficulty")
			fmt.Println(difficultyMgr.Difficulty())
		},
	}

	configCmd := &cobra.Command{
		Use:   "config [window] [initial] [target]",
		Args:  cobra.ExactArgs(3),
		Short: "Reconfigure difficulty manager",
		Run: func(cmd *cobra.Command, args []string) {
			w, _ := strconv.Atoi(args[0])
			initDiff, _ := strconv.ParseFloat(args[1], 64)
			target, _ := strconv.ParseFloat(args[2], 64)
			difficultyMgr = core.NewDifficultyManager(consensus, w, initDiff, target)
		},
	}

	cmd.AddCommand(sampleCmd, valueCmd, configCmd)
	rootCmd.AddCommand(cmd)
}
