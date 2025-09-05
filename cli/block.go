package cli

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	sbList    []*core.SubBlock
	lastBlock *core.Block
)

func init() {
	blockCmd := &cobra.Command{
		Use:   "block",
		Short: "Sub-block and block utilities",
	}

	subCreateCmd := &cobra.Command{
		Use:   "sub-create [validator] [from] [to] [amount] [fee] [nonce]",
		Args:  cobra.ExactArgs(6),
		Short: "Create a sub-block with a single transaction",
		RunE: func(cmd *cobra.Command, args []string) error {
			amt, err := strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid amount: %w", err)
			}
			fee, err := strconv.ParseUint(args[4], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid fee: %w", err)
			}
			nonce, err := strconv.ParseUint(args[5], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid nonce: %w", err)
			}
			tx := core.NewTransaction(args[1], args[2], amt, fee, nonce)
			sb := core.NewSubBlock([]*core.Transaction{tx}, args[0])
			sbList = append(sbList, sb)
			fmt.Println(sb.PohHash)
			return nil
		},
	}

	subVerifyCmd := &cobra.Command{
		Use:   "sub-verify",
		Short: "Verify the latest sub-block signature",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(sbList) == 0 {
				return fmt.Errorf("no sub-blocks")
			}
			fmt.Println(sbList[len(sbList)-1].VerifySignature())
			return nil
		},
	}

	blockCreateCmd := &cobra.Command{
		Use:   "create [prevhash]",
		Args:  cobra.ExactArgs(1),
		Short: "Create a block from existing sub-blocks",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(sbList) == 0 {
				return fmt.Errorf("no sub-blocks to assemble")
			}
			lastBlock = core.NewBlock(sbList, args[0])
			fmt.Printf("block with %d sub-blocks created\n", len(sbList))
			return nil
		},
	}

	headerCmd := &cobra.Command{
		Use:   "header [nonce]",
		Args:  cobra.ExactArgs(1),
		Short: "Compute header hash for the latest block",
		RunE: func(cmd *cobra.Command, args []string) error {
			if lastBlock == nil {
				return fmt.Errorf("no block")
			}
			nonce, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid nonce: %w", err)
			}
			fmt.Println(lastBlock.HeaderHash(nonce))
			return nil
		},
	}

	blockCmd.AddCommand(subCreateCmd, subVerifyCmd, blockCreateCmd, headerCmd)
	rootCmd.AddCommand(blockCmd)
}
