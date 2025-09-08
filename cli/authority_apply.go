package cli

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			approve, _ := strconv.ParseBool(args[2])
			sigHex, _ := cmd.Flags().GetString("sig")
			pubHex, _ := cmd.Flags().GetString("pub")
			sig, err := hex.DecodeString(sigHex)
			if err != nil {
				return err
			}
			pubBytes, err := hex.DecodeString(pubHex)
			if err != nil {
				return err
			}
			return applyManager.Vote(args[0], args[1], approve, sig, ed25519.PublicKey(pubBytes))
		},
	}
	voteCmd.Flags().String("sig", "", "hex-encoded signature of 'id:approve'")
	voteCmd.Flags().String("pub", "", "hex-encoded public key of voter")

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

	var getJSON bool
	var listJSON bool

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
			if getJSON {
				enc, _ := json.Marshal(app)
				fmt.Println(string(enc))
				return
			}
			fmt.Printf("%s %s approvals:%d rejections:%d finalized:%v\n", app.ID, app.Candidate, len(app.Approvals), len(app.Rejections), app.Finalized)
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List applications",
		Run: func(cmd *cobra.Command, args []string) {
			apps := applyManager.List()
			if listJSON {
				enc, _ := json.Marshal(apps)
				fmt.Println(string(enc))
				return
			}
			for _, app := range apps {
				fmt.Printf("%s %s finalized:%v\n", app.ID, app.Candidate, app.Finalized)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

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
