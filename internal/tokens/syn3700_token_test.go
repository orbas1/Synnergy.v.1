package tokens

import "testing"

func TestSYN3700IndexOperations(t *testing.T) {
	token := NewSYN3700Token("Syn Index", "SYNX")
	if err := token.AddComponent("SYN", 0.5); err != nil {
		t.Fatalf("add component: %v", err)
	}
	if err := token.AddComponent("BTC", 0.3); err != nil {
		t.Fatalf("add component: %v", err)
	}
	if err := token.UpdateComponent("BTC", 0.4); err != nil {
		t.Fatalf("update component: %v", err)
	}
	token.NormalizeWeights()
	weight, err := token.ComponentWeight("BTC")
	if err != nil {
		t.Fatalf("component weight: %v", err)
	}
	if weight <= 0 {
		t.Fatalf("expected positive weight")
	}
	prices := map[string]float64{"SYN": 2, "BTC": 10}
	value, contributions, err := token.ValueDetailed(prices)
	if err != nil {
		t.Fatalf("value detailed: %v", err)
	}
	if value <= 0 {
		t.Fatalf("expected positive value")
	}
	if _, ok := contributions["BTC"]; !ok {
		t.Fatalf("expected BTC contribution")
	}
	if err := token.RemoveComponent("SYN"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, err := token.ComponentWeight("SYN"); err != ErrComponentNotFound {
		t.Fatalf("expected ErrComponentNotFound, got %v", err)
	}
	if _, _, err := token.ValueDetailed(map[string]float64{"BTC": 10}); err != nil {
		t.Fatalf("value detailed after removal: %v", err)
	}
}
