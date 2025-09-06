package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var proposalMgr = core.NewProposalManager()

func init() {
	propCmd := &cobra.Command{
		Use:   "dao-proposal",
		Short: "Manage DAO proposals",
	}

	var createJSON bool
	var createPub, createMsg, createSig string
	createCmd := &cobra.Command{
		Use:   "create <daoID> <creator> <description>",
		Args:  cobra.MinimumNArgs(3),
		Short: "Create a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("CreateProposal")
			ok, err := VerifySignature(createPub, createMsg, createSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			desc := strings.Join(args[2:], " ")
			p := proposalMgr.CreateProposal(dao, args[1], desc)
			if createJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"id": p.ID})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), p.ID)
		},
	}
	createCmd.Flags().BoolVar(&createJSON, "json", false, "output as JSON")
	createCmd.Flags().StringVar(&createPub, "pub", "", "hex encoded public key")
	createCmd.Flags().StringVar(&createMsg, "msg", "", "hex encoded message")
	createCmd.Flags().StringVar(&createSig, "sig", "", "hex encoded signature")

	var voteJSON bool
	var votePub, voteMsg, voteSig string
	voteCmd := &cobra.Command{
		Use:   "vote <id> <voter> <weight> <yes|no>",
		Args:  cobra.ExactArgs(4),
		Short: "Cast a vote",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("Vote")
			ok, err := VerifySignature(votePub, voteMsg, voteSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			w, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid weight")
				return
			}
			support := strings.ToLower(args[3]) == "yes"
			if err := proposalMgr.Vote(args[0], args[1], w, support); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if voteJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "vote recorded"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "vote recorded")
		},
	}
	voteCmd.Flags().BoolVar(&voteJSON, "json", false, "output as JSON")
	voteCmd.Flags().StringVar(&votePub, "pub", "", "hex encoded public key")
	voteCmd.Flags().StringVar(&voteMsg, "msg", "", "hex encoded message")
	voteCmd.Flags().StringVar(&voteSig, "sig", "", "hex encoded signature")

	var resultsJSON bool
	resultsCmd := &cobra.Command{
		Use:   "results <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show voting results",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ProposalStatus")
			yes, no, err := proposalMgr.Results(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if resultsJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]uint64{"yes": yes, "no": no})
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "yes:%d no:%d\n", yes, no)
		},
	}
	resultsCmd.Flags().BoolVar(&resultsJSON, "json", false, "output as JSON")

	var execJSON bool
	var execPub, execMsg, execSig string
	executeCmd := &cobra.Command{
		Use:   "execute <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Mark proposal executed",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ExecuteProposal")
			ok, err := VerifySignature(execPub, execMsg, execSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			if err := proposalMgr.Execute(args[0]); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if execJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "executed"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "executed")
		},
	}
	executeCmd.Flags().BoolVar(&execJSON, "json", false, "output as JSON")
	executeCmd.Flags().StringVar(&execPub, "pub", "", "hex encoded public key")
	executeCmd.Flags().StringVar(&execMsg, "msg", "", "hex encoded message")
	executeCmd.Flags().StringVar(&execSig, "sig", "", "hex encoded signature")

	var getJSON bool
	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get proposal info",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("GetProposal")
			p, err := proposalMgr.Get(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if getJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(p)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s DAO:%s Desc:%s Executed:%v\n", p.ID, p.DAOID, p.Desc, p.Executed)
		},
	}
	getCmd.Flags().BoolVar(&getJSON, "json", false, "output as JSON")

	var listJSON bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ListProposals")
			if listJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(proposalMgr.List())
				return
			}
			for _, p := range proposalMgr.List() {
				fmt.Fprintf(cmd.OutOrStdout(), "%s DAO:%s Desc:%s Executed:%v\n", p.ID, p.DAOID, p.Desc, p.Executed)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	propCmd.AddCommand(createCmd, voteCmd, resultsCmd, executeCmd, getCmd, listCmd)
	rootCmd.AddCommand(propCmd)
}
