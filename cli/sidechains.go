package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	sidechainRegistry = core.NewSidechainRegistry()
	sidechainOps      = core.NewSidechainOps(sidechainRegistry)
)

func init() {
	scCmd := &cobra.Command{
		Use:   "sidechain",
		Short: "Manage side-chains",
	}

	registerCmd := &cobra.Command{
		Use:   "register <id> <meta> [validators...]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Register a new side-chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			sc, err := sidechainRegistry.Register(args[0], args[1], args[2:])
			if err != nil {
				return err
			}
			fmt.Println(sc.ID)
			return nil
		},
	}

	headerCmd := &cobra.Command{
		Use:   "header <id> <header>",
		Args:  cobra.ExactArgs(2),
		Short: "Submit a side-chain header",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sidechainRegistry.SubmitHeader(args[0], args[1])
		},
	}

	getHeaderCmd := &cobra.Command{
		Use:   "get-header <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get the latest header for a side-chain",
		Run: func(cmd *cobra.Command, args []string) {
			h, ok := sidechainRegistry.GetHeader(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Println(h)
		},
	}

	metaCmd := &cobra.Command{
		Use:   "meta <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Display side-chain metadata",
		Run: func(cmd *cobra.Command, args []string) {
			sc, ok := sidechainRegistry.Meta(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%s %s %v paused=%v\n", sc.ID, sc.Metadata, sc.Validators, sc.Paused)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered side-chains",
		Run: func(cmd *cobra.Command, args []string) {
			for _, sc := range sidechainRegistry.List() {
				fmt.Printf("%s %s paused=%v\n", sc.ID, sc.Header, sc.Paused)
			}
		},
	}

	pauseCmd := &cobra.Command{
		Use:   "pause <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Pause a side-chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sidechainRegistry.Pause(args[0])
		},
	}

	resumeCmd := &cobra.Command{
		Use:   "resume <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Resume a paused side-chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sidechainRegistry.Resume(args[0])
		},
	}

	updateValidatorsCmd := &cobra.Command{
		Use:   "update-validators <id> <v1> [v2...]",
		Args:  cobra.MinimumNArgs(2),
		Short: "Update validator set for a side-chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sidechainRegistry.UpdateValidators(args[0], args[1:])
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Remove a side-chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return sidechainRegistry.Remove(args[0])
		},
	}

	scCmd.AddCommand(registerCmd, headerCmd, getHeaderCmd, metaCmd, listCmd, pauseCmd, resumeCmd, updateValidatorsCmd, removeCmd)
	rootCmd.AddCommand(scCmd)
}
