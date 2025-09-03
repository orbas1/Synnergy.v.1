package core

import "testing"

func TestZeroTrustEngine(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := make([]byte, 32)
	if err := eng.OpenChannel("ch1", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	payload := []byte("secret")
	if _, err := eng.Send("ch1", payload); err != nil {
		t.Fatalf("send: %v", err)
	}
	pt, err := eng.Receive("ch1", 0)
	if err != nil || string(pt) != string(payload) {
		t.Fatalf("receive: %v", err)
	}
	if err := eng.CloseChannel("ch1"); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := eng.Send("ch1", payload); err == nil {
		t.Fatalf("expected error sending on closed channel")
	}
}
