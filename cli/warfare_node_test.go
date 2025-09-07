package cli

import (
	"testing"

	"synnergy/core"
)

func TestWarfareNodeLogistics(t *testing.T) {
	base := core.NewNode("n1", "addr", core.NewLedger())
	wn := core.NewWarfareNode(base)

	if err := wn.SecureCommand("move"); err != nil {
		t.Fatalf("secure command failed: %v", err)
	}

	wn.TrackLogistics("asset", "loc", "ok")
	logs := wn.Logistics()
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}

	assetLogs := wn.LogisticsByAsset("asset")
	if len(assetLogs) != 1 {
		t.Fatalf("expected asset log")
	}

	wn.ShareTactical("info")
}
