package cli

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	authorityValidators = core.NewValidatorManager(core.MinStake)
	authorityReg        = core.NewAuthorityNodeRegistry(ledger, authorityValidators, 1)
)

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
			sigHex, _ := cmd.Flags().GetString("sig")
			pubHex, _ := cmd.Flags().GetString("pub")
			sig, err := hex.DecodeString(sigHex)
			if err != nil {
				return err
			}
			pubBytes, err := hex.DecodeString(pubHex)
			if err != nil {
				return err
			}
			return authorityReg.Vote(args[0], args[1], sig, ed25519.PublicKey(pubBytes))
		},
	}
	voteCmd.Flags().String("sig", "", "hex-encoded signature of candidate address")
	voteCmd.Flags().String("pub", "", "hex-encoded public key of voter")

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

	var infoJSON bool
	var listJSON bool

	infoCmd := &cobra.Command{
		Use:   "info [address]",
		Args:  cobra.ExactArgs(1),
		Short: "Show information about an authority node",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := authorityReg.Info(args[0])
			if err != nil {
				return err
			}
			if infoJSON {
				enc, _ := json.Marshal(n)
				fmt.Println(string(enc))
			} else {
				fmt.Printf("address: %s role: %s votes: %d\n", n.Address, n.Role, len(n.Votes))
			}
			return nil
		},
	}
	infoCmd.Flags().BoolVar(&infoJSON, "json", false, "output as JSON")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all authority nodes",
		Run: func(cmd *cobra.Command, args []string) {
			nodes := authorityReg.List()
			if listJSON {
				enc, _ := json.Marshal(nodes)
				fmt.Println(string(enc))
				return
			}
			for _, n := range nodes {
				fmt.Printf("%s (%s) votes:%d\n", n.Address, n.Role, len(n.Votes))
			}
		},
	}
	listCmd.Flags().BoolVar(&listJSON, "json", false, "output as JSON")

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
