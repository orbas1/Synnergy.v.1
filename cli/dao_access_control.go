package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	memberCmd := &cobra.Command{
		Use:   "dao-members",
		Short: "Manage DAO membership roles",
	}

	var addJSON bool
	var addPub, addMsg, addSig string
	addCmd := &cobra.Command{
		Use:   "add <daoID> <addr> <role>",
		Args:  cobra.ExactArgs(3),
		Short: "Add member with role",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("AddMember")
			ok, err := VerifySignature(addPub, addMsg, addSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if err := dao.AddMember(args[1], args[2]); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if addJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "member added"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "member added")
		},
	}
	addCmd.Flags().BoolVar(&addJSON, "json", false, "output as JSON")
	addCmd.Flags().StringVar(&addPub, "pub", "", "hex encoded public key")
	addCmd.Flags().StringVar(&addMsg, "msg", "", "hex encoded message")
	addCmd.Flags().StringVar(&addSig, "sig", "", "hex encoded signature")

	var removeJSON bool
	var removePub, removeMsg, removeSig string
	removeCmd := &cobra.Command{
		Use:   "remove <daoID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Remove a member",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("RemoveMember")
			ok, err := VerifySignature(removePub, removeMsg, removeSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if err := dao.RemoveMember(args[1]); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if removeJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "member removed"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "member removed")
		},
	}
	removeCmd.Flags().BoolVar(&removeJSON, "json", false, "output as JSON")
	removeCmd.Flags().StringVar(&removePub, "pub", "", "hex encoded public key")
	removeCmd.Flags().StringVar(&removeMsg, "msg", "", "hex encoded message")
	removeCmd.Flags().StringVar(&removeSig, "sig", "", "hex encoded signature")

	var updateJSON bool
	var updatePub, updateMsg, updateSig string
	updateCmd := &cobra.Command{
		Use:   "update <daoID> <admin> <addr> <role>",
		Args:  cobra.ExactArgs(4),
		Short: "Update member role",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("UpdateMemberRole")
			ok, err := VerifySignature(updatePub, updateMsg, updateSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if err := dao.UpdateMemberRole(args[1], args[2], args[3]); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if updateJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "role updated"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "role updated")
		},
	}
	updateCmd.Flags().BoolVar(&updateJSON, "json", false, "output as JSON")
	updateCmd.Flags().StringVar(&updatePub, "pub", "", "hex encoded public key")
	updateCmd.Flags().StringVar(&updateMsg, "msg", "", "hex encoded message")
	updateCmd.Flags().StringVar(&updateSig, "sig", "", "hex encoded signature")

	var roleJSON bool
	roleCmd := &cobra.Command{
		Use:   "role <daoID> <addr>",
		Args:  cobra.ExactArgs(2),
		Short: "Get member role",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("MemberRole")
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			role, ok := dao.MemberRole(args[1])
			if !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "not found")
				return
			}
			if roleJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"role": role})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), role)
		},
	}
	roleCmd.Flags().BoolVar(&roleJSON, "json", false, "output as JSON")

	var listJSON bool
	listCmd := &cobra.Command{
		Use:   "list <daoID>",
		Args:  cobra.ExactArgs(1),
		Short: "List members",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("MembersList")
			dao, err := daoMgr.Info(args[0])
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if listJSON {
				members := make([]map[string]string, 0, len(dao.MembersList()))
				for addr, role := range dao.MembersList() {
					members = append(members, map[string]string{"addr": addr, "role": role})
				}
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(members)
				return
			}
			for addr, role := range dao.MembersList() {
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", addr, role)
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

	memberCmd.AddCommand(addCmd, updateCmd, removeCmd, roleCmd, listCmd)
	rootCmd.AddCommand(memberCmd)
}
