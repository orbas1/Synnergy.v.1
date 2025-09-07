package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn2500 = core.NewSyn2500Registry()

func parseMeta(meta string) map[string]string {
	m := make(map[string]string)
	if meta == "" {
		return m
	}
	for _, kv := range strings.Split(meta, ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}

func init() {
	cmd := &cobra.Command{
		Use:   "syn2500",
		Short: "DAO membership registry",
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a member",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			power, _ := cmd.Flags().GetUint64("power")
			if id == "" || addr == "" || power == 0 {
				return errors.New("id, addr and power are required")
			}
			meta, _ := cmd.Flags().GetString("meta")
			m := core.NewSyn2500Member(id, addr, power, parseMeta(meta))
			syn2500.AddMember(m)
			fmt.Println("member added")
			return nil
		},
	}
	addCmd.Flags().String("id", "", "member id")
	addCmd.Flags().String("addr", "", "address")
	addCmd.Flags().Uint64("power", 0, "voting power")
	addCmd.Flags().String("meta", "", "metadata key=value,comma-separated")
	cmd.AddCommand(addCmd)

	getCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get member info",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, ok := syn2500.GetMember(args[0])
			if !ok {
				return errors.New("not found")
			}
			fmt.Printf("%+v\n", *m)
			return nil
		},
	}
	cmd.AddCommand(getCmd)

	delCmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove member",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, ok := syn2500.GetMember(args[0]); !ok {
				return errors.New("not found")
			}
			syn2500.RemoveMember(args[0])
			fmt.Println("member removed")
			return nil
		},
	}
	cmd.AddCommand(delCmd)

	updateCmd := &cobra.Command{
		Use:   "update <id> <power>",
		Short: "Update voting power for a member",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			m, ok := syn2500.GetMember(args[0])
			if !ok {
				return errors.New("not found")
			}
			var power uint64
			if _, err := fmt.Sscanf(args[1], "%d", &power); err != nil {
				return err
			}
			m.UpdateVotingPower(power)
			fmt.Println("voting power updated")
			return nil
		},
	}
	cmd.AddCommand(updateCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List members",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, m := range syn2500.ListMembers() {
				fmt.Printf("%s %s %d\n", m.ID, m.Address, m.VotingPower)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
