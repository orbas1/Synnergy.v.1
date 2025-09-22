package optimizationnodes

import "sync"

// AdaptiveOptimizer implements a rolling-window optimisation strategy that
// considers CPU, memory, latency and error rate trends when recommending
// scaling actions.
type AdaptiveOptimizer struct {
	id           string
	mu           sync.Mutex
	history      []Metrics
	historyLimit int
	baseReplicas int
	lastDecision Suggestion
}

// NewAdaptiveOptimizer constructs an optimizer with the supplied base replicas
// and history size for smoothing decisions.
func NewAdaptiveOptimizer(id string, baseReplicas, historyLimit int) *AdaptiveOptimizer {
	if historyLimit <= 0 {
		historyLimit = 5
	}
	if baseReplicas <= 0 {
		baseReplicas = 1
	}
	return &AdaptiveOptimizer{id: id, historyLimit: historyLimit, baseReplicas: baseReplicas}
}

// ID returns the optimizer identifier.
func (o *AdaptiveOptimizer) ID() string { return o.id }

// Optimize analyses metrics and returns a scaling suggestion.
func (o *AdaptiveOptimizer) Optimize(m Metrics) Suggestion {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.history = append(o.history, m)
	if len(o.history) > o.historyLimit {
		o.history = append([]Metrics(nil), o.history[len(o.history)-o.historyLimit:]...)
	}

	avg := averageMetrics(o.history)
	suggestion := Suggestion{TargetReplicas: o.baseReplicas}

	scaleUpTriggers := 0
	if avg.CPUUsage > 0.75 {
		scaleUpTriggers++
	}
	if avg.MemoryUsage > 0.75 {
		scaleUpTriggers++
	}
	if avg.LatencyMS > 200 {
		scaleUpTriggers++
	}
	if avg.ErrorRate > 0.05 {
		scaleUpTriggers++
	}
	if avg.QueueDepth > 100 {
		scaleUpTriggers++
	}

	if scaleUpTriggers >= 2 {
		suggestion.ScaleResources = true
		suggestion.TargetReplicas = o.baseReplicas + scaleUpTriggers
		suggestion.Notes = "scaling up: sustained pressure across multiple indicators"
	} else if avg.CPUUsage < 0.35 && avg.MemoryUsage < 0.35 && avg.QueueDepth < 10 && avg.ErrorRate < 0.01 {
		suggestion.ScaleResources = false
		if o.baseReplicas > 1 {
			suggestion.TargetReplicas = o.baseReplicas - 1
			suggestion.Notes = "scaling down: utilisation well below thresholds"
		} else {
			suggestion.Notes = "resources within optimal range"
		}
	} else {
		suggestion.ScaleResources = false
		suggestion.Notes = "resources within optimal range"
	}

	o.lastDecision = suggestion
	return suggestion
}

// RecordDecision stores the last applied decision. Tooling can call this after
// applying a scaling operation to keep the optimizer aware of external changes.
func (o *AdaptiveOptimizer) RecordDecision(s Suggestion) {
	o.mu.Lock()
	o.lastDecision = s
	o.mu.Unlock()
}

// LastDecision returns the most recently recorded decision.
func (o *AdaptiveOptimizer) LastDecision() Suggestion {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.lastDecision
}

func averageMetrics(history []Metrics) Metrics {
	if len(history) == 0 {
		return Metrics{}
	}
	var sum Metrics
	var queueSum float64
	for _, m := range history {
		sum.CPUUsage += m.CPUUsage
		sum.MemoryUsage += m.MemoryUsage
		sum.LatencyMS += m.LatencyMS
		sum.Throughput += m.Throughput
		sum.ErrorRate += m.ErrorRate
		queueSum += float64(m.QueueDepth)
	}
	count := float64(len(history))
	return Metrics{
		CPUUsage:    sum.CPUUsage / count,
		MemoryUsage: sum.MemoryUsage / count,
		LatencyMS:   sum.LatencyMS / count,
		Throughput:  sum.Throughput / count,
		ErrorRate:   sum.ErrorRate / count,
		QueueDepth:  int(queueSum / count),
	}
}
