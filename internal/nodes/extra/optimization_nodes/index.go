package optimizationnodes

// BaseNode defines the minimal behaviour required by optimisation nodes.
type BaseNode interface {
	// ID returns a unique identifier for the node.
	ID() string
}

// Metrics captures runtime statistics used for optimisation decisions.
type Metrics struct {
	CPUUsage    float64 // 0-1 fraction of CPU utilised
	MemoryUsage float64 // 0-1 fraction of memory utilised
	LatencyMS   float64 // average network latency in milliseconds
	Throughput  float64 // transactions per second
	ErrorRate   float64 // fraction of failed operations
	QueueDepth  int     // pending operations in queue
}

// Suggestion provides guidance on whether resources should be scaled.
type Suggestion struct {
	ScaleResources bool
	TargetReplicas int
	Notes          string
}

// OptimizationNode exposes hooks for runtime performance tuning.
type OptimizationNode interface {
	BaseNode
	// Optimize analyses the provided metrics and applies adjustments.
	Optimize(Metrics) Suggestion
	// RecordDecision allows implementations to persist the applied change.
	RecordDecision(Suggestion)
}
