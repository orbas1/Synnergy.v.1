package optimizationnodes

import "testing"

func TestAdaptiveOptimizerImplementsInterface(t *testing.T) {
	var node OptimizationNode = NewAdaptiveOptimizer("opt", 1, 3)
	if node.ID() != "opt" {
		t.Fatalf("unexpected id %s", node.ID())
	}
	suggestion := node.Optimize(Metrics{})
	node.RecordDecision(suggestion)
	if suggestion.TargetReplicas <= 0 {
		t.Fatalf("expected positive replicas")
	}
}
