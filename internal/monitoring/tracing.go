package monitoring

// Tracer is a placeholder for distributed tracing.
type Tracer struct{}

// NewTracer creates a new Tracer.
func NewTracer() *Tracer { return &Tracer{} }

// StartSpan is a placeholder that would start a trace span.
func (t *Tracer) StartSpan(name string) func() {
	return func() {}
}
