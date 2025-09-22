package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestGatewayLifecycleAndRouteHandling(t *testing.T) {
	gateway := NewGateway(WithRateLimiter(NewRateLimiter(10)))
	t.Cleanup(func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := gateway.Stop(shutdownCtx); err != nil {
			t.Fatalf("gateway stop: %v", err)
		}
	})

	type payload struct {
		Miner  string `json:"miner"`
		Amount int    `json:"amount"`
	}
	type response struct {
		Accepted bool   `json:"accepted"`
		Message  string `json:"message"`
	}

	if err := gateway.RegisterRoute(Route{
		Path:        "/api/stakes",
		Methods:     []string{http.MethodPost},
		RequireAuth: true,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			var p payload
			if err := json.Unmarshal(body, &p); err != nil {
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}
			if p.Miner == "" || p.Amount <= 0 {
				http.Error(w, "invalid stake", http.StatusBadRequest)
				return
			}
			gateway.writeJSON(w, http.StatusOK, response{Accepted: true, Message: "stake registered"})
		},
	}); err != nil {
		t.Fatalf("register route: %v", err)
	}

	if err := gateway.Start(); err != nil {
		t.Fatalf("start gateway: %v", err)
	}

	addr := "http://" + gateway.Address()
	client := &http.Client{Timeout: 2 * time.Second}

	waitForGateway(t, client, addr+"/readyz")

	// readiness probes are available without authentication
	readyResp, err := client.Get(addr + "/readyz")
	if err != nil {
		t.Fatalf("readyz request failed after wait: %v", err)
	}
	defer readyResp.Body.Close()
	if readyResp.StatusCode != http.StatusOK {
		t.Fatalf("expected readyz 200, got %d", readyResp.StatusCode)
	}

	// method enforcement should reject unsupported verbs
	resp, err := client.Get(addr + "/api/stakes")
	if err != nil {
		t.Fatalf("unexpected get error: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for GET, got %d", resp.StatusCode)
	}

	// missing auth is rejected
	stakeBody, _ := json.Marshal(payload{Miner: "demo", Amount: 100})
	noAuthReq, _ := http.NewRequest(http.MethodPost, addr+"/api/stakes", bytes.NewReader(stakeBody))
	resp, err = client.Do(noAuthReq)
	if err != nil {
		t.Fatalf("no auth request failed: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing auth, got %d", resp.StatusCode)
	}

	// authenticated request is accepted
	authReq, _ := http.NewRequest(http.MethodPost, addr+"/api/stakes", bytes.NewReader(stakeBody))
	authReq.Header.Set("Authorization", "Bearer trusted-token")
	resp, err = client.Do(authReq)
	if err != nil {
		t.Fatalf("auth request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for authorized request, got %d", resp.StatusCode)
	}
	var success response
	if err := json.NewDecoder(resp.Body).Decode(&success); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !success.Accepted || success.Message == "" {
		t.Fatalf("unexpected response payload: %+v", success)
	}

	metrics := gateway.Metrics()
	if metrics.TotalRequests < 4 {
		t.Fatalf("expected at least 4 total requests, got %+v", metrics)
	}
	if metrics.Unauthorized == 0 {
		t.Fatalf("expected unauthorized counter to be incremented")
	}
	if metrics.LastRequest.IsZero() {
		t.Fatalf("expected last request timestamp to be set")
	}
}

func TestRegisterRouteValidation(t *testing.T) {
	gateway := NewGateway()
	if err := gateway.RegisterRoute(Route{}); err == nil {
		t.Fatal("expected error for missing path and handler")
	}
	if err := gateway.RegisterRoute(Route{Path: "/invalid"}); err == nil {
		t.Fatal("expected error for missing handler")
	}
	handler := func(w http.ResponseWriter, r *http.Request) {}
	if err := gateway.RegisterRoute(Route{Path: " ", Handler: handler}); err == nil {
		t.Fatal("expected error for blank path")
	}
	if err := gateway.RegisterRoute(Route{Path: "/valid", Handler: handler}); err != nil {
		t.Fatalf("unexpected error for valid route: %v", err)
	}
}

func waitForGateway(t *testing.T, client *http.Client, url string) {
	t.Helper()
	deadline := time.Now().Add(500 * time.Millisecond)
	for {
		resp, err := client.Get(url)
		if err == nil {
			_ = resp.Body.Close()
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("gateway did not become ready: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
