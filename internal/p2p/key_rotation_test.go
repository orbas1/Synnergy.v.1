package p2p

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"synnergy/internal/security"
)

func TestKeyRotatorRotate(t *testing.T) {
	km := security.NewKeyManager()
	kr := NewKeyRotator(time.Millisecond, km)
	var called int32
	kr.Register(func(keys RotationKeys) error {
		if len(keys.NoiseStatic) != 32 || len(keys.EnvelopeKey) != 32 {
			t.Fatalf("unexpected key lengths")
		}
		atomic.AddInt32(&called, 1)
		return nil
	})
	if err := kr.Rotate(); err != nil {
		t.Fatalf("rotate: %v", err)
	}
	if atomic.LoadInt32(&called) != 1 {
		t.Fatalf("expected handler called once")
	}
	if err := kr.Rotate(); err != nil {
		t.Fatalf("rotate second: %v", err)
	}
}

func TestKeyRotatorStart(t *testing.T) {
	km := security.NewKeyManager()
	kr := NewKeyRotator(time.Millisecond, km)
	var called int32
	kr.Register(func(keys RotationKeys) error {
		atomic.AddInt32(&called, 1)
		return nil
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
	defer cancel()
	kr.Start(ctx)
	if atomic.LoadInt32(&called) == 0 {
		t.Fatalf("expected background rotation")
	}
}
