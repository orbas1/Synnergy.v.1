package core

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "testing"
)

func TestSubBlockCreationAndVerification(t *testing.T) {
    tx := NewTransaction("a", "b", 1, 0, 0)
    sb := NewSubBlock([]*Transaction{tx}, "val")
    if sb.PohHash != sb.Hash() {
        t.Fatalf("poh hash not set correctly")
    }
    if !sb.VerifySignature() {
        t.Fatalf("signature verification failed")
    }
}

func TestBlockHeaderHash(t *testing.T) {
    tx := NewTransaction("a", "b", 1, 0, 0)
    sb := NewSubBlock([]*Transaction{tx}, "v1")
    b := NewBlock([]*SubBlock{sb}, "prevhash")
    nonce := uint64(7)
    got := b.HeaderHash(nonce)
    h := sha256.New()
    h.Write([]byte("prevhash"))
    h.Write([]byte(sb.PohHash))
    h.Write([]byte(fmt.Sprintf("%d%d", b.Timestamp, nonce)))
    expected := hex.EncodeToString(h.Sum(nil))
    if got != expected {
        t.Fatalf("header hash mismatch")
    }
}

