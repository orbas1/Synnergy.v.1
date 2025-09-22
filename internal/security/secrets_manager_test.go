package security

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
)

type secretsClock struct {
	mu  sync.Mutex
	now time.Time
}

func newSecretsClock(start time.Time) *secretsClock {
	return &secretsClock{now: start}
}

func (c *secretsClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *secretsClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

func TestSecretsManagerLifecycle(t *testing.T) {
	start := time.Unix(0, 0)
	clock := newSecretsClock(start)
	key := bytes.Repeat([]byte{0x42}, chacha20poly1305.KeySize)
	sm := NewSecretsManager(WithMasterKey(key), WithSecretsClock(clock.Now))

	if err := sm.Store("api", "secret", WithTTL(time.Minute), WithMetadata(map[string]string{"owner": "cli"})); err != nil {
		t.Fatalf("store: %v", err)
	}
	if err := sm.Store("", "value"); err == nil {
		t.Fatalf("expected validation error for empty key")
	}
	val, err := sm.Retrieve("api")
	if err != nil || val != "secret" {
		t.Fatalf("unexpected value %q err %v", val, err)
	}

	meta, err := sm.Metadata("api")
	if err != nil {
		t.Fatalf("metadata: %v", err)
	}
	if meta.Version != 1 || meta.Metadata["owner"] != "cli" {
		t.Fatalf("unexpected metadata: %+v", meta)
	}

	if err := sm.Store("api", "secret2"); err != nil {
		t.Fatalf("store second version: %v", err)
	}
	meta, _ = sm.Metadata("api")
	if meta.Version != 2 {
		t.Fatalf("expected version increment got %d", meta.Version)
	}

	keys := sm.ListKeys()
	if len(keys) != 1 || keys[0] != "api" {
		t.Fatalf("unexpected keys %v", keys)
	}

	clock.Advance(2 * time.Minute)
	if _, err := sm.Retrieve("api"); err != ErrSecretExpired {
		t.Fatalf("expected expired error got %v", err)
	}
	if len(sm.ListKeys()) != 0 {
		t.Fatalf("expected expired secret purged")
	}

	if err := sm.Store("api", "fresh", WithTTL(time.Minute)); err != nil {
		t.Fatalf("store fresh: %v", err)
	}
	newKey := bytes.Repeat([]byte{0x7a}, chacha20poly1305.KeySize)
	if err := sm.RotateMasterKey(newKey); err != nil {
		t.Fatalf("rotate master: %v", err)
	}
	if val, err := sm.Retrieve("api"); err != nil || val != "fresh" {
		t.Fatalf("expected fresh secret after rotate got %q err %v", val, err)
	}

	sm.Delete("api")
	if _, err := sm.Retrieve("api"); err != ErrSecretNotFound {
		t.Fatalf("expected deletion to remove secret got %v", err)
	}
}
