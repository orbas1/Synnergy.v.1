package cli

import (
	"testing"
	"time"
)

func TestSyn3600Lifecycle(t *testing.T) {
	contract = nil
	exp := time.Now().Add(time.Hour).Format(time.RFC3339)
	if _, err := execCommand("syn3600", "create", "--underlying", "BTC", "--qty", "2", "--price", "1000", "--expiration", exp); err != nil {
		t.Fatalf("create: %v", err)
	}
	out, err := execCommand("syn3600", "status")
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if out != "false" {
		t.Fatalf("expected false, got %s", out)
	}
	out, err = execCommand("syn3600", "settle", "1100")
	if err != nil {
		t.Fatalf("settle: %v", err)
	}
	if out != "pnl 200" {
		t.Fatalf("expected pnl 200, got %s", out)
	}
}

func TestSyn3600Validation(t *testing.T) {
	contract = nil
	exp := time.Now().Add(time.Hour).Format(time.RFC3339)
	if _, err := execCommand("syn3600", "create", "--underlying", "", "--qty", "1", "--price", "1000", "--expiration", exp); err == nil {
		t.Fatal("expected error for missing underlying")
	}
	if _, err := execCommand("syn3600", "create", "--underlying", "BTC", "--qty", "0", "--price", "1000", "--expiration", exp); err == nil {
		t.Fatal("expected error for zero quantity")
	}
	if _, err := execCommand("syn3600", "create", "--underlying", "BTC", "--qty", "1", "--price", "1000", "--expiration", "bad"); err == nil {
		t.Fatal("expected error for invalid expiration")
	}
}
