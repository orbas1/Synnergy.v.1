package tokens

import (
	"errors"
	"sort"
	"sync"
)

var ErrWeightTooLow = errors.New("tokens: weight would underflow")

// SYN3600Token stores governance weights for addresses.
type SYN3600Token struct {
	mu      sync.RWMutex
	weights map[string]uint64
	total   uint64
}

// NewSYN3600Token creates an empty governance weight ledger.
func NewSYN3600Token() *SYN3600Token {
	return &SYN3600Token{weights: make(map[string]uint64)}
}

// SetWeight assigns voting weight to an address.
func (t *SYN3600Token) SetWeight(addr string, w uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	prev := t.weights[addr]
	t.weights[addr] = w
	if w >= prev {
		t.total += w - prev
	} else {
		t.total -= prev - w
	}
}

// IncreaseWeight increments the weight for an address.
func (t *SYN3600Token) IncreaseWeight(addr string, delta uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.weights[addr] += delta
	t.total += delta
}

// DecreaseWeight reduces the weight for an address.
func (t *SYN3600Token) DecreaseWeight(addr string, delta uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	current := t.weights[addr]
	if current < delta {
		return ErrWeightTooLow
	}
	t.weights[addr] = current - delta
	t.total -= delta
	return nil
}

// Weight returns the voting weight of an address.
func (t *SYN3600Token) Weight(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.weights[addr]
}

// TotalWeight returns the sum of all weights.
func (t *SYN3600Token) TotalWeight() uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.total
}

// Snapshot returns a copy of the weight map.
func (t *SYN3600Token) Snapshot() map[string]uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	cp := make(map[string]uint64, len(t.weights))
	for addr, w := range t.weights {
		cp[addr] = w
	}
	return cp
}

// TopHolders returns the top n addresses by weight.
func (t *SYN3600Token) TopHolders(n int) []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	type kv struct {
		addr   string
		weight uint64
	}
	arr := make([]kv, 0, len(t.weights))
	for addr, w := range t.weights {
		arr = append(arr, kv{addr: addr, weight: w})
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].weight == arr[j].weight {
			return arr[i].addr < arr[j].addr
		}
		return arr[i].weight > arr[j].weight
	})
	if n <= 0 || n > len(arr) {
		n = len(arr)
	}
	res := make([]string, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, arr[i].addr)
	}
	return res
}
