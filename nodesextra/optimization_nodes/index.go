package optimizationnodes

// BaseNode defines the minimal behaviour required by optimisation nodes.
type BaseNode interface {
	// ID returns a unique identifier for the node.
	ID() string
}

// OptimizationNode exposes hooks for runtime performance tuning.
type OptimizationNode interface {
	BaseNode
	// Optimize analyses the provided metrics and applies adjustments.
	Optimize(Metrics) Suggestion
}
