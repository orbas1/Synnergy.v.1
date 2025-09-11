package core

import (
	"errors"
	"testing"
)

func TestIDRegistry(t *testing.T) {
	reg := NewIDRegistry()
	if err := reg.Register("addr1", "metadata"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if _, ok := reg.Info("addr1"); !ok {
		t.Fatalf("expected wallet registered")
	}
	if err := reg.Register("addr1", "other"); !errors.Is(err, ErrWalletExists) {
		t.Fatalf("expected ErrWalletExists got %v", err)
	}
	if err := reg.Register("", "info"); !errors.Is(err, ErrEmptyAddress) {
		t.Fatalf("expected ErrEmptyAddress got %v", err)
	}
	if !reg.Unregister("addr1") {
		t.Fatalf("expected unregister to succeed")
	}
	if _, ok := reg.Info("addr1"); ok {
		t.Fatalf("wallet should be removed")
	}
}
