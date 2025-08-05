package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	authorityRegistry = core.NewAuthorityNodeRegistry()
	applyManager      = core.NewAuthorityApplicationManager(authorityRegistry, time.Hour)
)

func init() {
	applyCmd := &cobra.Command{Use: "authority_apply", Short: "Authority node applications"}

	submitCmd := &cobra.Command{
		Use:   "submit [candidate] [role] [desc]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit an authority node application",
		Run: func(cmd *cobra.Command, args []string) {
			id := applyManager.Submit(args[0], args[1], args[2])
			fmt.Println("application:", id)
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [voter] [id] [approve]",
		Args:  cobra.ExactArgs(3),
		Short: "Vote on an application",
		Run: func(cmd *cobra.Command, args []string) {
			approve, _ := strconv.ParseBool(args[2])
			if err := applyManager.Vote(args[0], args[1], approve); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	finalizeCmd := &cobra.Command{
		Use:   "finalize [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Finalize an application",
		Run: func(cmd *cobra.Command, args []string) {
			if err := applyManager.Finalize(args[0]); err != nil {
				fmt.Println("error:", err)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get [id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get application details",
		Run: func(cmd *cobra.Command, args []string) {
			app, err := applyManager.Get(args[0])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Printf("%s %s approvals:%d rejections:%d finalized:%v\n", app.ID, app.Candidate, len(app.Approvals), len(app.Rejections), app.Finalized)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List applications",
		Run: func(cmd *cobra.Command, args []string) {
			for _, app := range applyManager.List() {
				fmt.Printf("%s %s finalized:%v\n", app.ID, app.Candidate, app.Finalized)
			}
		},
	}

	tickCmd := &cobra.Command{
		Use:   "tick",
		Short: "Remove expired applications",
		Run: func(cmd *cobra.Command, args []string) {
			applyManager.Tick(time.Now())
		},
	}

	applyCmd.AddCommand(submitCmd, voteCmd, finalizeCmd, getCmd, listCmd, tickCmd)
	rootCmd.AddCommand(applyCmd)
}
