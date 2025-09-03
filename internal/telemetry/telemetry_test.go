package telemetry

import "testing"

func TestTracerNotNil(t *testing.T) {
	if Tracer() == nil {
		t.Fatalf("expected non-nil tracer")
	}
}
