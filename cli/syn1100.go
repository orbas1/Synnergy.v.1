package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"synnergy/internal/tokens"
)

var syn1100 = tokens.NewSYN1100Token()

func init() {
	cmd := &cobra.Command{
		Use:   "syn1100",
		Short: "Healthcare record token",
	}

	addCmd := &cobra.Command{
		Use:   "add <id> <owner> <data>",
		Short: "Add a health record",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			if err := syn1100.AddRecord(tokens.TokenID(id), args[1], []byte(args[2])); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("record added")
		},
	}
	cmd.AddCommand(addCmd)

	grantCmd := &cobra.Command{
		Use:   "grant <id> <grantee>",
		Short: "Grant access",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			if err := syn1100.GrantAccess(tokens.TokenID(id), args[1]); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("access granted")
		},
	}
	cmd.AddCommand(grantCmd)

	revokeCmd := &cobra.Command{
		Use:   "revoke <id> <grantee>",
		Short: "Revoke access",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			if err := syn1100.RevokeAccess(tokens.TokenID(id), args[1]); err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("access revoked")
		},
	}
	cmd.AddCommand(revokeCmd)

	getCmd := &cobra.Command{
		Use:   "get <id> <caller>",
		Short: "Retrieve record",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			data, err := syn1100.GetRecord(tokens.TokenID(id), args[1])
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println(string(data))
		},
	}
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
