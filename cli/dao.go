package cli

import (
	"encoding/json"
	"fmt"

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
			gasPrint("CreateDAO")
			dao := daoMgr.Create(args[0], args[1])
			if dao == nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid parameters")
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), dao.ID)
		},
	}

	joinCmd := &cobra.Command{
		Use:   "join <id> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Join a DAO",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("JoinDAO")
			if err := daoMgr.Join(args[0], args[1]); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
			}
		},
	}

	leaveCmd := &cobra.Command{
		Use:   "leave <id> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Leave a DAO",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("LeaveDAO")
			if err := daoMgr.Leave(args[0], args[1]); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
			}
		},
	}

	var infoJSON bool
	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show DAO information",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("DAOInfo")
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			if infoJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(dao)
				return
			}
			fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", dao.ID, dao.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "creator: %s\n", dao.Creator)
			fmt.Fprintf(cmd.OutOrStdout(), "members: %d\n", len(dao.Members))
		},
	}
	infoCmd.Flags().BoolVar(&infoJSON, "json", false, "output as JSON")

	var listJSON bool
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all DAOs",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("ListDAOs")
			daos := daoMgr.List()
			if listJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(daos)
				return
			}
			for _, d := range daos {
				fmt.Fprintf(cmd.OutOrStdout(), "%s %s\n", d.ID, d.Name)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	daoCmd.AddCommand(createCmd, joinCmd, leaveCmd, infoCmd, listCmd)
	rootCmd.AddCommand(daoCmd)
}
