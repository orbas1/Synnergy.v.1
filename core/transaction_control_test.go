package core

import (
	"testing"
	"time"
)

func TestScheduleAndCancel(t *testing.T) {
	tx := NewTransaction("a", "b", 1, 0, 0)
	exec := time.Now().Add(time.Hour)
	st := ScheduleTransaction(tx, exec)
	if !CancelTransaction(st) {
		t.Fatal("expected cancel to succeed")
	}
	if !st.Canceled {
		t.Fatal("scheduled transaction not marked canceled")
	}
}

func TestReverseTransaction(t *testing.T) {
	l := NewLedger()
	l.Credit("a", 20)
	tx := NewTransaction("a", "b", 10, 2, 0)
	if err := l.ApplyTransaction(tx); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if err := ReverseTransaction(l, tx); err != nil {
		t.Fatalf("reverse: %v", err)
	}
	if l.GetBalance("a") != 20 || l.GetBalance("b") != 0 {
		t.Fatalf("unexpected balances: a=%d b=%d", l.GetBalance("a"), l.GetBalance("b"))
	}
}

func TestConvertToPrivate(t *testing.T) {
	key := []byte("example key 1234")
	tx := NewTransaction("a", "b", 5, 1, 0)
	pt, err := ConvertToPrivate(tx, key)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	dec, err := pt.Decrypt(key)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if dec.ID != tx.ID {
		t.Fatalf("expected %s got %s", tx.ID, dec.ID)
	}
}

func TestReceiptStore(t *testing.T) {
	rs := NewReceiptStore()
	r := GenerateReceipt("tx1", "ok", "details")
	rs.Store(r)
	if _, ok := rs.Get("tx1"); !ok {
		t.Fatal("receipt not stored")
	}
	res := rs.Search("ok")
	if len(res) != 1 {
		t.Fatalf("expected 1 result got %d", len(res))
	}
}
