package cli

import (
	"strings"
	"testing"
)

// TestConsensusAdaptiveWeights ensures weights command prints expected values.
func TestConsensusAdaptiveWeights(t *testing.T) {
	out, err := execCommand("consensus-adaptive", "weights")
	if err != nil {
		t.Fatalf("weights failed: %v", err)
	}
	if !strings.Contains(out, "PoW") {
		t.Fatalf("unexpected output: %s", out)
	}
}
