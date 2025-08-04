package core

import "time"

// VestingEntry defines a point in time when a portion becomes available.
type VestingEntry struct {
	ReleaseTime time.Time
	Amount      uint64
	Claimed     bool
}

// VestingSchedule represents a series of vesting entries.
type VestingSchedule struct {
	Entries []VestingEntry
}

// NewVestingSchedule creates a schedule with the given entries.
func NewVestingSchedule(entries []VestingEntry) *VestingSchedule {
	return &VestingSchedule{Entries: entries}
}

// Claim releases all matured, unclaimed amounts up to now.
func (v *VestingSchedule) Claim(now time.Time) uint64 {
	var total uint64
	for i := range v.Entries {
		e := &v.Entries[i]
		if !e.Claimed && !now.Before(e.ReleaseTime) {
			e.Claimed = true
			total += e.Amount
		}
	}
	return total
}

// Pending returns the amount still locked after the provided time.
func (v *VestingSchedule) Pending(now time.Time) uint64 {
	var total uint64
	for _, e := range v.Entries {
		if !e.Claimed && now.Before(e.ReleaseTime) {
			total += e.Amount
		}
	}
	return total
}
