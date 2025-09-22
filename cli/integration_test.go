package cli

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestIntegrationStatusCommandJSON(t *testing.T) {
	cmd := newIntegrationStatusCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v -- %s", err, buf.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("unmarshal: %v -- %s", err, buf.String())
	}
	if _, ok := payload["diagnostics"]; !ok {
		t.Fatalf("expected diagnostics in payload")
	}
	if _, ok := payload["enterprise"]; !ok {
		t.Fatalf("expected enterprise readiness in payload")
	}
}

func TestIntegrationStatusCommandTable(t *testing.T) {
	cmd := newIntegrationStatusCommand()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	if err := cmd.Flags().Set("format", "table"); err != nil {
		t.Fatalf("set flag: %v", err)
	}
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v -- %s", err, buf.String())
	}
	out := buf.String()
	if !strings.Contains(out, "Integration Health") {
		t.Fatalf("expected table header, got %s", out)
	}
	if !strings.Contains(out, "Enterprise readiness") {
		t.Fatalf("expected enterprise readiness section, got %s", out)
	}
}
