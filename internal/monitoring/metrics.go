package monitoring

import "sync"

// Metrics provides a simple counter registry.
type Metrics struct {
	mu       sync.Mutex
	counters map[string]int
}

// NewMetrics creates a Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{counters: make(map[string]int)}
}

// Inc increments the named counter.
func (m *Metrics) Inc(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name]++
}

// Get retrieves the current value for the named counter.
func (m *Metrics) Get(name string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.counters[name]
}
