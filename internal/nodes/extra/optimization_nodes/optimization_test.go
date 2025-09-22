package optimizationnodes

import "testing"

func TestAdaptiveOptimizerScaleUp(t *testing.T) {
	opt := NewAdaptiveOptimizer("opt", 2, 3)
	metrics := Metrics{CPUUsage: 0.8, MemoryUsage: 0.6, LatencyMS: 250, ErrorRate: 0.1, QueueDepth: 150}
	suggestion := opt.Optimize(metrics)
	if !suggestion.ScaleResources || suggestion.TargetReplicas <= 2 {
		t.Fatalf("expected scale up suggestion %+v", suggestion)
	}
}

func TestAdaptiveOptimizerScaleDown(t *testing.T) {
	opt := NewAdaptiveOptimizer("opt", 3, 3)
	metrics := Metrics{CPUUsage: 0.2, MemoryUsage: 0.2, LatencyMS: 100, ErrorRate: 0.0, QueueDepth: 2}
	suggestion := opt.Optimize(metrics)
	if suggestion.ScaleResources {
		t.Fatalf("expected no scale up suggestion %+v", suggestion)
	}
	if suggestion.TargetReplicas != 2 {
		t.Fatalf("expected to scale down to 2 replicas got %+v", suggestion)
	}
}
