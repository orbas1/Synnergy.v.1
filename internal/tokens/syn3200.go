package tokens

import "sync"

// SYN3200Token converts units using a configurable ratio. It can model
// convertible assets such as wrapped tokens.
type SYN3200Token struct {
	mu    sync.RWMutex
	ratio float64
}

// NewSYN3200Token creates a new converter with the given ratio.
func NewSYN3200Token(ratio float64) *SYN3200Token {
	return &SYN3200Token{ratio: ratio}
}

// Convert returns the amount multiplied by the current ratio.
func (t *SYN3200Token) Convert(amount uint64) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return uint64(float64(amount) * t.ratio)
}

// SetRatio updates the conversion ratio.
func (t *SYN3200Token) SetRatio(r float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.ratio = r
}
