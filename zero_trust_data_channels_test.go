package synnergy

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
)

// TestZeroTrustEngineBasicFlow verifies opening, sending, retrieving,
// and closing a channel works end-to-end.
func TestZeroTrustEngineBasicFlow(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := make([]byte, 32)
	if err := eng.OpenChannel("ch1", key); err != nil {
		t.Fatalf("open: %v", err)
	}

	payload := []byte("secret")
	ct, err := eng.Send("ch1", payload)
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	if bytes.Equal(ct, payload) {
		t.Fatalf("ciphertext should differ from plaintext")
	}
	pt, err := Decrypt(key, ct)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if !bytes.Equal(pt, payload) {
		t.Fatalf("expected %q got %q", payload, pt)
	}
	msgs := eng.Messages("ch1")
	if len(msgs) != 1 || !bytes.Equal(msgs[0], ct) {
		t.Fatalf("message retrieval failed")
	}
	if err := eng.CloseChannel("ch1"); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := eng.Send("ch1", payload); err == nil {
		t.Fatalf("expected error sending on closed channel")
	}
}

// TestZeroTrustEngineDuplicateOpen ensures opening the same channel twice errors.
func TestZeroTrustEngineDuplicateOpen(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := make([]byte, 32)
	if err := eng.OpenChannel("dup", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := eng.OpenChannel("dup", key); err == nil {
		t.Fatalf("expected error on duplicate open")
	}
}

// TestZeroTrustEngineSendErrors covers sending on unknown or closed channels.
func TestZeroTrustEngineSendErrors(t *testing.T) {
	eng := NewZeroTrustEngine()
	if _, err := eng.Send("missing", []byte("data")); err == nil {
		t.Fatalf("expected error sending to missing channel")
	}

	key := make([]byte, 32)
	if err := eng.OpenChannel("ch", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	if err := eng.CloseChannel("ch"); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := eng.Send("ch", []byte("data")); err == nil {
		t.Fatalf("expected error sending on closed channel")
	}
}

// TestZeroTrustEngineMessagesIsolation ensures Messages returns copies.
func TestZeroTrustEngineMessagesIsolation(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := make([]byte, 32)
	if err := eng.OpenChannel("iso", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	ct, err := eng.Send("iso", []byte("hello"))
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	msgs := eng.Messages("iso")
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
	msgs[0][0] ^= 0xFF // mutate returned slice
	msgs2 := eng.Messages("iso")
	if !bytes.Equal(msgs2[0], ct) {
		t.Fatalf("internal message mutated via external modification")
	}
}

// TestZeroTrustEngineCloseChannel verifies closing missing channels errors.
func TestZeroTrustEngineCloseChannel(t *testing.T) {
	eng := NewZeroTrustEngine()
	if err := eng.CloseChannel("nope"); err == nil {
		t.Fatalf("expected error closing unknown channel")
	}
}

// TestZeroTrustEngineConcurrentSend checks thread-safety of concurrent sends.
func TestZeroTrustEngineConcurrentSend(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := make([]byte, 32)
	if err := eng.OpenChannel("con", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	const n = 50
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			payload := []byte(fmt.Sprintf("msg-%d", i))
			if _, err := eng.Send("con", payload); err != nil {
				t.Errorf("send %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()
	if got := len(eng.Messages("con")); got != n {
		t.Fatalf("expected %d messages, got %d", n, got)
	}
}

// TestZeroTrustEngineMultipleChannels ensures isolation across channels.
func TestZeroTrustEngineMultipleChannels(t *testing.T) {
	eng := NewZeroTrustEngine()
	k1 := make([]byte, 32)
	k2 := make([]byte, 32)
	if err := eng.OpenChannel("a", k1); err != nil {
		t.Fatalf("open a: %v", err)
	}
	if err := eng.OpenChannel("b", k2); err != nil {
		t.Fatalf("open b: %v", err)
	}
	if _, err := eng.Send("a", []byte("foo")); err != nil {
		t.Fatalf("send a: %v", err)
	}
	if _, err := eng.Send("b", []byte("bar")); err != nil {
		t.Fatalf("send b: %v", err)
	}
	if len(eng.Messages("a")) != 1 || len(eng.Messages("b")) != 1 {
		t.Fatalf("messages should be isolated between channels")
	}
}

// TestZeroTrustEngineMessagesUnknown confirms unknown channels return nil slice.
func TestZeroTrustEngineMessagesUnknown(t *testing.T) {
	eng := NewZeroTrustEngine()
	if m := eng.Messages("missing"); m != nil {
		t.Fatalf("expected nil for unknown channel")
	}
}
