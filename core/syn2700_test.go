package core

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestVestingSchedule(t *testing.T) {
	now := time.Unix(0, 0)
	entries := []VestingEntry{{ReleaseTime: now.Add(time.Hour), Amount: 50}, {ReleaseTime: now.Add(2 * time.Hour), Amount: 50}}
	schedule := NewVestingSchedule(entries)

	if claimed := schedule.Claim(now); claimed != 0 {
		t.Fatalf("expected nothing claimable at start")
	}

	after := now.Add(time.Hour + time.Minute)
	if claimed := schedule.Claim(after); claimed != 50 {
		t.Fatalf("expected 50, got %d", claimed)
	}
	if pending := schedule.Pending(after); pending != 50 {
		t.Fatalf("expected 50 pending, got %d", pending)
	}
}

func TestVestingScheduleConcurrentClaim(t *testing.T) {
	now := time.Unix(0, 0)
	entries := []VestingEntry{{ReleaseTime: now.Add(time.Hour), Amount: 50}, {ReleaseTime: now.Add(2 * time.Hour), Amount: 50}}
	schedule := NewVestingSchedule(entries)
	after := now.Add(3 * time.Hour)
	var total uint64
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddUint64(&total, schedule.Claim(after))
		}()
	}
	wg.Wait()
	if total != 100 {
		t.Fatalf("expected total 100, got %d", total)
	}
}
