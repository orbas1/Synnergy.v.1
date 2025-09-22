package auth

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStdAuditLoggerMasksSensitiveFields(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStdAuditLogger(&buf)
	logger.Log("user", Permission("write"), true, map[string]any{"token": "abc", "detail": "ok"})
	out := strings.TrimSpace(buf.String())
	if out == "" {
		t.Fatalf("expected log output")
	}
	// log lines include timestamp prefix; find JSON payload
	idx := strings.Index(out, "{")
	if idx == -1 {
		t.Fatalf("expected JSON payload")
	}
	var event AuditEvent
	if err := json.Unmarshal([]byte(out[idx:]), &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if event.Metadata["token"] != "***" {
		t.Fatalf("expected token to be masked, got %v", event.Metadata["token"])
	}
	if event.Metadata["detail"] != "ok" {
		t.Fatalf("expected detail to remain, got %v", event.Metadata["detail"])
	}
}

func TestMemoryAuditLoggerStoresEvents(t *testing.T) {
	logger := NewMemoryAuditLogger()
	logger.Log("user", Permission("read"), false, map[string]any{"reason": "denied"})
	events := logger.Events()
	if len(events) != 1 {
		t.Fatalf("expected one event")
	}
	if events[0].Allowed {
		t.Fatalf("expected allowed=false")
	}
	if events[0].Metadata["reason"] != "denied" {
		t.Fatalf("unexpected metadata: %v", events[0].Metadata)
	}
}
