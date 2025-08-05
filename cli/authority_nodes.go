package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var authorityReg = core.NewAuthorityNodeRegistry()

func init() {
	authCmd := &cobra.Command{
		Use:   "authority",
		Short: "Manage authority nodes",
	}

	registerCmd := &cobra.Command{
		Use:   "register [address] [role]",
		Args:  cobra.ExactArgs(2),
		Short: "Register a new authority node",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := authorityReg.Register(args[0], args[1])
			if err == nil {
				fmt.Println("registered")
			}
			return err
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote [voter] [candidate]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for a candidate authority node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return authorityReg.Vote(args[0], args[1])
		},
	}

	electCmd := &cobra.Command{
		Use:   "elect [n]",
		Args:  cobra.ExactArgs(1),
		Short: "Sample an electorate of size n",
		Run: func(cmd *cobra.Command, args []string) {
			n, _ := strconv.Atoi(args[0])
			for _, addr := range authorityReg.Electorate(n) {
				fmt.Println(addr)
			}
		},
	}

	infoCmd := &cobra.Command{
		Use:   "info [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Show information about an authority node",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := authorityReg.Info(args[0])
			if err != nil {
				return err
			}
			fmt.Printf("address: %s role: %s votes: %d\n", n.Address, n.Role, len(n.Votes))
			return nil
		},
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all authority nodes",
		Run: func(cmd *cobra.Command, args []string) {
			for _, n := range authorityReg.List() {
				fmt.Printf("%s (%s) votes:%d\n", n.Address, n.Role, len(n.Votes))
			}
		},
	}

	isCmd := &cobra.Command{
		Use:   "is [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Check if address is an authority node",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(authorityReg.IsAuthorityNode(args[0]))
		},
	}

	deregCmd := &cobra.Command{
		Use:   "deregister [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Remove an authority node",
		Run: func(cmd *cobra.Command, args []string) {
			authorityReg.Deregister(args[0])
		},
	}

	authCmd.AddCommand(registerCmd, voteCmd, electCmd, infoCmd, listCmd, isCmd, deregCmd)
	rootCmd.AddCommand(authCmd)
}
