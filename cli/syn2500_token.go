package cli

import (
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
		Run: func(cmd *cobra.Command, args []string) {
			id, _ := cmd.Flags().GetString("id")
			addr, _ := cmd.Flags().GetString("addr")
			power, _ := cmd.Flags().GetUint64("power")
			meta, _ := cmd.Flags().GetString("meta")
			m := core.NewSyn2500Member(id, addr, power, parseMeta(meta))
			syn2500.AddMember(m)
			fmt.Println("member added")
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
		Run: func(cmd *cobra.Command, args []string) {
			m, ok := syn2500.GetMember(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			fmt.Printf("%+v\n", *m)
		},
	}
	cmd.AddCommand(getCmd)

	delCmd := &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove member",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			syn2500.RemoveMember(args[0])
			fmt.Println("member removed")
		},
	}
	cmd.AddCommand(delCmd)

	updateCmd := &cobra.Command{
		Use:   "update <id> <power>",
		Short: "Update voting power for a member",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			m, ok := syn2500.GetMember(args[0])
			if !ok {
				fmt.Println("not found")
				return
			}
			var power uint64
			fmt.Sscanf(args[1], "%d", &power)
			m.UpdateVotingPower(power)
			fmt.Println("voting power updated")
		},
	}
	cmd.AddCommand(updateCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List members",
		Run: func(cmd *cobra.Command, args []string) {
			for _, m := range syn2500.ListMembers() {
				fmt.Printf("%s %s %d\n", m.ID, m.Address, m.VotingPower)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
