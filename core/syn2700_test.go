package core

import (
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
