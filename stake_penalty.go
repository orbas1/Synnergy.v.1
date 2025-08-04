package synnergy

import (
	"sync"
	"time"
)

// PenaltyRecord captures information about a stake penalty applied to a
// validator.
type PenaltyRecord struct {
	Points int
	Reason string
	Time   time.Time
}

// StakePenaltyManager tracks validator stake levels and penalty points.  It
// provides helper methods used by the CLI to adjust stake, record penalties and
// fetch information for a validator.
type StakePenaltyManager struct {
	mu        sync.RWMutex
	stakes    map[string]uint64
	penalties map[string]int
	history   map[string][]PenaltyRecord
}

// NewStakePenaltyManager creates an empty manager instance.
func NewStakePenaltyManager() *StakePenaltyManager {
	return &StakePenaltyManager{
		stakes:    make(map[string]uint64),
		penalties: make(map[string]int),
		history:   make(map[string][]PenaltyRecord),
	}
}

// AdjustStake increases or decreases stake for a validator by delta.  Delta may
// be negative to reduce stake.
func (m *StakePenaltyManager) AdjustStake(addr string, delta int64) {
	m.mu.Lock()
	current := int64(m.stakes[addr]) + delta
	if current < 0 {
		current = 0
	}
	m.stakes[addr] = uint64(current)
	m.mu.Unlock()
}

// Penalize records penalty points against the validator with an optional reason.
func (m *StakePenaltyManager) Penalize(addr string, points int, reason string) {
	m.mu.Lock()
	m.penalties[addr] += points
	rec := PenaltyRecord{Points: points, Reason: reason, Time: time.Now().UTC()}
	m.history[addr] = append(m.history[addr], rec)
	m.mu.Unlock()
}

// Info returns the current stake amount, total penalty points and history for
// the given validator address.
func (m *StakePenaltyManager) Info(addr string) (stake uint64, penalty int, history []PenaltyRecord) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stakes[addr], m.penalties[addr], append([]PenaltyRecord(nil), m.history[addr]...)
}
