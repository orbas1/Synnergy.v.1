package monitoring

import "testing"

func TestTracerStartSpan(t *testing.T) {
	tr := NewTracer()
	end := tr.StartSpan("op")
	end() // should not panic
}
