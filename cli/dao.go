package cli

import (
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
			dao := daoMgr.Create(args[0], args[1])
			fmt.Println(dao.ID)
		},
	}

	joinCmd := &cobra.Command{
		Use:   "join <id> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Join a DAO",
		Run: func(cmd *cobra.Command, args []string) {
			if err := daoMgr.Join(args[0], args[1]); err != nil {
				fmt.Println(err)
			}
		},
	}

	leaveCmd := &cobra.Command{
		Use:   "leave <id> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Leave a DAO",
		Run: func(cmd *cobra.Command, args []string) {
			if err := daoMgr.Leave(args[0], args[1]); err != nil {
				fmt.Println(err)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info <id>",
		Args:  cobra.ExactArgs(1),
		Short: "Show DAO information",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("%s %s\n", dao.ID, dao.Name)
			fmt.Printf("creator: %s\n", dao.Creator)
			fmt.Printf("members: %d\n", len(dao.Members))
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all DAOs",
		Run: func(cmd *cobra.Command, args []string) {
			for _, d := range daoMgr.List() {
				fmt.Printf("%s %s\n", d.ID, d.Name)
			}
		},
	}

	daoCmd.AddCommand(createCmd, joinCmd, leaveCmd, infoCmd, listCmd)
	rootCmd.AddCommand(daoCmd)
}
