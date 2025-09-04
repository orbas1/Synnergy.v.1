package p2p

import (
	"testing"
	"time"
)

func TestKeyRotatorRotate(t *testing.T) {
	kr := NewKeyRotator(time.Second)
	if kr.Interval != time.Second {
		t.Fatalf("unexpected interval: %v", kr.Interval)
	}
	kr.Rotate() // should not panic
}
