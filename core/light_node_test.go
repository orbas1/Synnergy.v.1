package core

import (
	"testing"

	"synnergy/internal/nodes"
)

func TestLightNodeHeaders(t *testing.T) {
	ln := NewLightNode(nodes.Address("l1"))
	h := nodes.BlockHeader{Hash: "h1", Height: 1}
	ln.AddHeader(h)
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
}
