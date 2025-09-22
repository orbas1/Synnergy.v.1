package security

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"synnergy/internal/telemetry"
)

type manualClock struct {
	mu  sync.Mutex
	now time.Time
}

func newManualClock(start time.Time) *manualClock {
	return &manualClock{now: start}
}

func (c *manualClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *manualClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

func TestRateLimiterBurstAndRetry(t *testing.T) {
	clock := newManualClock(time.Unix(0, 0))
	sink := telemetry.NewInMemoryRateLimitSink(8)
	limiter := NewRateLimiter(10*time.Millisecond, WithBurst(2), WithClock(clock.Now), WithSink(sink), WithComponent("rpc"))

	if ok, _ := limiter.AllowN("alice", 1); !ok {
		t.Fatalf("expected first call allowed")
	}
	clock.Advance(5 * time.Millisecond)
	if ok, _ := limiter.AllowN("alice", 1); !ok {
		t.Fatalf("expected burst to allow second call")
	}
	if ok, retry := limiter.AllowN("alice", 1); ok || retry < 5*time.Millisecond {
		t.Fatalf("expected denial with retry >= 5ms, ok=%v retry=%s", ok, retry)
	}
	clock.Advance(10 * time.Millisecond)
	if ok, _ := limiter.AllowN("alice", 1); !ok {
		t.Fatalf("expected token refilled after wait")
	}

	samples := sink.Snapshot()
	if len(samples) != 4 {
		t.Fatalf("expected 4 samples got %d", len(samples))
	}
	if samples[len(samples)-1].Identity != "alice" || !samples[len(samples)-1].Allowed {
		t.Fatalf("unexpected final sample %+v", samples[len(samples)-1])
	}
}

func TestRateLimiterSnapshotAndReconfigure(t *testing.T) {
	clock := newManualClock(time.Unix(0, 0))
	limiter := NewRateLimiter(20*time.Millisecond, WithBurst(3), WithClock(clock.Now))
	if ok, _ := limiter.AllowN("cli", 2); !ok {
		t.Fatalf("expected first call allowed")
	}
	state := limiter.Snapshot()
	if len(state) != 1 {
		t.Fatalf("expected 1 identity")
	}
	if st := state["cli"]; st.Remaining >= 3 || st.LastRefill.IsZero() {
		t.Fatalf("unexpected snapshot %#v", st)
	}

	limiter.Reconfigure(5*time.Millisecond, 5)
	clock.Advance(20 * time.Millisecond)
	if ok, _ := limiter.AllowN("cli", 5); !ok {
		t.Fatalf("expected burst after reconfigure")
	}
}

func TestRateLimiterConcurrentIdentities(t *testing.T) {
	sink := telemetry.NewInMemoryRateLimitSink(32)
	limiter := NewRateLimiter(time.Millisecond, WithBurst(1), WithSink(sink))
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			id := "id-" + strconv.Itoa(idx)
			allowed, _ := limiter.AllowN(id, 1)
			if !allowed {
				t.Errorf("expected first request for %s allowed", id)
			}
			allowed, _ = limiter.AllowN(id, 1)
			if allowed {
				t.Errorf("expected second request for %s to be throttled", id)
			}
		}(i)
	}
	wg.Wait()
	if len(limiter.Snapshot()) != 8 {
		t.Fatalf("expected 8 identities tracked")
	}
}
