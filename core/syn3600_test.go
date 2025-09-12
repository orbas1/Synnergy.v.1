package core

import (
	"sync"
	"testing"
	"time"
)

func TestFuturesContractConcurrentSettlement(t *testing.T) {
	exp := time.Unix(100, 0)
	f := NewFuturesContract("BTC", 2, 1000, exp)

	var wg sync.WaitGroup
	successes := make(chan int64, 10)
	errs := make(chan error, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pnl, err := f.Settle(1100)
			if err != nil {
				errs <- err
				return
			}
			successes <- pnl
		}()
	}
	wg.Wait()
	close(successes)
	close(errs)

	var okCount int
	var pnl int64
	for p := range successes {
		okCount++
		pnl = p
	}
	if okCount != 1 {
		t.Fatalf("expected exactly one successful settlement, got %d", okCount)
	}
	if pnl != 200 {
		t.Fatalf("expected pnl 200, got %d", pnl)
	}
	for err := range errs {
		if err.Error() != "contract already settled" {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if !f.IsExpired(time.Unix(200, 0)) {
		t.Fatalf("expected expired")
	}
}
