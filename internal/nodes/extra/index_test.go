package nodes

import (
	"errors"
	"testing"
)

type stubNode struct {
	id        string
	running   bool
	failStart bool
	failStop  bool
}

func (s *stubNode) ID() string { return s.id }

func (s *stubNode) Start() error {
	if s.failStart {
		return errors.New("start failed")
	}
	s.running = true
	return nil
}

func (s *stubNode) Stop() error {
	if s.failStop {
		return errors.New("stop failed")
	}
	s.running = false
	return nil
}

func TestRegistryRegisterStartStop(t *testing.T) {
	reg := NewRegistry()
	node := &stubNode{id: "n1"}
	if err := reg.Register(node); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Register(node); err == nil {
		t.Fatalf("expected duplicate registration error")
	}

	if err := reg.Start("n1"); err != nil {
		t.Fatalf("start: %v", err)
	}
	metrics, ok := reg.Metrics("n1")
	if !ok || !metrics.Running {
		t.Fatalf("unexpected metrics %+v", metrics)
	}
	if err := reg.Stop("n1"); err != nil {
		t.Fatalf("stop: %v", err)
	}
	metrics, ok = reg.Metrics("n1")
	if !ok || metrics.Running {
		t.Fatalf("expected node stopped got %+v", metrics)
	}
}

func TestRegistryStartAllStopAll(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(&stubNode{id: "a"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Register(&stubNode{id: "b"}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.StartAll(); err != nil {
		t.Fatalf("start all: %v", err)
	}
	if err := reg.StopAll(); err != nil {
		t.Fatalf("stop all: %v", err)
	}
}

func TestRegistryErrorPropagation(t *testing.T) {
	reg := NewRegistry()
	if err := reg.Register(&stubNode{id: "n1", failStart: true}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Start("n1"); err == nil {
		t.Fatalf("expected start error")
	}
	if err := reg.Register(&stubNode{id: "n2", failStop: true}); err != nil {
		t.Fatalf("register: %v", err)
	}
	if err := reg.Stop("n2"); err == nil {
		t.Fatalf("expected stop error")
	}
}
