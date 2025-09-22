package tokens

import (
	"sync"
	"testing"
	"time"
)

func TestRegistryLifecycle(t *testing.T) {
	var mu sync.Mutex
	events := make([]RegistryEvent, 0)
	reg := NewRegistry(WithRegistryClock(func() time.Time { return time.Unix(0, 0) }))
	reg.RegisterObserver(func(evt RegistryEvent) {
		mu.Lock()
		events = append(events, evt)
		mu.Unlock()
	})
	tok := NewBaseToken(reg.NextID(), "Token", "TOK", 0)
	reg.Register(tok)
	if info, ok := reg.Info(tok.ID()); !ok || info.Symbol != "TOK" {
		t.Fatalf("unexpected info %+v ok=%v", info, ok)
	}
	if listed := reg.List(); len(listed) != 1 || listed[0].ID != tok.ID() {
		t.Fatalf("unexpected list %+v", listed)
	}
	if removed := reg.Remove(tok.ID()); !removed {
		t.Fatalf("expected removal")
	}
	if _, ok := reg.Get(tok.ID()); ok {
		t.Fatalf("token should be absent after removal")
	}
	mu.Lock()
	if len(events) != 2 || events[0].Type != RegistryEventRegistered || events[1].Type != RegistryEventRemoved {
		t.Fatalf("unexpected events %+v", events)
	}
	mu.Unlock()
}

func TestRegistryGetBySymbol(t *testing.T) {
	reg := NewRegistry()
	first := NewBaseToken(reg.NextID(), "First", "FST", 0)
	second := NewBaseToken(reg.NextID(), "Second", "SND", 0)
	reg.Register(first)
	reg.Register(second)
	if tok, ok := reg.GetBySymbol("SND"); !ok || tok.ID() != second.ID() {
		t.Fatalf("unexpected token %+v ok=%v", tok, ok)
	}
	if !reg.Remove(second.ID()) {
		t.Fatalf("expected removal")
	}
	if len(reg.List()) != 1 {
		t.Fatalf("expected registry to contain single token")
	}
}
