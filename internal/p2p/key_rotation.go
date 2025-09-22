package p2p

import (
	"context"
	"fmt"
	"sync"
	"time"

	"synnergy/internal/security"
)

// RotationKeys captures the key material published to subscribers.
type RotationKeys struct {
	NoiseStatic    []byte
	EnvelopeKey    []byte
	SigningKey     []byte
	SigningVersion int
	RotatedAt      time.Time
}

// RotationHandler receives notifications when keys rotate.
type RotationHandler func(RotationKeys) error

// KeyRotator coordinates key rotation across transports, the VM sandbox and CLI
// tooling. Subscribers are invoked synchronously in registration order.
type KeyRotator struct {
	Interval time.Duration
	manager  *security.KeyManager

	mu        sync.RWMutex
	handlers  []RotationHandler
	inFlight  bool
	lastError error
}

// NewKeyRotator creates a KeyRotator with the given interval.
func NewKeyRotator(d time.Duration, manager *security.KeyManager) *KeyRotator {
	if manager == nil {
		manager = security.NewKeyManager()
	}
	if d <= 0 {
		d = time.Hour
	}
	return &KeyRotator{Interval: d, manager: manager}
}

// Register adds a handler to receive rotation events.
func (k *KeyRotator) Register(handler RotationHandler) {
	if handler == nil {
		return
	}
	k.mu.Lock()
	defer k.mu.Unlock()
	k.handlers = append(k.handlers, handler)
}

// Rotate executes a single rotation cycle and notifies subscribers.
func (k *KeyRotator) Rotate() error {
	k.mu.Lock()
	if k.inFlight {
		k.mu.Unlock()
		return fmt.Errorf("p2p: rotation already running")
	}
	k.inFlight = true
	k.mu.Unlock()

	defer func() {
		k.mu.Lock()
		k.inFlight = false
		k.mu.Unlock()
	}()

	noiseVersion, noiseKey, err := k.manager.GenerateSymmetricKey(security.PurposeNoiseStatic, "p2p-rotator")
	if err != nil {
		k.setLastError(err)
		return err
	}
	envelopeVersion, envelopeKey, err := k.manager.GenerateSymmetricKey(security.PurposeEnvelope, "p2p-rotator")
	if err != nil {
		k.setLastError(err)
		return err
	}
	sigPub, sigVersion, err := k.manager.GenerateSigningKey(security.PurposeStateSigning, "p2p-rotator")
	if err != nil {
		k.setLastError(err)
		return err
	}

	keys := RotationKeys{
		NoiseStatic:    noiseKey,
		EnvelopeKey:    envelopeKey,
		SigningKey:     append([]byte(nil), sigPub...),
		SigningVersion: sigVersion,
		RotatedAt:      time.Now().UTC(),
	}

	k.mu.RLock()
	handlers := append([]RotationHandler(nil), k.handlers...)
	k.mu.RUnlock()
	for _, handler := range handlers {
		if err := handler(keys); err != nil {
			k.setLastError(err)
			return err
		}
	}
	k.setLastError(nil)
	_ = noiseVersion
	_ = envelopeVersion
	return nil
}

// Start launches a ticker that triggers rotation until the context is done.
func (k *KeyRotator) Start(ctx context.Context) {
	ticker := time.NewTicker(k.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = k.Rotate()
		}
	}
}

// LastError returns the most recent failure.
func (k *KeyRotator) LastError() error {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.lastError
}

func (k *KeyRotator) setLastError(err error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.lastError = err
}
