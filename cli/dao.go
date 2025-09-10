package cli

import (
	"github.com/spf13/cobra"
	"synnergy/core"
)

var daoMgr = core.NewDAOManager()

func init() {
	daoCmd := &cobra.Command{
		Use:   "dao",
		Short: "Manage decentralised autonomous organisations",
	}

	createCmd := &cobra.Command{
		Use:   "create <name> <creator>",
		Args:  cobra.ExactArgs(2),
		Short: "Create a new DAO",
		Run: func(cmd *cobra.Command, args []string) {
			daoMgr.AuthorizeRelayer(args[1])
			dao, err := daoMgr.Create(args[0], args[1])
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("CreateDAO")
			printOutput(map[string]any{"id": dao.ID})
		},
	}

	joinCmd := &cobra.Command{
		Use:   "join <id> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Join a DAO",
		Run: func(cmd *cobra.Command, args []string) {
			daoMgr.AuthorizeRelayer(args[1])
			if err := daoMgr.Join(args[0], args[1]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("JoinDAO")
			printOutput(map[string]any{"status": "joined", "id": args[0], "address": args[1]})
		},
	}

	leaveCmd := &cobra.Command{
		Use:   "leave <id> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Leave a DAO",
		Run: func(cmd *cobra.Command, args []string) {
			daoMgr.AuthorizeRelayer(args[1])
			if err := daoMgr.Leave(args[0], args[1]); err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("LeaveDAO")
			printOutput(map[string]any{"status": "left", "id": args[0], "address": args[1]})
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show DAO information",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				printOutput(map[string]any{"error": err.Error()})
				return
			}
			gasPrint("DAOInfo")
			printOutput(map[string]any{
				"id":      dao.ID,
				"name":    dao.Name,
				"creator": dao.Creator,
				"members": len(dao.Members),
			})
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all DAOs",
		Run: func(cmd *cobra.Command, args []string) {
			daos := daoMgr.List()
			gasPrint("ListDAOs")
			printOutput(daos)
		},
	}

	daoCmd.AddCommand(createCmd, joinCmd, leaveCmd, infoCmd, listCmd)
	rootCmd.AddCommand(daoCmd)
}
