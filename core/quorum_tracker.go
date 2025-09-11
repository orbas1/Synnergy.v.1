package core

import "sync"

// QuorumTracker tracks validator participation to ensure quorum is met.
type QuorumTracker struct {
	mu       sync.RWMutex
	required int
	present  map[string]struct{}
}

// NewQuorumTracker creates a tracker requiring the specified number of
// participants.
func NewQuorumTracker(required int) *QuorumTracker {
	if required < 0 {
		required = 0
	}
	return &QuorumTracker{required: required, present: make(map[string]struct{})}
}

// Join marks a validator as present.
func (qt *QuorumTracker) Join(id string) {
	qt.mu.Lock()
	defer qt.mu.Unlock()
	qt.present[id] = struct{}{}
}

// Leave removes a validator from the present set.
func (qt *QuorumTracker) Leave(id string) {
	qt.mu.Lock()
	defer qt.mu.Unlock()
	delete(qt.present, id)
}

// Count returns the current number of active validators.
func (qt *QuorumTracker) Count() int {
	qt.mu.RLock()
	defer qt.mu.RUnlock()
	return len(qt.present)
}

// Reached reports whether the quorum threshold has been met.
func (qt *QuorumTracker) Reached() bool {
	qt.mu.RLock()
	defer qt.mu.RUnlock()
	return len(qt.present) >= qt.required
}

// Reset clears the tracker of all present validators.
func (qt *QuorumTracker) Reset() {
	qt.mu.Lock()
	qt.present = make(map[string]struct{})
	qt.mu.Unlock()
}

// SetRequired updates the quorum threshold.
func (qt *QuorumTracker) SetRequired(required int) {
	if required < 0 {
		required = 0
	}
	qt.mu.Lock()
	qt.required = required
	qt.mu.Unlock()
}

// Required returns the current quorum threshold.
func (qt *QuorumTracker) Required() int {
	qt.mu.RLock()
	defer qt.mu.RUnlock()
	return qt.required
}
