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
		Use:   "add",
		Short: "Add a health record",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetUint64("id")
			owner, _ := cmd.Flags().GetString("owner")
			data, _ := cmd.Flags().GetString("data")
			if id == 0 || owner == "" || data == "" {
				return fmt.Errorf("id, owner and data must be provided")
			}
			if err := syn1100.AddRecord(tokens.TokenID(id), owner, []byte(data)); err != nil {
				return err
			}
			cmd.Println("record added")
			return nil
		},
	}
	addCmd.Flags().Uint64("id", 0, "record ID")
	addCmd.Flags().String("owner", "", "record owner")
	addCmd.Flags().String("data", "", "record data")
	addCmd.MarkFlagRequired("id")
	addCmd.MarkFlagRequired("owner")
	addCmd.MarkFlagRequired("data")
	cmd.AddCommand(addCmd)

	grantCmd := &cobra.Command{
		Use:   "grant",
		Short: "Grant access",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetUint64("id")
			grantee, _ := cmd.Flags().GetString("grantee")
			if id == 0 || grantee == "" {
				return fmt.Errorf("id and grantee must be provided")
			}
			if err := syn1100.GrantAccess(tokens.TokenID(id), grantee); err != nil {
				return err
			}
			cmd.Println("access granted")
			return nil
		},
	}
	grantCmd.Flags().Uint64("id", 0, "record ID")
	grantCmd.Flags().String("grantee", "", "grantee address")
	grantCmd.MarkFlagRequired("id")
	grantCmd.MarkFlagRequired("grantee")
	cmd.AddCommand(grantCmd)

	revokeCmd := &cobra.Command{
		Use:   "revoke",
		Short: "Revoke access",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetUint64("id")
			grantee, _ := cmd.Flags().GetString("grantee")
			if id == 0 || grantee == "" {
				return fmt.Errorf("id and grantee must be provided")
			}
			if err := syn1100.RevokeAccess(tokens.TokenID(id), grantee); err != nil {
				return err
			}
			cmd.Println("access revoked")
			return nil
		},
	}
	revokeCmd.Flags().Uint64("id", 0, "record ID")
	revokeCmd.Flags().String("grantee", "", "grantee address")
	revokeCmd.MarkFlagRequired("id")
	revokeCmd.MarkFlagRequired("grantee")
	cmd.AddCommand(revokeCmd)

	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Retrieve record",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetUint64("id")
			caller, _ := cmd.Flags().GetString("caller")
			if id == 0 || caller == "" {
				return fmt.Errorf("id and caller must be provided")
			}
			data, err := syn1100.GetRecord(tokens.TokenID(id), caller)
			if err != nil {
				return err
			}
			cmd.Println(string(data))
			return nil
		},
	}
	getCmd.Flags().Uint64("id", 0, "record ID")
	getCmd.Flags().String("caller", "", "caller address")
	getCmd.MarkFlagRequired("id")
	getCmd.MarkFlagRequired("caller")
	cmd.AddCommand(getCmd)

	rootCmd.AddCommand(cmd)
}
