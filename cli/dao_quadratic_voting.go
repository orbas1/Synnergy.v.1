package cli

import (
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

	weightCmd := &cobra.Command{
		Use:   "weight <tokens>",
		Args:  cobra.ExactArgs(1),
		Short: "Calculate quadratic weight",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("QuadraticWeight")
			tokens, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				fmt.Println("invalid tokens")
				return
			}
			fmt.Println(core.QuadraticWeight(tokens))
		},
	}

	voteCmd := &cobra.Command{
		Use:   "vote <id> <voter> <tokens> <yes|no>",
		Args:  cobra.ExactArgs(4),
		Short: "Cast a quadratic vote",
		Run: func(cmd *cobra.Command, args []string) {
			gasPrint("CastQuadraticVote")
			tokens, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				fmt.Println("invalid tokens")
				return
			}
			support := strings.ToLower(args[3]) == "yes"
			if err := proposalMgr.CastQuadraticVote(args[0], args[1], tokens, support); err != nil {
				fmt.Println(err)
			}
		},
	}

	qvCmd.AddCommand(weightCmd, voteCmd)
	rootCmd.AddCommand(qvCmd)
}
