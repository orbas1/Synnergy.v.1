package governance

import (
	"errors"
	"sort"
	"sync"
	"time"
)

// ErrReplayDetected indicates that a duplicate identifier was observed within the protection window.
var ErrReplayDetected = errors.New("replay detected")

// ReplayOption customises the replay protector.
type ReplayOption func(*ReplayProtector)

// ReplayProtector prevents processing of duplicate IDs.
type ReplayProtector struct {
	mu          sync.Mutex
	seen        map[string]time.Time
	window      time.Duration
	maxEntries  int
	onDuplicate func(string)
	onEvict     func(string)
}

// NewReplayProtector creates a ReplayProtector.
func NewReplayProtector(opts ...ReplayOption) *ReplayProtector {
	rp := &ReplayProtector{
		seen:       make(map[string]time.Time),
		window:     10 * time.Minute,
		maxEntries: 100_000,
	}
	for _, opt := range opts {
		opt(rp)
	}
	return rp
}

// WithWindow overrides the replay detection window.
func WithWindow(window time.Duration) ReplayOption {
	return func(r *ReplayProtector) {
		if window > 0 {
			r.window = window
		}
	}
}

// WithMaxEntries caps the number of stored identifiers.
func WithMaxEntries(max int) ReplayOption {
	return func(r *ReplayProtector) {
		if max > 0 {
			r.maxEntries = max
		}
	}
}

// WithDuplicateCallback registers a callback for duplicates.
func WithDuplicateCallback(cb func(string)) ReplayOption {
	return func(r *ReplayProtector) {
		r.onDuplicate = cb
	}
}

// WithEvictionCallback registers a callback when identifiers are evicted.
func WithEvictionCallback(cb func(string)) ReplayOption {
	return func(r *ReplayProtector) {
		r.onEvict = cb
	}
}

// Seen checks if the ID has been observed inside the configured window. It returns true if the identifier is a replay.
func (r *ReplayProtector) Seen(id string) bool {
	_, duplicate := r.check(id, time.Now().UTC())
	return duplicate
}

// Check evaluates the identifier at a given timestamp and returns ErrReplayDetected when the id is a duplicate.
func (r *ReplayProtector) Check(id string, ts time.Time) error {
	_, duplicate := r.check(id, ts)
	if duplicate {
		return ErrReplayDetected
	}
	return nil
}

func (r *ReplayProtector) check(id string, ts time.Time) (time.Time, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cleanup(ts)
	if firstSeen, ok := r.seen[id]; ok {
		if r.onDuplicate != nil {
			r.onDuplicate(id)
		}
		return firstSeen, true
	}

	if len(r.seen) >= r.maxEntries {
		r.evictOldest()
	}

	r.seen[id] = ts
	return ts, false
}

func (r *ReplayProtector) cleanup(now time.Time) {
	limit := now.Add(-r.window)
	for id, firstSeen := range r.seen {
		if firstSeen.Before(limit) {
			delete(r.seen, id)
			if r.onEvict != nil {
				r.onEvict(id)
			}
		}
	}
}

func (r *ReplayProtector) evictOldest() {
	if len(r.seen) == 0 {
		return
	}
	type entry struct {
		id string
		ts time.Time
	}
	entries := make([]entry, 0, len(r.seen))
	for id, ts := range r.seen {
		entries = append(entries, entry{id: id, ts: ts})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].ts.Before(entries[j].ts) })
	oldest := entries[0]
	delete(r.seen, oldest.id)
	if r.onEvict != nil {
		r.onEvict(oldest.id)
	}
}

// Size returns the number of tracked identifiers (useful for tests and metrics).
func (r *ReplayProtector) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.seen)
}
