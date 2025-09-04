package p2p

import "time"

// KeyRotator handles periodic key rotation for a channel.
type KeyRotator struct {
	Interval time.Duration
}

// NewKeyRotator creates a KeyRotator with the given interval.
func NewKeyRotator(d time.Duration) *KeyRotator {
	return &KeyRotator{Interval: d}
}

// Rotate is a placeholder rotation implementation.
func (k *KeyRotator) Rotate() {}
