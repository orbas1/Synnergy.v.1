package main

import (
	"bytes"
	"strings"
	"testing"

	"synnergy/internal/p2p"
)

func TestRunAddAndListPeers(t *testing.T) {
	mgr := p2p.NewManager()
	if code := runWithManager(mgr, []string{"add-peer", "-id", "p1", "-addr", "127.0.0.1:1"}, &bytes.Buffer{}); code != 0 {
		t.Fatalf("add-peer exit code %d", code)
	}
	var buf bytes.Buffer
	if code := runWithManager(mgr, []string{"list-peers"}, &buf); code != 0 {
		t.Fatalf("list-peers exit code %d", code)
	}
	if !strings.Contains(buf.String(), "p1") {
		t.Fatalf("expected peer in output, got %q", buf.String())
	}
}
