package core

import (
	"fmt"
	"sync"
	"testing"
)

func TestMusicToken(t *testing.T) {
	m := NewMusicToken("Song", "Artist", "Album")
	title, artist, album := m.Info()
	if title != "Song" || artist != "Artist" || album != "Album" {
		t.Fatalf("unexpected metadata: %s %s %s", title, artist, album)
	}

	m.Update("NewSong", "", "")
	if t2, _, _ := m.Info(); t2 != "NewSong" {
		t.Fatalf("expected title update, got %s", t2)
	}

	m.SetRoyaltyShare("addr1", 1)
	m.SetRoyaltyShare("addr2", 1)
	if share, ok := m.RoyaltyShare("addr1"); !ok || share != 1 {
		t.Fatalf("expected share of 1, got %d", share)
	}
	payouts, err := m.Distribute(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payouts["addr1"] != 50 || payouts["addr2"] != 50 {
		t.Fatalf("incorrect payouts: %#v", payouts)
	}

	if err := m.RemoveRoyaltyRecipient("addr1"); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
	if _, ok := m.RoyaltyShare("addr1"); ok {
		t.Fatalf("expected addr1 removal")
	}
}

func TestMusicTokenErrors(t *testing.T) {
	m := NewMusicToken("Song", "Artist", "Album")
	if _, err := m.Distribute(100); err != ErrNoRoyaltyRecipients {
		t.Fatalf("expected ErrNoRoyaltyRecipients, got %v", err)
	}
	if err := m.RemoveRoyaltyRecipient("addrX"); err != ErrRecipientNotFound {
		t.Fatalf("expected ErrRecipientNotFound, got %v", err)
	}
}

func TestMusicTokenConcurrentAccess(t *testing.T) {
	m := NewMusicToken("Song", "Artist", "Album")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			addr := fmt.Sprintf("addr%d", i)
			m.SetRoyaltyShare(addr, 1)
		}(i)
	}
	wg.Wait()
	payouts, err := m.Distribute(1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(payouts) != 100 {
		t.Fatalf("expected 100 payouts, got %d", len(payouts))
	}
}
