package synnergy

import "testing"

func TestSplitAndReconstructHolographic(t *testing.T) {
	data := []byte("hello world")
	frame := SplitHolographic("id", data, 3)
	if len(frame.Shards) != 3 {
		t.Fatalf("expected 3 shards, got %d", len(frame.Shards))
	}
	recon := ReconstructHolographic(frame)
	if string(recon) != string(data) {
		t.Fatalf("reconstructed data mismatch: %s", string(recon))
	}
	empty := SplitHolographic("id2", data, 0)
	if len(empty.Shards) != 0 {
		t.Fatalf("expected no shards for n<=0")
	}
}

func TestHolographicShardIsolation(t *testing.T) {
	data := []byte("enterprise-grade")
	frame := SplitHolographic("iso", data, 4)
	if len(frame.Shards) != 4 {
		t.Fatalf("expected 4 shards")
	}
	frame.Shards[0][0] = 'X'
	recon := ReconstructHolographic(frame)
	if string(recon[:len(data)]) == string(data) {
		t.Fatalf("expected reconstruction to reflect shard mutation")
	}
}
