package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn300 *core.SYN300Token

func parseBalances(s string) map[string]uint64 {
	m := make(map[string]uint64)
	if s == "" {
		return m
	}
	for _, kv := range strings.Split(s, ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		var amt uint64
		fmt.Sscanf(parts[1], "%d", &amt)
		m[parts[0]] = amt
	}
	return m
}

func init() {
	cmd := &cobra.Command{Use: "syn300", Short: "SYN300 governance token"}

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise token with balances",
		Run: func(cmd *cobra.Command, args []string) {
			balStr, _ := cmd.Flags().GetString("balances")
			syn300 = core.NewSYN300Token(parseBalances(balStr))
			fmt.Println("token initialised")
		},
	}
	initCmd.Flags().String("balances", "", "initial balances addr=amt,comma-separated")
	cmd.AddCommand(initCmd)

	delCmd := &cobra.Command{
		Use:   "delegate <owner> <delegate>",
		Short: "Delegate voting power",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			syn300.Delegate(args[0], args[1])
			fmt.Println("delegated")
		},
	}
	cmd.AddCommand(delCmd)

	revokeCmd := &cobra.Command{
		Use:   "revoke <owner>",
		Short: "Revoke delegation",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			syn300.RevokeDelegation(args[0])
			fmt.Println("revoked")
		},
	}
	cmd.AddCommand(revokeCmd)

	powerCmd := &cobra.Command{
		Use:   "power <addr>",
		Short: "Show voting power",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			fmt.Println(syn300.VotingPower(args[0]))
		},
	}
	cmd.AddCommand(powerCmd)

	propCmd := &cobra.Command{
		Use:   "propose <creator> <desc>",
		Short: "Create proposal",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			id := syn300.CreateProposal(args[0], args[1])
			fmt.Printf("proposal %d created\n", id)
		},
	}
	cmd.AddCommand(propCmd)

	voteCmd := &cobra.Command{
		Use:   "vote <id> <voter> <approve>",
		Short: "Cast vote",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			approve := args[2] == "true"
			if err := syn300.Vote(id, args[1], approve); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("vote recorded")
			}
		},
	}
	cmd.AddCommand(voteCmd)

	execCmd := &cobra.Command{
		Use:   "execute <id> <quorum>",
		Short: "Execute proposal if quorum met",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			var id, quorum uint64
			fmt.Sscanf(args[0], "%d", &id)
			fmt.Sscanf(args[1], "%d", &quorum)
			if err := syn300.Execute(id, quorum); err != nil {
				fmt.Printf("error: %v\n", err)
			} else {
				fmt.Println("proposal executed")
			}
		},
	}
	cmd.AddCommand(execCmd)

	statusCmd := &cobra.Command{
		Use:   "status <id>",
		Short: "Show proposal status",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			var id uint64
			fmt.Sscanf(args[0], "%d", &id)
			p, err := syn300.ProposalStatus(id)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return
			}
			fmt.Printf("%+v\n", *p)
		},
	}
	cmd.AddCommand(statusCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		Run: func(cmd *cobra.Command, args []string) {
			if syn300 == nil {
				fmt.Println("token not initialised")
				return
			}
			for _, p := range syn300.ListProposals() {
				fmt.Printf("%d %s executed:%v\n", p.ID, p.Description, p.Executed)
			}
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
