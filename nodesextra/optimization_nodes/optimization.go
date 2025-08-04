package optimizationnodes

import "sync"

// Metrics captures basic runtime statistics used for optimisation decisions.
type Metrics struct {
	CPUUsage    float64 // 0-1 fraction of CPU utilised
	MemoryUsage float64 // 0-1 fraction of memory utilised
	LatencyMS   float64 // average network latency in milliseconds
	Throughput  float64 // transactions per second
}

// Suggestion provides guidance on whether resources should be scaled.
type Suggestion struct {
	ScaleResources bool   // true when additional resources are recommended
	Notes          string // human readable reasoning for the suggestion
}

// SimpleOptimizer implements a trivial optimisation strategy.
type SimpleOptimizer struct {
	mu sync.Mutex
}

// Optimize analyses metrics and returns a scaling suggestion.
func (o *SimpleOptimizer) Optimize(m Metrics) Suggestion {
	o.mu.Lock()
	defer o.mu.Unlock()

	switch {
	case m.CPUUsage > 0.75 || m.MemoryUsage > 0.75:
		return Suggestion{ScaleResources: true, Notes: "high resource usage detected"}
	case m.LatencyMS > 200:
		return Suggestion{ScaleResources: true, Notes: "network latency above threshold"}
	default:
		return Suggestion{ScaleResources: false, Notes: "resources within optimal range"}
	}
}
