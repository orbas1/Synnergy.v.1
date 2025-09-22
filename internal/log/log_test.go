package log

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestLoggerWritesJSON(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLogger(Settings{Level: DebugLevel, Writers: []io.Writer{buf}, StaticFields: map[string]any{"service": "test"}})
	logger.Info("component ready", "status", "ok")
	if buf.Len() == 0 {
		t.Fatalf("expected log output")
	}
	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if payload["msg"] != "component ready" || payload["status"] != "ok" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
	if payload["service"] != "test" {
		t.Fatalf("expected static field propagated")
	}
}

func TestLoggerFiltersLevels(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLogger(Settings{Level: ErrorLevel, Writers: []io.Writer{buf}})
	logger.Info("suppressed")
	if buf.Len() != 0 {
		t.Fatalf("expected info log suppressed at error level")
	}
}

func TestLoggerIncludesTraceContext(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLogger(Settings{Level: DebugLevel, Writers: []io.Writer{buf}})
	sc := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: 1, Remote: false})
	ctx := trace.ContextWithSpanContext(context.Background(), sc)
	logger.WithContext(ctx).Info("with trace")
	if buf.Len() == 0 {
		t.Fatalf("expected log output")
	}
	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if payload["trace_id"] == nil || payload["span_id"] == nil {
		t.Fatalf("expected trace context fields")
	}
}

func TestLoggerTextFormat(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := NewLogger(Settings{Level: InfoLevel, Format: "text", Writers: []io.Writer{buf}})
	logger.Info("hello", "k", "v")
	if !strings.Contains(buf.String(), "hello k=v") {
		t.Fatalf("expected text formatted output, got %s", buf.String())
	}
}

func TestDefaultLoggerSingleton(t *testing.T) {
	Configure(Settings{Level: InfoLevel})
	first := Default()
	second := Default()
	if first != second {
		t.Fatalf("expected default logger singleton")
	}
}
