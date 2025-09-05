package cli

import (
	"strings"
	"testing"
)

// TestConsensusWeights ensures consensus weights command prints expected values.
func TestConsensusWeights(t *testing.T) {
	out, err := execCommand("consensus", "weights")
	if err != nil {
		t.Fatalf("weights failed: %v", err)
	}
	if !strings.Contains(out, "PoW") || !strings.Contains(out, "PoS") {
		t.Fatalf("unexpected output: %s", out)
	}
}
