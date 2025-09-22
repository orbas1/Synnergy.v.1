package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGatewayStartAndMetrics(t *testing.T) {
	gateway := NewGateway()
	t.Cleanup(func() {
		_ = gateway.Stop(context.Background())
	})

	if err := gateway.RegisterRoute(Route{Path: "/secured", Methods: []string{http.MethodGet}, RequireAuth: true, Handler: func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}}); err != nil {
		t.Fatalf("register route: %v", err)
	}

	if err := gateway.Start(); err != nil {
		t.Fatalf("start gateway: %v", err)
	}

	addr := gateway.Address()
	client := &http.Client{Timeout: 2 * time.Second}

	assertStatus(t, client, "http://"+addr+"/healthz", http.StatusOK, "\"status\":\"ok\"")
	assertStatus(t, client, "http://"+addr+"/readyz", http.StatusOK, "\"ready\":true")

	// metrics require auth
	resp, err := client.Get("http://" + addr + "/metrics")
	if err != nil {
		t.Fatalf("metrics request failed: %v", err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized status, got %d", resp.StatusCode)
	}

	metricsReq, _ := http.NewRequest(http.MethodGet, "http://"+addr+"/metrics", nil)
	metricsReq.Header.Set("Authorization", "Bearer token")
	metricsResp, err := client.Do(metricsReq)
	if err != nil {
		t.Fatalf("metrics request failed: %v", err)
	}
	defer metricsResp.Body.Close()
	if metricsResp.StatusCode != http.StatusOK {
		t.Fatalf("expected metrics 200, got %d", metricsResp.StatusCode)
	}
	var snapshot GatewayMetrics
	if err := json.NewDecoder(metricsResp.Body).Decode(&snapshot); err != nil {
		t.Fatalf("decode metrics: %v", err)
	}
	if snapshot.TotalRequests < 2 {
		t.Fatalf("expected metrics to record requests, got %+v", snapshot)
	}

	// secured route
	assertStatus(t, client, "http://"+addr+"/secured", http.StatusUnauthorized, "")
	secureReq, _ := http.NewRequest(http.MethodGet, "http://"+addr+"/secured", nil)
	secureReq.Header.Set("Authorization", "Bearer token")
	secureResp, err := client.Do(secureReq)
	if err != nil {
		t.Fatalf("secured request failed: %v", err)
	}
	_ = secureResp.Body.Close()
	if secureResp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204 from secured route, got %d", secureResp.StatusCode)
	}

	metrics := gateway.Metrics()
	if metrics.TotalRequests < 4 {
		t.Fatalf("expected metrics to reflect traffic, got %+v", metrics)
	}
	if metrics.Unauthorized == 0 {
		t.Fatalf("expected unauthorized counter to increment")
	}
}

func TestGatewayRateLimiting(t *testing.T) {
	gateway := NewGateway(WithRateLimiter(NewRateLimiter(2)))
	t.Cleanup(func() {
		_ = gateway.Stop(context.Background())
	})
	if err := gateway.Start(); err != nil {
		t.Fatalf("start gateway: %v", err)
	}

	addr := gateway.Address()
	client := &http.Client{Timeout: 2 * time.Second}

	assertStatus(t, client, "http://"+addr+"/healthz", http.StatusOK, "")
	assertStatus(t, client, "http://"+addr+"/healthz", http.StatusOK, "")
	assertStatus(t, client, "http://"+addr+"/healthz", http.StatusTooManyRequests, "rate limit exceeded")

	metrics := gateway.Metrics()
	if metrics.RateLimited == 0 {
		t.Fatalf("expected rate limited counter to increment")
	}
}

func assertStatus(t *testing.T, client *http.Client, url string, expected int, contains string) {
	t.Helper()
	deadline := time.Now().Add(500 * time.Millisecond)
	var resp *http.Response
	var err error
	for {
		resp, err = client.Get(url)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("request failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}
	defer resp.Body.Close()
	if resp.StatusCode != expected {
		t.Fatalf("expected status %d from %s, got %d", expected, url, resp.StatusCode)
	}
	if contains != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read response body: %v", err)
		}
		if !strings.Contains(string(body), contains) {
			t.Fatalf("expected response body %q to contain %q", body, contains)
		}
	}
}
