package core

import (
	"encoding/hex"
	"testing"
)

func TestKademliaStoreFind(t *testing.T) {
	k := NewKademlia()
	if err := k.Store("a1", []byte("value")); err != nil {
		t.Fatalf("store: %v", err)
	}
	v, ok, err := k.FindValue("a1")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if !ok || string(v) != "value" {
		t.Fatalf("expected to retrieve stored value")
	}
}

func TestKademliaStoreEmptyKey(t *testing.T) {
	k := NewKademlia()
	if err := k.Store("", []byte("v")); err == nil {
		t.Fatalf("expected error for empty key")
	}
}

func TestDistance(t *testing.T) {
	id1 := hex.EncodeToString([]byte{0x00, 0x01})
	id2 := hex.EncodeToString([]byte{0x01, 0x01})
	dist, err := Distance(id1, id2)
	if err != nil {
		t.Fatalf("distance: %v", err)
	}
	if dist.Sign() == 0 {
		t.Fatalf("expected non-zero distance")
	}
}

func TestDistanceInvalidHex(t *testing.T) {
	if _, err := Distance("zz", "aa"); err == nil {
		t.Fatalf("expected error for invalid hex")
	}
}
