package governance

import "sync"

// ReplayProtector prevents processing of duplicate IDs.
type ReplayProtector struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// NewReplayProtector creates a ReplayProtector.
func NewReplayProtector() *ReplayProtector {
	return &ReplayProtector{seen: make(map[string]struct{})}
}

// Seen checks if the ID has been observed.
func (r *ReplayProtector) Seen(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.seen[id]; ok {
		return true
	}
	r.seen[id] = struct{}{}
	return false
}
