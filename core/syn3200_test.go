package core

import (
	"testing"
	"time"
)

func TestBillRegistryLifecycle(t *testing.T) {
	r := NewBillRegistry()
	due := time.Now().Add(24 * time.Hour)
	b, err := r.Create("1", "issuer", "payer", 100, due, "meta")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if b.Amount != 100 {
		t.Fatalf("expected amount 100 got %d", b.Amount)
	}

	if err := r.Pay("1", "payer", 40); err != nil {
		t.Fatalf("pay: %v", err)
	}
	if b.Amount != 60 {
		t.Fatalf("expected amount 60 got %d", b.Amount)
	}

	if err := r.Adjust("1", 80); err != nil {
		t.Fatalf("adjust: %v", err)
	}
	if b.Amount != 80 {
		t.Fatalf("expected amount 80 got %d", b.Amount)
	}

	got, ok := r.Get("1")
	if !ok || got.ID != "1" {
		t.Fatalf("get failed")
	}
}

func TestBillRegistryConcurrentPay(t *testing.T) {
	r := NewBillRegistry()
	r.Create("2", "issuer", "payer", 100, time.Now(), "")

	done := make(chan struct{})
	for i := 0; i < 5; i++ {
		go func() {
			_ = r.Pay("2", "payer", 10)
			done <- struct{}{}
		}()
	}
	for i := 0; i < 5; i++ {
		<-done
	}
	b, _ := r.Get("2")
	if b.Amount != 50 {
		t.Fatalf("expected 50 got %d", b.Amount)
	}
}
