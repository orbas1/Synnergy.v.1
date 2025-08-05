package cli

import (
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

	createCmd := &cobra.Command{
		Use:   "create <daoID> <creator> <description>",
		Args:  cobra.MinimumNArgs(3),
		Short: "Create a proposal",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			desc := strings.Join(args[2:], " ")
			p := proposalMgr.CreateProposal(dao, args[1], desc)
			fmt.Println(p.ID)
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote <id> <voter> <weight> <yes|no>",
		Args:  cobra.ExactArgs(4),
		Short: "Cast a vote",
		Run: func(cmd *cobra.Command, args []string) {
			w, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid weight")
				return
			}
			support := strings.ToLower(args[3]) == "yes"
			if err := proposalMgr.Vote(args[0], args[1], w, support); err != nil {
				fmt.Println(err)
			}
		},
	}

	resultsCmd := &cobra.Command{
		Use:   "results <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show voting results",
		Run: func(cmd *cobra.Command, args []string) {
			yes, no, err := proposalMgr.Results(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("yes:%d no:%d\n", yes, no)
		},
	}

	executeCmd := &cobra.Command{
		Use:   "execute <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Mark proposal executed",
		Run: func(cmd *cobra.Command, args []string) {
			if err := proposalMgr.Execute(args[0]); err != nil {
				fmt.Println(err)
			}
		},
	}

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Get proposal info",
		Run: func(cmd *cobra.Command, args []string) {
			p, err := proposalMgr.Get(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%s DAO:%s Desc:%s Executed:%v\n", p.ID, p.DAOID, p.Desc, p.Executed)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		Run: func(cmd *cobra.Command, args []string) {
			for _, p := range proposalMgr.List() {
				fmt.Printf("%s DAO:%s Desc:%s Executed:%v\n", p.ID, p.DAOID, p.Desc, p.Executed)
			}
		},
	}

	propCmd.AddCommand(createCmd, voteCmd, resultsCmd, executeCmd, getCmd, listCmd)
	rootCmd.AddCommand(propCmd)
}
