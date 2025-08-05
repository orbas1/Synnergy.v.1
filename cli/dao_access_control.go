package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	memberCmd := &cobra.Command{
		Use:   "dao-members",
		Short: "Manage DAO membership roles",
	}

	addCmd := &cobra.Command{
		Use:   "add <daoID> <addr> <role>",
		Args:  cobra.ExactArgs(3),
		Short: "Add member with role",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			if err := dao.AddMember(args[1], args[2]); err != nil {
				fmt.Println(err)
			}
		},
	}

	removeCmd := &cobra.Command{
		Use:   "remove <daoID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Remove a member",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			dao.RemoveMember(args[1])
		},
	}

	roleCmd := &cobra.Command{
		Use:   "role <daoID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Get member role",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			role, ok := dao.MemberRole(args[1])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Println(role)
		},
	}

	listCmd := &cobra.Command{
		Use:   "list <daoID>",
		Args:  cobra.ExactArgs(1),
		Short: "List members",
		Run: func(cmd *cobra.Command, args []string) {
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			for addr, role := range dao.MembersList() {
				fmt.Printf("%s: %s\n", addr, role)
			}
		},
	}

	memberCmd.AddCommand(addCmd, removeCmd, roleCmd, listCmd)
	rootCmd.AddCommand(memberCmd)
}
