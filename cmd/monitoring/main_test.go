package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	synn "synnergy"
	"synnergy/internal/nodes/extra/watchtower"
)

// TestMetricsEndpoint ensures the monitoring server exposes watchtower metrics
// over HTTP.
func TestMetricsEndpoint(t *testing.T) {
	wt := synn.NewWatchtowerNode("t", nil)
	if err := wt.Start(context.Background()); err != nil {
		t.Fatalf("start watchtower: %v", err)
	}
	defer wt.Stop()

	srv := httptest.NewServer(newHandler(wt))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/metrics")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	var m watchtower.Metrics
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if m.PeerCount != 0 {
		t.Fatalf("unexpected metrics: %+v", m)
	}
}
