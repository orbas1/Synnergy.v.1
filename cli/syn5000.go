package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var syn5000Tokens = map[string]*core.SYN5000Token{}

func init() {
	cmd := &cobra.Command{
		Use:   "syn5000",
		Short: "SYN5000 gambling token",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a gambling token",
		Run: func(cmd *cobra.Command, args []string) {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint("dec")
			t := core.NewSYN5000Token(name, symbol, uint8(dec))
			syn5000Tokens[symbol] = t
			fmt.Println("token created", symbol)
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().Uint("dec", 0, "decimals")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")

	betCmd := &cobra.Command{
		Use:   "bet <bettor>",
		Args:  cobra.ExactArgs(1),
		Short: "Place a bet",
		Run: func(cmd *cobra.Command, args []string) {
			tokenID, _ := cmd.Flags().GetString("id")
			amt, _ := cmd.Flags().GetUint64("amt")
			odds, _ := cmd.Flags().GetFloat64("odds")
			game, _ := cmd.Flags().GetString("game")
			t, ok := syn5000Tokens[tokenID]
			if !ok {
				fmt.Println("token not found")
				return
			}
			id := t.PlaceBet(args[0], amt, odds, game)
			fmt.Println("bet placed", id)
		},
	}
	betCmd.Flags().String("id", "", "token symbol")
	betCmd.Flags().Uint64("amt", 0, "bet amount")
	betCmd.Flags().Float64("odds", 1, "betting odds")
	betCmd.Flags().String("game", "", "game type")
	betCmd.MarkFlagRequired("id")
	betCmd.MarkFlagRequired("amt")
	betCmd.MarkFlagRequired("odds")
	betCmd.MarkFlagRequired("game")

	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a bet",
		Run: func(cmd *cobra.Command, args []string) {
			tokenID, _ := cmd.Flags().GetString("id")
			betID, _ := cmd.Flags().GetUint64("bet")
			win, _ := cmd.Flags().GetBool("win")
			t, ok := syn5000Tokens[tokenID]
			if !ok {
				fmt.Println("token not found")
				return
			}
			payout, err := t.ResolveBet(betID, win)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
			fmt.Println("payout", payout)
		},
	}
	resolveCmd.Flags().String("id", "", "token symbol")
	resolveCmd.Flags().Uint64("bet", 0, "bet id")
	resolveCmd.Flags().Bool("win", false, "mark bet as won")
	resolveCmd.MarkFlagRequired("id")
	resolveCmd.MarkFlagRequired("bet")

	cmd.AddCommand(createCmd, betCmd, resolveCmd)
	rootCmd.AddCommand(cmd)
}
