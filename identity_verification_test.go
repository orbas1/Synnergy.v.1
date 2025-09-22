package synnergy

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"
)

func TestIdentityServiceLifecycle(t *testing.T) {
	ctx := context.Background()
	var events []IdentityEvent
	svc := NewIdentityService(WithIdentityEventHandler(IdentityEventHandlerFunc(func(ctx context.Context, event IdentityEvent) error {
		events = append(events, event)
		return nil
	})))

	if err := svc.Register("addr1", "Alice", "2000-01-01", "US"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if info, ok := svc.Info("addr1"); !ok || info.Name != "Alice" {
		t.Fatalf("unexpected info: %+v %v", info, ok)
	}
	if status, ok := svc.Status("addr1"); !ok || status != IdentityStatusPending {
		t.Fatalf("expected pending status, got %v %v", status, ok)
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	ts := time.Now().UTC()
	req := VerificationRequest{
		Method:          "passport",
		EvidenceHash:    []byte("passport#123"),
		SignerPublicKey: pub,
		Metadata:        map[string]string{"document": "passport"},
		Timestamp:       ts,
	}
	req.Signature = ed25519.Sign(priv, verificationDigest("addr1", req))
	if err := svc.VerifyContext(ctx, "addr1", req); err != nil {
		t.Fatalf("verify: %v", err)
	}

	logs := svc.Logs("addr1")
	if len(logs) != 1 || logs[0].Method != "passport" || !logs[0].Successful {
		t.Fatalf("unexpected logs: %+v", logs)
	}
	if status, ok := svc.Status("addr1"); !ok || status != IdentityStatusVerified {
		t.Fatalf("expected verified status, got %v %v", status, ok)
	}

	if err := svc.Revoke(ctx, "addr1", "kyc_revocation", map[string]string{"reason": "kyc"}); err != nil {
		t.Fatalf("revoke: %v", err)
	}
	if status, ok := svc.Status("addr1"); !ok || status != IdentityStatusRevoked {
		t.Fatalf("expected revoked status, got %v %v", status, ok)
	}
	logs = svc.Logs("addr1")
	if len(logs) != 2 || logs[1].Method != "revoke" || logs[1].Successful {
		t.Fatalf("expected revoke log entry, got %+v", logs)
	}

	if got := svc.Addresses(); len(got) != 1 || got[0] != "addr1" {
		t.Fatalf("unexpected addresses: %v", got)
	}
	if record, ok := svc.Record("addr1"); !ok || record.Revision != 3 {
		t.Fatalf("unexpected record: %+v %v", record, ok)
	}

	if len(events) != 3 { // registered, verified, revoked
		t.Fatalf("expected 3 events, got %d", len(events))
	}
	if events[0].Type != IdentityEventRegistered || events[1].Type != IdentityEventVerified || events[2].Type != IdentityEventRevoked {
		t.Fatalf("unexpected event order: %+v", events)
	}
}

func TestIdentityServiceErrorsAndLogRetention(t *testing.T) {
	svc := NewIdentityService(WithIdentityLogRetention(2))
	if err := svc.Register("addr", "Bob", "1990-02-02", "CA"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.Register("addr", "Bob", "1990-02-02", "CA"); !errors.Is(err, ErrIdentityExists) {
		t.Fatalf("expected ErrIdentityExists, got %v", err)
	}
	if err := svc.Verify("", "passport"); !errors.Is(err, ErrAddressRequired) {
		t.Fatalf("expected ErrAddressRequired, got %v", err)
	}
	if err := svc.Verify("addr", ""); !errors.Is(err, ErrMethodRequired) {
		t.Fatalf("expected ErrMethodRequired, got %v", err)
	}

	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	req := VerificationRequest{Method: "passport", SignerPublicKey: pub, Signature: []byte("bad"), Timestamp: time.Now().UTC()}
	if err := svc.VerifyContext(context.Background(), "addr", req); !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}

	good := VerificationRequest{Method: "passport"}
	if err := svc.VerifyContext(context.Background(), "addr", good); err != nil {
		t.Fatalf("verify: %v", err)
	}
	if err := svc.VerifyContext(context.Background(), "addr", VerificationRequest{Method: "passport2"}); err != nil {
		t.Fatalf("verify second: %v", err)
	}
	if err := svc.VerifyContext(context.Background(), "addr", VerificationRequest{Method: "passport3"}); err != nil {
		t.Fatalf("verify third: %v", err)
	}

	logs := svc.Logs("addr")
	if len(logs) != 2 || logs[0].Method != "passport2" || logs[1].Method != "passport3" {
		t.Fatalf("expected last two logs retained, got %+v", logs)
	}

	if err := svc.Revoke(context.Background(), "missing", "", nil); !errors.Is(err, ErrIdentityNotFound) {
		t.Fatalf("expected ErrIdentityNotFound, got %v", err)
	}
}

func TestIdentityServiceWatcherErrorPropagates(t *testing.T) {
	watchErr := errors.New("listener failed")
	svc := NewIdentityService(WithIdentityEventHandler(IdentityEventHandlerFunc(func(ctx context.Context, event IdentityEvent) error {
		return watchErr
	})))

	err := svc.Register("addr", "Alice", "2000-03-03", "UK")
	if !errors.Is(err, watchErr) {
		t.Fatalf("expected watcher error, got %v", err)
	}
	if _, ok := svc.Info("addr"); !ok {
		t.Fatalf("identity should still be registered")
	}
}

func TestIdentityRecordCopyIsolation(t *testing.T) {
	svc := NewIdentityService()
	if err := svc.Register("addr", "Bob", "1990-02-02", "CA"); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := svc.Verify("addr", "passport"); err != nil {
		t.Fatalf("verify: %v", err)
	}
	record, ok := svc.Record("addr")
	if !ok {
		t.Fatalf("expected record")
	}
	if record.Metadata == nil {
		record.Metadata = map[string]string{}
	}
	record.Metadata["tamper"] = "true"
	if fresh, _ := svc.Record("addr"); fresh.Metadata != nil && fresh.Metadata["tamper"] == "true" {
		t.Fatalf("metadata should be isolated, got %+v", fresh.Metadata)
	}
}
