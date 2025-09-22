package monitoring

import (
        "sync"
        "sync/atomic"
        "time"
)

const defaultHistoryLimit = 256

// Span represents a single traced operation captured by the in-process tracer.
type Span struct {
        ID         uint64
        Name       string
        Start      time.Time
        End        time.Time
        Duration   time.Duration
        Attributes map[string]string
        Error      error
        Status     string
}

// SpanExporter receives completed spans. Exporters should return quickly and
// must not panic – the tracer swallows exporter errors by design.
type SpanExporter interface {
        ExportSpan(Span) error
}

// SpanFinishOption mutates a span when it is completed.
type SpanFinishOption func(*Span)

// TracerOption configures optional behaviour when constructing a tracer.
type TracerOption func(*Tracer)

// TracerMetrics exposes aggregated statistics for emitted spans.
type TracerMetrics struct {
        SpansStarted   uint64
        SpansFailed    uint64
        AverageLatency time.Duration
}

// Tracer implements a lightweight in-memory tracer with optional exporters.
type Tracer struct {
        mu           sync.RWMutex
        history      []Span
        historyLimit int
        exporters    []SpanExporter
        clock        func() time.Time

        seq atomic.Uint64
        metrics struct {
                started    atomic.Uint64
                failed     atomic.Uint64
                durationNS atomic.Int64
        }
}

// NewTracer constructs a tracer using the provided options.
func NewTracer(opts ...TracerOption) *Tracer {
        tracer := &Tracer{
                historyLimit: defaultHistoryLimit,
                clock:        time.Now,
        }
        for _, opt := range opts {
                if opt != nil {
                        opt(tracer)
                }
        }
        return tracer
}

// WithHistoryLimit overrides the number of spans retained in memory. A limit of
// zero disables pruning and keeps all spans.
func WithHistoryLimit(limit int) TracerOption {
        return func(t *Tracer) {
                        if limit < 0 {
                                limit = 0
                        }
                        t.historyLimit = limit
        }
}

// WithExporter registers an exporter during tracer construction.
func WithExporter(exporter SpanExporter) TracerOption {
        return func(t *Tracer) {
                        t.RegisterExporter(exporter)
        }
}

// WithClock overrides the clock used to timestamp spans. Primarily used in
// tests to provide deterministic timings.
func WithClock(clock func() time.Time) TracerOption {
        return func(t *Tracer) {
                        if clock != nil {
                                t.clock = clock
                        }
        }
}

// StartSpan begins a span and returns a closure to finish it. The returned
// function is idempotent – multiple invocations record the span once.
func (t *Tracer) StartSpan(name string) func(...SpanFinishOption) {
        if name == "" {
                name = "unnamed"
        }
        start := t.now()
        id := t.seq.Add(1)
        var once sync.Once
        return func(opts ...SpanFinishOption) {
                once.Do(func() {
                        span := Span{ID: id, Name: name, Start: start}
                        span.End = t.now()
                        span.Duration = span.End.Sub(span.Start)
                        for _, opt := range opts {
                                if opt != nil {
                                        opt(&span)
                                }
                        }
                        t.record(span)
                })
        }
}

// RegisterExporter attaches an exporter to the tracer at runtime.
func (t *Tracer) RegisterExporter(exporter SpanExporter) {
        if exporter == nil {
                return
        }
        t.mu.Lock()
        t.exporters = append(t.exporters, exporter)
        t.mu.Unlock()
}

// Spans returns a copy of the recorded spans.
func (t *Tracer) Spans() []Span {
        t.mu.RLock()
        defer t.mu.RUnlock()
        out := make([]Span, len(t.history))
        copy(out, t.history)
        return out
}

// Reset clears the history and metrics.
func (t *Tracer) Reset() {
        t.mu.Lock()
        t.history = nil
        t.mu.Unlock()
        t.metrics.started.Store(0)
        t.metrics.failed.Store(0)
        t.metrics.durationNS.Store(0)
        t.seq.Store(0)
}

// Metrics returns aggregated span statistics.
func (t *Tracer) Metrics() TracerMetrics {
        started := t.metrics.started.Load()
        failed := t.metrics.failed.Load()
        duration := t.metrics.durationNS.Load()
        avg := time.Duration(0)
        if started > 0 {
                avg = time.Duration(duration / int64(started))
        }
        return TracerMetrics{SpansStarted: started, SpansFailed: failed, AverageLatency: avg}
}

// WithSpanAttributes merges the provided attributes into the span.
func WithSpanAttributes(attrs map[string]string) SpanFinishOption {
        return func(span *Span) {
                        if len(attrs) == 0 {
                                return
                        }
                        if span.Attributes == nil {
                                span.Attributes = make(map[string]string, len(attrs))
                        }
                        for k, v := range attrs {
                                span.Attributes[k] = v
                        }
        }
}

// WithSpanError records an error for the span.
func WithSpanError(err error) SpanFinishOption {
        return func(span *Span) {
                        span.Error = err
        }
}

// WithSpanStatus attaches a status string to the span.
func WithSpanStatus(status string) SpanFinishOption {
        return func(span *Span) {
                        span.Status = status
        }
}

func (t *Tracer) record(span Span) {
        if span.Attributes != nil {
                copied := make(map[string]string, len(span.Attributes))
                for k, v := range span.Attributes {
                        copied[k] = v
                }
                span.Attributes = copied
        }

        t.metrics.started.Add(1)
        if span.Error != nil {
                t.metrics.failed.Add(1)
        }
        t.metrics.durationNS.Add(span.Duration.Nanoseconds())

        t.mu.Lock()
        if t.historyLimit > 0 && len(t.history) >= t.historyLimit {
                trim := len(t.history) - t.historyLimit + 1
                t.history = append(t.history[trim:], span)
        } else {
                t.history = append(t.history, span)
        }
        exporters := append([]SpanExporter(nil), t.exporters...)
        t.mu.Unlock()

        for _, exporter := range exporters {
                if exporter == nil {
                        continue
                }
                _ = exporter.ExportSpan(span)
        }
}

func (t *Tracer) now() time.Time {
        if t.clock != nil {
                return t.clock()
        }
        return time.Now()
}
