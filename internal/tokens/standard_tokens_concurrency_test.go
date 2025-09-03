package tokens

import (
	"fmt"
	"sync"
	"testing"
)

// TestSYN20ConcurrentPause exercises the pause and unpause controls from
// multiple goroutines to ensure the additional mutex does not race.
func TestSYN20ConcurrentPause(t *testing.T) {
	tok := NewSYN20Token(99, "Util", "UTL", 0)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tok.Pause()
			tok.Unpause()
		}()
	}
	wg.Wait()
	if err := tok.Mint("alice", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestSYN70ConcurrentRegister registers assets concurrently verifying that the
// internal map is protected.
func TestSYN70ConcurrentRegister(t *testing.T) {
	tok := NewSYN70Token(100, "Game", "G", 0)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("a%d", i)
			tok.RegisterAsset(id, "owner", "Item", "Game")
		}(i)
	}
	wg.Wait()
	assets := tok.ListAssets()
	if len(assets) != 5 {
		t.Fatalf("expected 5 assets got %d", len(assets))
	}
}
