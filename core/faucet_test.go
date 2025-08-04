package core

import (
	"testing"
	"time"
)

func TestFaucet(t *testing.T) {
	f := NewFaucet(100, 10, time.Second)
	amt, err := f.Request("addr")
	if err != nil || amt != 10 {
		t.Fatalf("first request failed: %v %d", err, amt)
	}
	if _, err := f.Request("addr"); err == nil {
		t.Fatalf("expected cooldown error")
	}
	time.Sleep(time.Second)
	if _, err := f.Request("addr"); err != nil {
		t.Fatalf("second request after cooldown failed: %v", err)
	}
	f.UpdateConfig(5, time.Millisecond)
	if f.Balance() != 80 {
		t.Fatalf("unexpected balance: %d", f.Balance())
	}
}
