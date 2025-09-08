package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSubBlockCreationAndVerification(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "val")
	if err := sb.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
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
	b.Nonce = nonce
	b.Hash = got
	if err := b.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}

func TestSubBlockValidateRejectsDuplicateTransactions(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx, tx}, "val")
	if err := sb.Validate(); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected duplicate transaction error, got %v", err)
	}
}

func TestBlockValidateRejectsFutureTimestamp(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "v1")
	b := NewBlock([]*SubBlock{sb}, "prevhash")
	b.Timestamp = time.Now().Add(6 * time.Minute).Unix()
	b.Nonce = 1
	b.Hash = b.HeaderHash(b.Nonce)
	if err := b.Validate(); err == nil || !strings.Contains(err.Error(), "timestamp in future") {
		t.Fatalf("expected future timestamp error, got %v", err)
	}
}

func TestBlockValidateRejectsSubBlockAfterBlockTimestamp(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "v1")
	b := NewBlock([]*SubBlock{sb}, "prevhash")
	b.Timestamp = sb.Timestamp - 10
	b.Nonce = 1
	b.Hash = b.HeaderHash(b.Nonce)
	if err := b.Validate(); err == nil || !strings.Contains(err.Error(), "sub-block timestamp") {
		t.Fatalf("expected sub-block timestamp error, got %v", err)
	}
}

func TestBlockValidateRequiresHash(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "v1")
	b := NewBlock([]*SubBlock{sb}, "prevhash")
	b.Nonce = 1
	if err := b.Validate(); err == nil || !strings.Contains(err.Error(), "hash required") {
		t.Fatalf("expected hash required error, got %v", err)
	}
}

func TestBlockValidateDetectsHashMismatch(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	sb := NewSubBlock([]*Transaction{tx}, "v1")
	b := NewBlock([]*SubBlock{sb}, "prevhash")
	b.Nonce = 1
	b.Hash = "bad"
	if err := b.Validate(); err == nil || !strings.Contains(err.Error(), "hash mismatch") {
		t.Fatalf("expected hash mismatch error, got %v", err)
	}
}
