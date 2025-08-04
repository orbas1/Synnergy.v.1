package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var consensus = core.NewSynnergyConsensus()

func init() {
	consensusCmd := &cobra.Command{
		Use:   "consensus",
		Short: "Consensus operations",
	}
	mineCmd := &cobra.Command{
		Use:   "mine",
		Short: "Mine a block",
		Run: func(cmd *cobra.Command, args []string) {
			sb := core.SubBlock{Transactions: []*core.Transaction{}}
			b := core.Block{SubBlocks: []core.SubBlock{sb}}
			consensus.MineBlock(&b)
			fmt.Println("block mined")
		},
	}
	consensusCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(consensusCmd)
}
