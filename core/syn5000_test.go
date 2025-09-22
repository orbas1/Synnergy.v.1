package core

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSYN5000PlaceResolveLifecycle(t *testing.T) {
	token := NewSYN5000Token("Gamble", "GMB", 2, WithSYN5000Config(SYN5000Config{
		MinBet:             10,
		MaxBet:             0,
		MinOddsBasisPoints: 12000,
		MaxOddsBasisPoints: 25000,
		MaxExposurePerGame: 1000,
		MaxGlobalExposure:  5000,
	}))

	placement := BetPlacement{
		Bettor:          "alice",
		Amount:          100,
		OddsBasisPoints: 15000,
		Game:            "dice",
		Metadata:        "early-morning",
	}
	betID, err := token.PlaceBetDetailed(placement)
	if err != nil {
		t.Fatalf("place bet: %v", err)
	}

	bet, ok := token.GetBet(betID)
	if !ok {
		t.Fatalf("bet %d missing", betID)
	}
	if bet.State != BetStatePending {
		t.Fatalf("expected pending state got %s", bet.State)
	}
	if bet.PotentialPayout != 150 {
		t.Fatalf("expected payout 150 got %d", bet.PotentialPayout)
	}
	if bet.Liability != 50 {
		t.Fatalf("expected liability 50 got %d", bet.Liability)
	}

	snapshot := token.Snapshot()
	if snapshot.GlobalExposure != bet.Liability {
		t.Fatalf("expected exposure %d got %d", bet.Liability, snapshot.GlobalExposure)
	}
	if snapshot.Pending != 1 || snapshot.TotalBets != 1 {
		t.Fatalf("unexpected snapshot %+v", snapshot)
	}

	payout, err := token.ResolveBet(betID, true)
	if err != nil {
		t.Fatalf("resolve bet: %v", err)
	}
	if payout != bet.PotentialPayout {
		t.Fatalf("expected payout %d got %d", bet.PotentialPayout, payout)
	}

	post, ok := token.GetBet(betID)
	if !ok || post.State != BetStateWon {
		t.Fatalf("expected resolved bet got %#v", post)
	}
	if !post.ResolvedAt.After(post.CreatedAt) {
		t.Fatalf("expected resolved timestamp to be set")
	}

	if token.Snapshot().GlobalExposure != 0 {
		t.Fatalf("expected exposure reset")
	}
}

func TestSYN5000ExposureLimits(t *testing.T) {
	cfg := SYN5000Config{
		MinBet:             5,
		MaxBet:             500,
		MinOddsBasisPoints: 11000,
		MaxOddsBasisPoints: 20000,
		MaxExposurePerGame: 9,
		MaxGlobalExposure:  9,
	}
	token := NewSYN5000Token("Risk", "RSK", 0, WithSYN5000Config(cfg))

	placement := BetPlacement{Bettor: "bob", Amount: 10, OddsBasisPoints: 15000, Game: "roulette"}
	if _, err := token.PlaceBetDetailed(placement); err != nil {
		t.Fatalf("first bet: %v", err)
	}
	if _, err := token.PlaceBetDetailed(placement); err == nil {
		t.Fatalf("expected exposure failure")
	} else if !errorsIs(err, ErrExposureLimit) {
		t.Fatalf("expected ErrExposureLimit got %v", err)
	}

	bets := token.ListBets(BetFilter{})
	if len(bets) != 1 {
		t.Fatalf("expected 1 bet recorded got %d", len(bets))
	}

	if err := token.CancelBet(bets[0].ID, "risk reset"); err != nil {
		t.Fatalf("cancel bet: %v", err)
	}
	if token.Snapshot().GlobalExposure != 0 {
		t.Fatalf("exposure should be zero after cancel")
	}
}

func TestSYN5000SignatureValidation(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	token := NewSYN5000Token("Secure", "SEC", 0, WithSYN5000Config(SYN5000Config{
		MinBet:             1,
		MinOddsBasisPoints: 12000,
		RequireSignature:   true,
	}), WithSYN5000Participant("carol", pub))

	placement := BetPlacement{
		Bettor:          "carol",
		Amount:          50,
		OddsBasisPoints: 16000,
		Game:            "slots",
		Nonce:           "nonce-1",
	}
	digest := placement.signingDigest()
	signature := ed25519.Sign(priv, digest[:])

	betID, err := token.PlaceBetWithSignature(placement, signature)
	if err != nil {
		t.Fatalf("signed placement: %v", err)
	}

	badSig := make([]byte, len(signature))
	copy(badSig, signature)
	badSig[0] ^= 0xFF
	placement.Nonce = "nonce-2"
	if _, err := token.PlaceBetWithSignature(placement, badSig); err == nil {
		t.Fatalf("expected invalid signature error")
	}

	if _, err := token.PlaceBetDetailed(BetPlacement{Bettor: "carol", Amount: 10, OddsBasisPoints: 13000, Game: "slots"}); err == nil {
		t.Fatalf("expected unsigned placement to fail when signatures required")
	}

	if _, err := token.ResolveBet(betID, false); err != nil {
		t.Fatalf("resolve signed bet: %v", err)
	}
}

func TestSYN5000ListFiltering(t *testing.T) {
	token := NewSYN5000Token("Filter", "FLT", 0)
	ids := make([]uint64, 0, 3)
	base := BetPlacement{Amount: 20, OddsBasisPoints: 14000, Game: "cards"}
	for i := 0; i < 3; i++ {
		placement := base
		placement.Bettor = []string{"alice", "bob", "alice"}[i]
		placement.Game = []string{"cards", "dice", "cards"}[i]
		placement.Metadata = "batch"
		id, err := token.PlaceBetDetailed(placement)
		if err != nil {
			t.Fatalf("place bet %d: %v", i, err)
		}
		ids = append(ids, id)
	}

	if _, err := token.ResolveBet(ids[0], true); err != nil {
		t.Fatalf("resolve bet: %v", err)
	}

	bets := token.ListBets(BetFilter{Bettor: "alice", State: BetStatePending})
	if len(bets) != 1 {
		t.Fatalf("expected 1 pending alice bet got %d", len(bets))
	}

	bets = token.ListBets(BetFilter{Game: "cards"})
	if len(bets) != 2 {
		t.Fatalf("expected 2 card bets got %d", len(bets))
	}
}

func TestSYN5000ConcurrentPlacements(t *testing.T) {
	token := NewSYN5000Token("Stress", "STR", 0, WithSYN5000Config(SYN5000Config{
		MaxExposurePerGame: 0,
		MaxGlobalExposure:  0,
	}))
	base := BetPlacement{Amount: 5, OddsBasisPoints: 15000, Game: "arena"}

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			placement := base
			placement.Bettor = fmt.Sprintf("user-%d", idx)
			if _, err := token.PlaceBetDetailed(placement); err != nil {
				t.Errorf("placement %d failed: %v", idx, err)
			}
		}(i)
	}
	wg.Wait()

	snapshot := token.Snapshot()
	if snapshot.TotalBets == 0 {
		t.Fatalf("expected bets recorded after concurrent placement")
	}
	if snapshot.Pending != snapshot.TotalBets {
		t.Fatalf("expected all bets pending in snapshot")
	}
}

func errorsIs(err, target error) bool {
	return err != nil && target != nil && (err == target || errors.Is(err, target))
}

func init() {
	// Seed monotonic timers for tests relying on timestamp ordering.
	time.Now()
}
