package core

import (
	"context"
	"sync"

	ierr "synnergy/internal/errors"
	"synnergy/internal/telemetry"
)

// ValidatorManager tracks validator stakes and slashing state.
type ValidatorManager struct {
	mu       sync.RWMutex
	stakes   map[string]uint64
	slashed  map[string]bool
	evidence map[string]string
	minStake uint64
}

// NewValidatorManager creates a manager requiring the provided minimum stake.
func NewValidatorManager(minStake uint64) *ValidatorManager {
	return &ValidatorManager{
		stakes:   make(map[string]uint64),
		slashed:  make(map[string]bool),
		evidence: make(map[string]string),
		minStake: minStake,
	}
}

// Add registers a validator with a given stake.
func (vm *ValidatorManager) Add(ctx context.Context, addr string, stake uint64) error {
	ctx, span := telemetry.Tracer("core.consensus").Start(ctx, "ValidatorManager.Add")
	defer span.End()

	vm.mu.Lock()
	defer vm.mu.Unlock()
	if stake < vm.minStake {
		return ierr.New(ierr.Invalid, "stake below minimum")
	}
	vm.stakes[addr] = stake
	delete(vm.slashed, addr)
	return nil
}

// Remove deletes a validator from the set.
func (vm *ValidatorManager) Remove(ctx context.Context, addr string) {
	ctx, span := telemetry.Tracer("core.consensus").Start(ctx, "ValidatorManager.Remove")
	defer span.End()

	vm.mu.Lock()
	defer vm.mu.Unlock()
	delete(vm.stakes, addr)
	delete(vm.slashed, addr)
}

// Slash halves the stake of the validator and marks it as slashed.
func (vm *ValidatorManager) Slash(ctx context.Context, addr string) {
	vm.SlashWithEvidence(ctx, addr, "")
}

// Reward increases the stake of the validator by the provided amount.
func (vm *ValidatorManager) Reward(ctx context.Context, addr string, amount uint64) {
	ctx, span := telemetry.Tracer("core.consensus").Start(ctx, "ValidatorManager.Reward")
	defer span.End()

	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.stakes[addr] += amount
}

// SlashWithEvidence halves the stake of the validator, marks it as slashed and
// records the provided evidence string for later auditing.
func (vm *ValidatorManager) SlashWithEvidence(ctx context.Context, addr, evidence string) {
	ctx, span := telemetry.Tracer("core.consensus").Start(ctx, "ValidatorManager.Slash")
	defer span.End()

	vm.mu.Lock()
	defer vm.mu.Unlock()
	if stake, ok := vm.stakes[addr]; ok {
		vm.stakes[addr] = stake / 2
		vm.slashed[addr] = true
		vm.evidence[addr] = evidence
	}
}

// Eligible returns a copy of the current validators eligible for selection.
func (vm *ValidatorManager) Eligible() map[string]uint64 {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	eligible := make(map[string]uint64)
	for addr, stake := range vm.stakes {
		if stake >= vm.minStake && !vm.slashed[addr] {
			eligible[addr] = stake
		}
	}
	return eligible
}

// Stake returns the current stake for a validator.
func (vm *ValidatorManager) Stake(addr string) uint64 {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.stakes[addr]
}

// Evidence returns the recorded misbehaviour proof for a validator, if any.
func (vm *ValidatorManager) Evidence(addr string) string {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	return vm.evidence[addr]
}
