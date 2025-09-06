package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestConsensusWeights ensures consensus weights command prints expected values.
func TestConsensusWeights(t *testing.T) {
	out, err := execCommand("--json", "consensus", "weights")
	if err != nil {
		t.Fatalf("weights failed: %v", err)
	}
	start := strings.LastIndex(out, "{")
	end := strings.LastIndex(out, "}")
	if start != -1 && end != -1 {
		out = out[start : end+1]
	}
	var res map[string]float64
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := res["pow"]; !ok {
		t.Fatalf("missing pow weight: %v", res)
	}
}
