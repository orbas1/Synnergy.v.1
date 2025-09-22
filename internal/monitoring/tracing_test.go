package monitoring

import (
        "errors"
        "testing"
        "time"
)

type recordingExporter struct {
        spans []Span
}

func (r *recordingExporter) ExportSpan(span Span) error {
        r.spans = append(r.spans, span)
        return nil
}

func TestTracerRecordsSpanAndExporter(t *testing.T) {
        clockTimes := []time.Time{
                time.Unix(0, 0),
                time.Unix(1, 0),
        }
        idx := 0
        tracer := NewTracer(WithClock(func() time.Time {
                if idx >= len(clockTimes) {
                        return clockTimes[len(clockTimes)-1]
                }
                t := clockTimes[idx]
                idx++
                return t
        }))
        exporter := &recordingExporter{}
        tracer.RegisterExporter(exporter)

        finish := tracer.StartSpan("operation")
        finish(WithSpanAttributes(map[string]string{"key": "value"}), WithSpanStatus("ok"), WithSpanError(errors.New("boom")))
        finish() // idempotent

        spans := tracer.Spans()
        if len(spans) != 1 {
                t.Fatalf("expected single span recorded, got %d", len(spans))
        }
        span := spans[0]
        if span.Name != "operation" || span.Duration != time.Second {
                t.Fatalf("unexpected span data: %+v", span)
        }
        if span.Attributes["key"] != "value" || span.Status != "ok" || span.Error == nil {
                t.Fatalf("span metadata missing")
        }
        if len(exporter.spans) != 1 {
                t.Fatalf("exporter should capture span")
        }

        metrics := tracer.Metrics()
        if metrics.SpansStarted != 1 || metrics.SpansFailed != 1 {
                t.Fatalf("unexpected metrics: %+v", metrics)
        }
        if metrics.AverageLatency != time.Second {
                t.Fatalf("expected average latency to be 1s, got %s", metrics.AverageLatency)
        }
}

func TestTracerHistoryLimit(t *testing.T) {
        tracer := NewTracer(WithHistoryLimit(2))
        tracer.StartSpan("one")()
        tracer.StartSpan("two")()
        tracer.StartSpan("three")()

        spans := tracer.Spans()
        if len(spans) != 2 {
                t.Fatalf("expected history to cap at 2, got %d", len(spans))
        }
        if spans[0].Name != "two" || spans[1].Name != "three" {
                t.Fatalf("unexpected spans retained: %+v", spans)
        }
}

func TestTracerResetClearsState(t *testing.T) {
        tracer := NewTracer()
        tracer.StartSpan("reset")()
        tracer.Reset()

        if len(tracer.Spans()) != 0 {
                t.Fatalf("expected history to be cleared")
        }
        metrics := tracer.Metrics()
        if metrics.SpansStarted != 0 || metrics.SpansFailed != 0 || metrics.AverageLatency != 0 {
                t.Fatalf("metrics should reset: %+v", metrics)
        }
}
