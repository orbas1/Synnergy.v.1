package errors

import (
	"encoding/json"
	goerrors "errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestErrorWrapping(t *testing.T) {
	base := fmt.Errorf("boom")
	err := Wrap(Invalid, "failed", base, WithSeverity(SeverityWarn), WithCorrelation("corr-123"), WithRetryable(true))
	if err.Error() == "" {
		t.Fatalf("empty error")
	}
	if !IsCode(err, Invalid) {
		t.Fatalf("code mismatch")
	}
	if err.Unwrap() != base {
		t.Fatalf("unwrap mismatch")
	}
	if err.Severity != SeverityWarn {
		t.Fatalf("expected severity warn, got %s", err.Severity)
	}
	if !err.Retryable {
		t.Fatalf("expected retryable flag")
	}
	if err.CorrelationID != "corr-123" {
		t.Fatalf("expected correlation id propagated")
	}
}

func TestMarshalJSONOmitsWrappedError(t *testing.T) {
	err := Wrap(NotFound, "missing resource", fmt.Errorf("cause"))
	payload, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Fatalf("marshal error: %v", marshalErr)
	}
	if string(payload) == "" {
		t.Fatalf("expected json output")
	}
	if goerrors.Is(fmt.Errorf("%s", payload), err.Err) {
		t.Fatalf("wrapped error should not be marshalled")
	}
}

func TestToMapContainsMetadata(t *testing.T) {
	err := Wrap(Unauthorized, "denied", nil, WithDetail("role", "auditor"), WithStatus(http.StatusForbidden))
	err.Timestamp = time.Unix(0, 0).UTC()
	out := err.ToMap()
	if out["status"] != http.StatusForbidden {
		t.Fatalf("expected custom status reflected")
	}
	details := out["details"].(map[string]any)
	if details["role"] != "auditor" {
		t.Fatalf("expected detail preserved")
	}
	if _, ok := out["correlation_id"]; ok {
		t.Fatalf("unexpected correlation id")
	}
	if out["timestamp"].(string) == "" {
		t.Fatalf("timestamp should be present")
	}
}

func TestParseReturnsStructuredError(t *testing.T) {
	original := Wrap(Conflict, "version mismatch", nil)
	wrapped := fmt.Errorf("outer: %w", original)
	parsed, ok := Parse(wrapped)
	if !ok {
		t.Fatalf("expected to parse structured error")
	}
	if parsed.Code != Conflict {
		t.Fatalf("expected conflict code")
	}
}

func TestDefaultStatusMapping(t *testing.T) {
	cases := map[Code]int{
		NotFound:     http.StatusNotFound,
		Invalid:      http.StatusBadRequest,
		Unauthorized: http.StatusUnauthorized,
		Conflict:     http.StatusConflict,
		Internal:     http.StatusInternalServerError,
	}
	for code, status := range cases {
		if got := defaultStatus(code); got != status {
			t.Fatalf("expected %d for code %s got %d", status, code, got)
		}
	}
}
