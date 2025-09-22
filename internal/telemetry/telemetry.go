package telemetry

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Config defines the metadata advertised when creating tracers.
type Config struct {
	Service     string
	Environment string
	Version     string
}

var (
	cfgMu sync.RWMutex
	cfg   = Config{Service: "synnergy"}

	sinkMu     sync.RWMutex
	globalSink RateLimitRecorder = NewInMemoryRateLimitSink(512)
)

// Configure updates the global telemetry metadata. The function is safe for
// repeated invocation; empty fields are ignored so callers can update a subset
// of values at runtime.
func Configure(c Config) {
	cfgMu.Lock()
	defer cfgMu.Unlock()
	if c.Service != "" {
		cfg.Service = c.Service
	}
	if c.Environment != "" {
		cfg.Environment = c.Environment
	}
	if c.Version != "" {
		cfg.Version = c.Version
	}
}

// Tracer returns an OpenTelemetry tracer namespaced by component.
func Tracer(component string) trace.Tracer {
	cfgMu.RLock()
	service := cfg.Service
	cfgMu.RUnlock()
	name := service
	if component != "" {
		name = service + "." + component
	}
	return otel.Tracer(name)
}

// StartSpan starts a span for the provided component and operation while
// applying standard attributes such as environment and version.
func StartSpan(ctx context.Context, component, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	cfgMu.RLock()
	environment := cfg.Environment
	version := cfg.Version
	cfgMu.RUnlock()
	kv := []attribute.KeyValue{
		attribute.String("environment", environment),
		attribute.String("version", version),
	}
	kv = append(kv, attrs...)
	tracer := Tracer(component)
	return tracer.Start(ctx, operation, trace.WithAttributes(kv...))
}

// RecordError annotates the span with an error and marks the status accordingly.
func RecordError(span trace.Span, err error) {
	if err == nil || span == nil {
		return
	}
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// RateLimitSample captures a rate-limit decision used by dashboards and the web
// UI to render historical behaviour.
type RateLimitSample struct {
	Component  string
	Identity   string
	Allowed    bool
	Remaining  float64
	RetryAfter time.Duration
	RecordedAt time.Time
}

// RateLimitRecorder is implemented by sinks that consume limiter decisions.
type RateLimitRecorder interface {
	RecordRateLimit(component, identity string, allowed bool, remaining float64, retryAfter time.Duration)
	Snapshot() []RateLimitSample
	Reset()
}

// GlobalRateLimitSink returns the process-wide recorder used by rate limiters.
func GlobalRateLimitSink() RateLimitRecorder {
	sinkMu.RLock()
	defer sinkMu.RUnlock()
	return globalSink
}

// SetGlobalRateLimitSink overrides the global rate limit recorder. Primarily
// used by tests or when wiring external telemetry systems.
func SetGlobalRateLimitSink(recorder RateLimitRecorder) {
	sinkMu.Lock()
	defer sinkMu.Unlock()
	if recorder == nil {
		globalSink = NewInMemoryRateLimitSink(512)
		return
	}
	globalSink = recorder
}

// RateLimitHistory returns the recorded decisions and is safe for concurrent
// use. The most recent events are ordered last in the slice.
func RateLimitHistory() []RateLimitSample {
	return GlobalRateLimitSink().Snapshot()
}

// ResetRateLimitHistory clears all recorded samples.
func ResetRateLimitHistory() {
	GlobalRateLimitSink().Reset()
}

// InMemoryRateLimitSink retains the most recent decisions in a ring buffer. The
// implementation is concurrency safe and resilient to slow consumers which makes
// it suitable for both CLI dashboards and the JS function web.
type InMemoryRateLimitSink struct {
	mu       sync.RWMutex
	capacity int
	entries  []RateLimitSample
}

// NewInMemoryRateLimitSink creates a sink storing up to capacity samples.
func NewInMemoryRateLimitSink(capacity int) *InMemoryRateLimitSink {
	if capacity <= 0 {
		capacity = 256
	}
	return &InMemoryRateLimitSink{capacity: capacity}
}

// RecordRateLimit stores a new sample trimming the oldest entries when the ring
// buffer is full.
func (s *InMemoryRateLimitSink) RecordRateLimit(component, identity string, allowed bool, remaining float64, retryAfter time.Duration) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sample := RateLimitSample{
		Component:  component,
		Identity:   identity,
		Allowed:    allowed,
		Remaining:  remaining,
		RetryAfter: retryAfter,
		RecordedAt: time.Now().UTC(),
	}
	if len(s.entries) >= s.capacity {
		copy(s.entries, s.entries[1:])
		s.entries[len(s.entries)-1] = sample
		return
	}
	s.entries = append(s.entries, sample)
}

// Snapshot returns a copy of the recorded entries for safe consumption.
func (s *InMemoryRateLimitSink) Snapshot() []RateLimitSample {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]RateLimitSample, len(s.entries))
	copy(out, s.entries)
	return out
}

// Reset removes all recorded entries.
func (s *InMemoryRateLimitSink) Reset() {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.entries = nil
	s.mu.Unlock()
}
