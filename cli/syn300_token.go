package cli

import (
	"errors"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			balStr, _ := cmd.Flags().GetString("balances")
			if balStr == "" {
				return errors.New("balances required")
			}
			syn300 = core.NewSYN300Token(parseBalances(balStr))
			fmt.Println("token initialised")
			return nil
		},
	}
	initCmd.Flags().String("balances", "", "initial balances addr=amt,comma-separated")
	cmd.AddCommand(initCmd)

	delCmd := &cobra.Command{
		Use:   "delegate <owner> <delegate>",
		Short: "Delegate voting power",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			if err := syn300.Delegate(args[0], args[1]); err != nil {
				return err
			}
			fmt.Println("delegated")
			return nil
		},
	}
	cmd.AddCommand(delCmd)

	revokeCmd := &cobra.Command{
		Use:   "revoke <owner>",
		Short: "Revoke delegation",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			if err := syn300.RevokeDelegation(args[0]); err != nil {
				return err
			}
			fmt.Println("revoked")
			return nil
		},
	}
	cmd.AddCommand(revokeCmd)

	powerCmd := &cobra.Command{
		Use:   "power <addr>",
		Short: "Show voting power",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			fmt.Println(syn300.VotingPower(args[0]))
			return nil
		},
	}
	cmd.AddCommand(powerCmd)

	propCmd := &cobra.Command{
		Use:   "propose <creator> <desc>",
		Short: "Create proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			id, err := syn300.CreateProposal(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Printf("proposal %d created\n", id)
			return nil
		},
	}
	cmd.AddCommand(propCmd)

	voteCmd := &cobra.Command{
		Use:   "vote <id> <voter> <approve>",
		Short: "Cast vote",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			var id uint64
			if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
				return err
			}
			approve := args[2] == "true"
			if err := syn300.Vote(id, args[1], approve); err != nil {
				return err
			}
			fmt.Println("vote recorded")
			return nil
		},
	}
	cmd.AddCommand(voteCmd)

	execCmd := &cobra.Command{
		Use:   "execute <id> <quorum>",
		Short: "Execute proposal if quorum met",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			var id, quorum uint64
			if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
				return err
			}
			if _, err := fmt.Sscanf(args[1], "%d", &quorum); err != nil {
				return err
			}
			if err := syn300.Execute(id, quorum); err != nil {
				return err
			}
			fmt.Println("proposal executed")
			return nil
		},
	}
	cmd.AddCommand(execCmd)

	statusCmd := &cobra.Command{
		Use:   "status <id>",
		Short: "Show proposal status",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			var id uint64
			if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
				return err
			}
			p, err := syn300.ProposalStatus(id)
			if err != nil {
				return err
			}
			fmt.Printf("%+v\n", *p)
			return nil
		},
	}
	cmd.AddCommand(statusCmd)

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List proposals",
		RunE: func(cmd *cobra.Command, args []string) error {
			if syn300 == nil {
				return errors.New("token not initialised")
			}
			for _, p := range syn300.ListProposals() {
				fmt.Printf("%d %s executed:%v\n", p.ID, p.Description, p.Executed)
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd)

	rootCmd.AddCommand(cmd)
}
