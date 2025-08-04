package core

import (
	"encoding/hex"
	"testing"
)

func TestKademliaStoreFind(t *testing.T) {
	k := NewKademlia()
	k.Store("a1", []byte("value"))
	v, ok := k.FindValue("a1")
	if !ok || string(v) != "value" {
		t.Fatalf("expected to retrieve stored value")
	}
}

func TestDistance(t *testing.T) {
	id1 := hex.EncodeToString([]byte{0x00, 0x01})
	id2 := hex.EncodeToString([]byte{0x01, 0x01})
	dist := Distance(id1, id2)
	if dist.Sign() == 0 {
		t.Fatalf("expected non-zero distance")
	}
}
