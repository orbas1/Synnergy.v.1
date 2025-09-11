package core

import (
	"testing"

	"synnergy/internal/nodes"
)

func TestLightNodeHeaders(t *testing.T) {
	ln := NewLightNode(nodes.Address("l1"))
	h := nodes.BlockHeader{Hash: "h1", Height: 1}
	if err := ln.AddHeader(h); err != nil {
		t.Fatalf("add header: %v", err)
	}
	if latest, ok := ln.LatestHeader(); !ok || latest.Hash != "h1" {
		t.Fatalf("unexpected latest header: %+v", latest)
	}
	headers := ln.Headers()
	if len(headers) != 1 || headers[0].Hash != "h1" {
		t.Fatalf("unexpected headers slice")
	}
	// Ensure returned slice is a copy
	headers[0].Hash = "modified"
	if latest, _ := ln.LatestHeader(); latest.Hash != "h1" {
		t.Fatalf("internal headers modified via returned slice")
	}
	if err := ln.AddHeader(nodes.BlockHeader{}); err == nil {
		t.Fatalf("expected error for empty hash")
	}
}
