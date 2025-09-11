package core

import (
	"errors"
	"testing"
)

func TestIdentityService(t *testing.T) {
	svc := NewIdentityService()
	if err := svc.Register("addr1", "Alice", "2000-01-01", "US"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.Register("", "Alice", "2000-01-01", "US"); !errors.Is(err, ErrAddressRequired) {
		t.Fatalf("expected ErrAddressRequired got %v", err)
	}
	// duplicate registration
	if err := svc.Register("addr1", "Alice", "2000-01-01", "US"); !errors.Is(err, ErrIdentityExists) {
		t.Fatalf("expected ErrIdentityExists got %v", err)
	}
	if err := svc.Verify("addr1", "passport"); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := svc.Verify("addr1", ""); !errors.Is(err, ErrMethodRequired) {
		t.Fatalf("expected ErrMethodRequired got %v", err)
	}
	if err := svc.Verify("unknown", "passport"); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("expected ErrIdentityNotFound got %v", err)
	}
	if err := svc.Verify("", "passport"); !errors.Is(err, ErrAddressRequired) {
		t.Fatalf("expected ErrAddressRequired got %v", err)
	}
	info, ok := svc.Info("addr1")
	if !ok || info.Name != "Alice" {
		t.Fatalf("unexpected info: %v", info)
	}
	logs := svc.Logs("addr1")
	if len(logs) != 1 || logs[0].Method != "passport" {
		t.Fatalf("unexpected logs: %v", logs)
	}
}
