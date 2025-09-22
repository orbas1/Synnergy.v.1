package synnergy

import (
	"context"
	"crypto/rand"
	"errors"
	"testing"
)

func TestIDRegistryLifecycle(t *testing.T) {
	ctx := context.Background()
	var events []RegistryEvent
	reg := NewIDRegistry(WithRegistryEventHandler(RegistryEventHandlerFunc(func(ctx context.Context, event RegistryEvent) error {
		events = append(events, event)
		return nil
	})))

	if err := reg.RegisterContext(ctx, "addr1", "info1", map[string]string{"region": "us"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	info, metadata, ok := reg.InfoWithMetadata("addr1")
	if !ok || info != "info1" || metadata["region"] != "us" {
		t.Fatalf("unexpected info: %q %+v %v", info, metadata, ok)
	}
	if string(reg.wallets["addr1"].ciphertext) == "info1" { // ensure encrypted
		t.Fatalf("ciphertext should not equal plaintext")
	}

	if err := reg.Update(ctx, "addr1", "info2", map[string]string{"region": "eu"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	info, metadata, ok = reg.InfoWithMetadata("addr1")
	if !ok || info != "info2" || metadata["region"] != "eu" {
		t.Fatalf("unexpected updated info: %q %+v %v", info, metadata, ok)
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand: %v", err)
	}
	if err := reg.RotateKey(key); err != nil {
		t.Fatalf("rotate key: %v", err)
	}
	if info, _, ok := reg.InfoWithMetadata("addr1"); !ok || info != "info2" {
		t.Fatalf("info lost after rotate: %q %v", info, ok)
	}

	if !reg.IsRegistered("addr1") {
		t.Fatalf("expected addr1 to be registered")
	}
	if got := reg.Addresses(); len(got) != 1 || got[0] != "addr1" {
		t.Fatalf("unexpected addresses: %v", got)
	}

	if !reg.Unregister(ctx, "addr1") {
		t.Fatalf("expected unregister to return true")
	}
	if reg.IsRegistered("addr1") {
		t.Fatalf("expected addr1 to be removed")
	}

	if len(events) != 3 || events[0].Type != RegistryEventRegistered || events[1].Type != RegistryEventUpdated || events[2].Type != RegistryEventRemoved {
		t.Fatalf("unexpected events: %+v", events)
	}
}

func TestIDRegistryErrors(t *testing.T) {
	reg := NewIDRegistry()
	if err := reg.Register("", "info"); !errors.Is(err, ErrEmptyAddress) {
		t.Fatalf("expected ErrEmptyAddress, got %v", err)
	}
	if err := reg.Register("addr", "info"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Register("addr", "info"); !errors.Is(err, ErrWalletExists) {
		t.Fatalf("expected ErrWalletExists, got %v", err)
	}
	if err := reg.Update(context.Background(), "missing", "info", nil); !errors.Is(err, ErrWalletNotFound) {
		t.Fatalf("expected ErrWalletNotFound, got %v", err)
	}
	if reg.Unregister(context.Background(), "missing") {
		t.Fatalf("expected unregister missing to be false")
	}
}

func TestIDRegistryWatcherError(t *testing.T) {
	watchErr := errors.New("registry watcher failure")
	reg := NewIDRegistry(WithRegistryEventHandler(RegistryEventHandlerFunc(func(ctx context.Context, event RegistryEvent) error {
		return watchErr
	})))
	if err := reg.Register("addr", "info"); !errors.Is(err, watchErr) {
		t.Fatalf("expected watcher error, got %v", err)
	}
	if _, ok := reg.Info("addr"); !ok {
		t.Fatalf("wallet should still be stored")
	}
}
