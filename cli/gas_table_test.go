package cli

import (
	"testing"

	"synnergy/core"
)

func TestGasTableSnapshot(t *testing.T) {
	snap := core.GasTableSnapshot()
	if len(snap) == 0 {
		t.Fatalf("expected snapshot to contain entries")
	}
}
