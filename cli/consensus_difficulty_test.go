package cli

import "testing"

// TestConsensusDifficultySample ensures sample command outputs a value.
func TestConsensusDifficultySample(t *testing.T) {
	out, err := execCommand("consensus-difficulty", "sample", "2")
	if err != nil {
		t.Fatalf("sample failed: %v", err)
	}
	if out == "" {
		t.Fatalf("expected difficulty output")
	}
}
