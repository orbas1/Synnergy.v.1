package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "syn5000_index",
		Short: "Inspect registered SYN5000 tokens",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List registered gambling tokens",
		Run: func(cmd *cobra.Command, args []string) {
			symbols := syn5000TokenIndex.Symbols()
			if len(symbols) == 0 {
				cmd.Println("no tokens registered")
				return
			}
			for _, sym := range symbols {
				cmd.Println(sym)
			}
		},
	}

	summaryCmd := &cobra.Command{
		Use:   "summary",
		Short: "Show exposure summary for all tokens",
		Run: func(cmd *cobra.Command, args []string) {
			symbols := syn5000TokenIndex.Symbols()
			if len(symbols) == 0 {
				cmd.Println("no tokens registered")
				return
			}
			sort.Strings(symbols)
			for _, sym := range symbols {
				snap, ok := syn5000TokenIndex.Snapshot(sym)
				if !ok {
					cmd.Printf("%s: unavailable\n", sym)
					continue
				}
				cmd.Printf("%s pending=%d won=%d lost=%d cancelled=%d exposure=%d\n", sym, snap.Pending, snap.Won, snap.Lost, snap.Cancelled, snap.GlobalExposure)
			}
		},
	}

	detailCmd := &cobra.Command{
		Use:   "detail",
		Short: "Show snapshot for a specific token",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("id")
			if symbol == "" {
				return fmt.Errorf("id required")
			}
			snap, ok := syn5000TokenIndex.Snapshot(symbol)
			if !ok {
				return fmt.Errorf("token not found")
			}
			cmd.Printf("Token %s (%s)\n", snap.Name, snap.Symbol)
			cmd.Printf("Pending=%d Won=%d Lost=%d Cancelled=%d Exposure=%d\n", snap.Pending, snap.Won, snap.Lost, snap.Cancelled, snap.GlobalExposure)
			if len(snap.ExposureByGame) == 0 {
				cmd.Println("Exposure by game: none")
				return nil
			}
			games := make([]string, 0, len(snap.ExposureByGame))
			for game := range snap.ExposureByGame {
				games = append(games, game)
			}
			sort.Strings(games)
			for _, game := range games {
				cmd.Printf("  %s: %d\n", game, snap.ExposureByGame[game])
			}
			return nil
		},
	}
	detailCmd.Flags().String("id", "", "token symbol")

	cmd.AddCommand(listCmd, summaryCmd, detailCmd)
	rootCmd.AddCommand(cmd)
}
