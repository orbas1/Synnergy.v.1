package telemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer returns the global tracer used across the codebase.
func Tracer() trace.Tracer {
	return otel.Tracer("synnergy")
}
