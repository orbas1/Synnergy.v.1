package telemetry

import (
	"context"
	"testing"
)

func TestInMemoryRateLimitSink(t *testing.T) {
	sink := NewInMemoryRateLimitSink(2)
	sink.Reset()
	sink.RecordRateLimit("rpc", "alice", true, 3, 0)
	sink.RecordRateLimit("rpc", "bob", false, 0, 5)
	sink.RecordRateLimit("rpc", "charlie", true, 1, 0)
	samples := sink.Snapshot()
	if len(samples) != 2 {
		t.Fatalf("expected ring buffer to keep 2 entries got %d", len(samples))
	}
	if samples[0].Identity != "bob" || samples[1].Identity != "charlie" {
		t.Fatalf("unexpected order %+v", samples)
	}
	sink.Reset()
	if len(sink.Snapshot()) != 0 {
		t.Fatalf("expected reset to clear entries")
	}
}

func TestGlobalRateLimitSink(t *testing.T) {
	sink := NewInMemoryRateLimitSink(4)
	SetGlobalRateLimitSink(sink)
	ResetRateLimitHistory()
	GlobalRateLimitSink().RecordRateLimit("cli", "operator", true, 4, 0)
	GlobalRateLimitSink().RecordRateLimit("cli", "operator", false, 0, 3)
	history := RateLimitHistory()
	if len(history) != 2 {
		t.Fatalf("expected 2 samples got %d", len(history))
	}
	if !history[0].Allowed || history[1].Allowed {
		t.Fatalf("unexpected allowed flags %+v", history)
	}
}

func TestStartSpan(t *testing.T) {
	Configure(Config{Service: "synnergy", Environment: "test", Version: "1.0"})
	ctx, span := StartSpan(context.Background(), "rpc", "execute")
	if span == nil {
		t.Fatalf("expected span")
	}
	span.End()
	if trace := Tracer("rpc"); trace == nil {
		t.Fatalf("expected tracer")
	} else {
		// ensure we can start a nested span without panic
		_, nested := trace.Start(ctx, "nested")
		nested.End()
	}
}
