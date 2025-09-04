package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewWalletHandler(t *testing.T) {
	srv := newServer()
	req := httptest.NewRequest(http.MethodPost, "/wallet/new", nil)
	w := httptest.NewRecorder()
	srv.newWalletHandler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	var resp struct{ Address string }
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Address == "" {
		t.Fatal("address empty")
	}
}

func TestHealthHandler(t *testing.T) {
	srv := newServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	srv.healthHandler(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}
