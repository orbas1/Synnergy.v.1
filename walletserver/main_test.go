package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunRegistersHandlers(t *testing.T) {
	oldListen := httpListenAndServe
	defer func() { httpListenAndServe = oldListen }()

	var capturedAddr string
	var handler http.Handler
	httpListenAndServe = func(addr string, h http.Handler) error {
		capturedAddr = addr
		handler = h
		return nil
	}

	if err := run(":9999", newServer()); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if capturedAddr != ":9999" {
		t.Fatalf("expected addr :9999, got %s", capturedAddr)
	}
	if handler == nil {
		t.Fatal("expected handler to be registered")
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("health status: %d", rr.Code)
	}
	var health map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &health); err != nil {
		t.Fatalf("decode health: %v", err)
	}
	if health["status"] != "ok" {
		t.Fatalf("unexpected health response: %v", health)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/wallet/new", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected method not allowed, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/wallet/new", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("wallet new status: %d", rr.Code)
	}
	var wallet map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &wallet); err != nil {
		t.Fatalf("decode wallet: %v", err)
	}
	if wallet["address"] == "" {
		t.Fatalf("expected non-empty address: %v", wallet)
	}
}
