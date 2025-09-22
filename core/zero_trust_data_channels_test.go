package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
	"time"
)

func TestZeroTrustEngineSendAndReceive(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := bytesKey(32)
	if err := eng.OpenChannel("ch1", key, WithOwner("ops")); err != nil {
		t.Fatalf("open: %v", err)
	}
	info, err := eng.ChannelInfo("ch1")
	if err != nil || info.Owner != "ops" {
		t.Fatalf("channel info mismatch: %+v %v", info, err)
	}
	if _, err := eng.SendAs("ch1", "ops", []byte("secret")); err != nil {
		t.Fatalf("send: %v", err)
	}
	pt, err := eng.Receive("ch1", 0)
	if err != nil || string(pt) != "secret" {
		t.Fatalf("receive: %v", err)
	}
	events := eng.Events()
	if len(events) < 2 {
		t.Fatalf("expected events, got %+v", events)
	}
}

func TestZeroTrustEngineAuthorizeAndRotate(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := bytesKey(32)
	if err := eng.OpenChannel("ch1", key); err != nil {
		t.Fatalf("open: %v", err)
	}
	pub, _, _ := ed25519.GenerateKey(nil)
	if err := eng.AuthorizePeer("ch1", "ally", hex.EncodeToString(pub)); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	if _, err := eng.SendAs("ch1", "ally", []byte("joint")); err != nil {
		t.Fatalf("send as ally: %v", err)
	}
	if err := eng.RotateKey("ch1", bytesKey(32)); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if err := eng.CloseChannel("ch1"); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := eng.Send("ch1", []byte("blocked")); err == nil {
		t.Fatalf("expected send to fail on closed channel")
	}
}

func TestZeroTrustEngineEventsSubscription(t *testing.T) {
	eng := NewZeroTrustEngine()
	key := bytesKey(32)
	if err := eng.OpenChannel("ch", key, WithOwner("ops")); err != nil {
		t.Fatalf("open: %v", err)
	}
	ch, cancel := eng.SubscribeEvents(4)
	defer cancel()
	select {
	case ev := <-ch:
		if ev.Type != ChannelEventOpened {
			t.Fatalf("expected open event, got %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected initial open event")
	}

	if _, err := eng.Send("ch", []byte("msg")); err != nil {
		t.Fatalf("send: %v", err)
	}
	select {
	case ev := <-ch:
		if ev.Type != ChannelEventMessage {
			t.Fatalf("unexpected event %+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatalf("expected message event")
	}
}

func bytesKey(size int) []byte {
	key := make([]byte, size)
	for i := range key {
		key[i] = byte(i)
	}
	return key
}
