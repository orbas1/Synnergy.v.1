package cli

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"synnergy/core"
)

var (
	syn5000Tokens     = map[string]*core.SYN5000Token{}
	syn5000TokenIndex = core.NewSYN5000Index()
)

func init() {
	cmd := &cobra.Command{
		Use:   "syn5000",
		Short: "Manage SYN5000 gambling tokens",
	}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a SYN5000 token",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			symbol, _ := cmd.Flags().GetString("symbol")
			dec, _ := cmd.Flags().GetUint("dec")
			if name == "" || symbol == "" {
				return fmt.Errorf("name and symbol required")
			}
			t := core.NewSYN5000Token(name, symbol, uint8(dec))

			cfg := core.DefaultSYN5000Config()
			minBet, _ := cmd.Flags().GetUint64("min-bet")
			maxBet, _ := cmd.Flags().GetUint64("max-bet")
			minOdds, _ := cmd.Flags().GetFloat64("min-odds")
			maxOdds, _ := cmd.Flags().GetFloat64("max-odds")
			maxGameExposure, _ := cmd.Flags().GetUint64("max-game-exposure")
			maxGlobalExposure, _ := cmd.Flags().GetUint64("max-global-exposure")
			requireSignature, _ := cmd.Flags().GetBool("require-signature")

			if minBet > 0 {
				cfg.MinBet = minBet
			}
			if maxBet > 0 {
				cfg.MaxBet = maxBet
			}
			if minOdds > 0 {
				bp, err := core.ConvertOddsToBasisPoints(minOdds)
				if err != nil {
					return fmt.Errorf("convert min-odds: %w", err)
				}
				cfg.MinOddsBasisPoints = bp
			}
			if maxOdds > 0 {
				bp, err := core.ConvertOddsToBasisPoints(maxOdds)
				if err != nil {
					return fmt.Errorf("convert max-odds: %w", err)
				}
				cfg.MaxOddsBasisPoints = bp
			}
			if maxGameExposure > 0 {
				cfg.MaxExposurePerGame = maxGameExposure
			}
			if maxGlobalExposure > 0 {
				cfg.MaxGlobalExposure = maxGlobalExposure
			}
			cfg.RequireSignature = requireSignature
			t.SetConfig(cfg)

			syn5000Tokens[symbol] = t
			if err := syn5000TokenIndex.Register(symbol, t); err != nil {
				return err
			}
			cmd.Printf("token created %s\n", symbol)
			return nil
		},
	}
	createCmd.Flags().String("name", "", "token name")
	createCmd.Flags().String("symbol", "", "token symbol")
	createCmd.Flags().Uint("dec", 0, "token decimals")
	createCmd.Flags().Uint64("min-bet", 0, "minimum bet amount")
	createCmd.Flags().Uint64("max-bet", 0, "maximum bet amount")
	createCmd.Flags().Float64("min-odds", 0, "minimum odds multiplier")
	createCmd.Flags().Float64("max-odds", 0, "maximum odds multiplier")
	createCmd.Flags().Uint64("max-game-exposure", 0, "maximum liability per game")
	createCmd.Flags().Uint64("max-global-exposure", 0, "maximum global liability")
	createCmd.Flags().Bool("require-signature", false, "require ed25519 signatures for bets")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("symbol")

	registerCmd := &cobra.Command{
		Use:   "register <bettor>",
		Args:  cobra.ExactArgs(1),
		Short: "Register a bettor's public key",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("id")
			token, err := loadSYN5000Token(symbol)
			if err != nil {
				return err
			}
			pubKeyStr, _ := cmd.Flags().GetString("pubkey")
			if pubKeyStr == "" {
				return fmt.Errorf("pubkey required")
			}
			raw, err := base64.StdEncoding.DecodeString(pubKeyStr)
			if err != nil {
				return fmt.Errorf("decode pubkey: %w", err)
			}
			if len(raw) != ed25519.PublicKeySize {
				return fmt.Errorf("expected %d byte ed25519 key", ed25519.PublicKeySize)
			}
			if err := token.RegisterParticipant(args[0], ed25519.PublicKey(raw)); err != nil {
				return err
			}
			cmd.Printf("bettor %s registered\n", args[0])
			return nil
		},
	}
	registerCmd.Flags().String("id", "", "token symbol")
	registerCmd.Flags().String("pubkey", "", "bettor ed25519 public key (base64)")
	registerCmd.MarkFlagRequired("id")
	registerCmd.MarkFlagRequired("pubkey")

	betCmd := &cobra.Command{
		Use:   "bet <bettor>",
		Args:  cobra.ExactArgs(1),
		Short: "Place a bet",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("id")
			token, err := loadSYN5000Token(symbol)
			if err != nil {
				return err
			}
			amt, _ := cmd.Flags().GetUint64("amt")
			odds, _ := cmd.Flags().GetFloat64("odds")
			game, _ := cmd.Flags().GetString("game")
			metadata, _ := cmd.Flags().GetString("metadata")
			nonce, _ := cmd.Flags().GetString("nonce")
			sigStr, _ := cmd.Flags().GetString("signature")
			if game == "" {
				return fmt.Errorf("game required")
			}
			oddsBP, err := core.ConvertOddsToBasisPoints(odds)
			if err != nil {
				return err
			}
			placement := core.BetPlacement{
				Bettor:          args[0],
				Amount:          amt,
				OddsBasisPoints: oddsBP,
				Game:            game,
				Metadata:        metadata,
				Nonce:           nonce,
			}

			var betID uint64
			if sigStr != "" {
				rawSig, err := base64.StdEncoding.DecodeString(sigStr)
				if err != nil {
					return fmt.Errorf("decode signature: %w", err)
				}
				betID, err = token.PlaceBetWithSignature(placement, rawSig)
				if err != nil {
					return err
				}
			} else {
				var err error
				betID, err = token.PlaceBetDetailed(placement)
				if err != nil {
					return err
				}
			}
			cmd.Printf("bet placed %d\n", betID)
			return nil
		},
	}
	betCmd.Flags().String("id", "", "token symbol")
	betCmd.Flags().Uint64("amt", 0, "bet amount")
	betCmd.Flags().Float64("odds", 1, "betting odds multiplier")
	betCmd.Flags().String("game", "", "game or market identifier")
	betCmd.Flags().String("metadata", "", "optional bet metadata")
	betCmd.Flags().String("nonce", "", "custom placement nonce")
	betCmd.Flags().String("signature", "", "base64 ed25519 signature")
	betCmd.MarkFlagRequired("id")
	betCmd.MarkFlagRequired("amt")
	betCmd.MarkFlagRequired("odds")
	betCmd.MarkFlagRequired("game")

	resolveCmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a bet",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("id")
			token, err := loadSYN5000Token(symbol)
			if err != nil {
				return err
			}
			betID, _ := cmd.Flags().GetUint64("bet")
			win, _ := cmd.Flags().GetBool("win")
			note, _ := cmd.Flags().GetString("note")
			payout, err := token.ResolveBetWithNote(betID, win, note)
			if err != nil {
				return err
			}
			cmd.Printf("payout %d\n", payout)
			return nil
		},
	}
	resolveCmd.Flags().String("id", "", "token symbol")
	resolveCmd.Flags().Uint64("bet", 0, "bet identifier")
	resolveCmd.Flags().Bool("win", false, "mark bet as winning")
	resolveCmd.Flags().String("note", "", "resolution note")
	resolveCmd.MarkFlagRequired("id")
	resolveCmd.MarkFlagRequired("bet")

	cancelCmd := &cobra.Command{
		Use:   "cancel",
		Short: "Cancel a pending bet",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("id")
			token, err := loadSYN5000Token(symbol)
			if err != nil {
				return err
			}
			betID, _ := cmd.Flags().GetUint64("bet")
			note, _ := cmd.Flags().GetString("note")
			if err := token.CancelBet(betID, note); err != nil {
				return err
			}
			cmd.Println("bet cancelled")
			return nil
		},
	}
	cancelCmd.Flags().String("id", "", "token symbol")
	cancelCmd.Flags().Uint64("bet", 0, "bet identifier")
	cancelCmd.Flags().String("note", "", "cancellation note")
	cancelCmd.MarkFlagRequired("id")
	cancelCmd.MarkFlagRequired("bet")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List bets for a token",
		RunE: func(cmd *cobra.Command, args []string) error {
			symbol, _ := cmd.Flags().GetString("id")
			if symbol == "" {
				return fmt.Errorf("id required")
			}
			bettor, _ := cmd.Flags().GetString("bettor")
			game, _ := cmd.Flags().GetString("game")
			stateStr, _ := cmd.Flags().GetString("state")
			limit, _ := cmd.Flags().GetInt("limit")
			var state core.BetState
			if stateStr != "" {
				switch strings.ToLower(stateStr) {
				case string(core.BetStatePending):
					state = core.BetStatePending
				case string(core.BetStateWon):
					state = core.BetStateWon
				case string(core.BetStateLost):
					state = core.BetStateLost
				case string(core.BetStateCancelled):
					state = core.BetStateCancelled
				default:
					return fmt.Errorf("unknown state %s", stateStr)
				}
			}
			bets, ok := syn5000TokenIndex.Bets(symbol, core.BetFilter{Bettor: bettor, Game: game, State: state, Limit: limit})
			if !ok {
				return fmt.Errorf("token not found")
			}
			sort.Slice(bets, func(i, j int) bool { return bets[i].ID < bets[j].ID })
			for _, bet := range bets {
				cmd.Printf("%d %s %s %d liability=%d state=%s\n", bet.ID, bet.Placement.Bettor, bet.Placement.Game, bet.Placement.Amount, bet.Liability, bet.State)
			}
			if len(bets) == 0 {
				cmd.Println("no bets found")
			}
			return nil
		},
	}
	listCmd.Flags().String("id", "", "token symbol")
	listCmd.Flags().String("bettor", "", "filter by bettor")
	listCmd.Flags().String("game", "", "filter by game")
	listCmd.Flags().String("state", "", "filter by state (pending|won|lost|cancelled)")
	listCmd.Flags().Int("limit", 0, "maximum results")

	snapshotCmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Show token exposure snapshot",
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
			cmd.Printf("Total Bets: %d (pending=%d won=%d lost=%d cancelled=%d)\n", snap.TotalBets, snap.Pending, snap.Won, snap.Lost, snap.Cancelled)
			cmd.Printf("Global Exposure: %d\n", snap.GlobalExposure)
			if len(snap.ExposureByGame) == 0 {
				cmd.Println("Exposure By Game: none")
			} else {
				cmd.Println("Exposure By Game:")
				games := make([]string, 0, len(snap.ExposureByGame))
				for game := range snap.ExposureByGame {
					games = append(games, game)
				}
				sort.Strings(games)
				for _, game := range games {
					cmd.Printf("  %s: %d\n", game, snap.ExposureByGame[game])
				}
			}
			return nil
		},
	}
	snapshotCmd.Flags().String("id", "", "token symbol")

	cmd.AddCommand(createCmd, registerCmd, betCmd, resolveCmd, cancelCmd, listCmd, snapshotCmd)
	rootCmd.AddCommand(cmd)
}

func loadSYN5000Token(symbol string) (*core.SYN5000Token, error) {
	if symbol == "" {
		return nil, fmt.Errorf("id required")
	}
	token, ok := syn5000Tokens[symbol]
	if !ok {
		return nil, fmt.Errorf("token %s not found", symbol)
	}
	return token, nil
}
