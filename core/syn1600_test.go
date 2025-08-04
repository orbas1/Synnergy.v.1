package core

import "testing"

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
	payouts, err := m.Distribute(100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payouts["addr1"] != 50 || payouts["addr2"] != 50 {
		t.Fatalf("incorrect payouts: %#v", payouts)
	}
}
