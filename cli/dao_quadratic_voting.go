package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

func init() {
	qvCmd := &cobra.Command{
		Use:   "dao-qv",
		Short: "Quadratic voting operations",
	}

	var weightJSON bool
	weightCmd := &cobra.Command{
		Use:   "weight <tokens>",
		Args:  cobra.ExactArgs(1),
		Short: "Calculate quadratic weight",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("QuadraticWeight")
			tokens, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid tokens")
				return
			}
			w := core.QuadraticWeight(tokens)
			if weightJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]uint64{"weight": w})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), w)
		},
	}
	weightCmd.Flags().BoolVar(&weightJSON, "json", false, "output as JSON")

	var voteJSON bool
	var votePub, voteMsg, voteSig string
	voteCmd := &cobra.Command{
		Use:   "vote <id> <voter> <tokens> <yes|no>",
		Args:  cobra.ExactArgs(4),
		Short: "Cast a quadratic vote",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("CastQuadraticVote")
			ok, err := VerifySignature(votePub, voteMsg, voteSig)
			if err != nil || !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "signature verification failed")
				return
			}
			tokens, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), "invalid tokens")
				return
			}
			support := strings.ToLower(args[3]) == "yes"
			if err := proposalMgr.CastQuadraticVote(args[0], args[1], tokens, support); err != nil {
				fmt.Fprintln(cmd.OutOrStdout(), err)
				return
			}
			if voteJSON {
				_ = json.NewEncoder(cmd.OutOrStdout()).Encode(map[string]string{"status": "vote recorded"})
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), "vote recorded")
		},
	}
	voteCmd.Flags().BoolVar(&voteJSON, "json", false, "output as JSON")
	voteCmd.Flags().StringVar(&votePub, "pub", "", "hex encoded public key")
	voteCmd.Flags().StringVar(&voteMsg, "msg", "", "hex encoded message")
	voteCmd.Flags().StringVar(&voteSig, "sig", "", "hex encoded signature")

	qvCmd.AddCommand(weightCmd, voteCmd)
	rootCmd.AddCommand(qvCmd)
}
