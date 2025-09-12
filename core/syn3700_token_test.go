package core

import (
	"sync"
	"testing"
)

func TestSYN3700TokenLifecycle(t *testing.T) {
	tok := NewSYN3700Token("Index", "IDX")
	tok.AddComponent("AAA", 0.5)
	tok.AddComponent("BBB", 0.5)
	if err := tok.RemoveComponent("CCC"); err == nil {
		t.Fatalf("expected missing component error")
	}
	if err := tok.RemoveComponent("BBB"); err != nil {
		t.Fatalf("remove: %v", err)
	}
	comps := tok.ListComponents()
	if len(comps) != 1 || comps[0].Token != "AAA" {
		t.Fatalf("unexpected components %+v", comps)
	}
	val := tok.Value(map[string]float64{"AAA": 2})
	if val != 1 { // 0.5 * 2
		t.Fatalf("expected value 1 got %f", val)
	}
}

func TestSYN3700TokenConcurrentAdd(t *testing.T) {
	tok := NewSYN3700Token("Index", "IDX")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tok.AddComponent(string(rune('A'+i)), 0.1)
		}(i)
	}
	wg.Wait()
	if len(tok.ListComponents()) != 10 {
		t.Fatalf("expected 10 components")
	}
}
