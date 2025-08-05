package synnergy

import (
	"testing"
)

func TestMobileMiningNodeBatteryThreshold(t *testing.T) {
	m := NewMobileMiningNode("m1", 1, 0.5, 0.3)
	if err := m.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	if !m.IsRunning() {
		t.Fatalf("expected running")
	}
	m.UpdateBattery(0.2)
	if m.IsRunning() {
		t.Fatalf("expected stopped after low battery")
	}
	if m.Battery() != 0.2 {
		t.Fatalf("battery not updated: %v", m.Battery())
	}
	m.UpdateBattery(0.8)
	if err := m.Start(); err != nil {
		t.Fatalf("restart: %v", err)
	}
	if !m.IsRunning() {
		t.Fatalf("expected running after restart")
	}
	if m.Threshold() != 0.3 {
		t.Fatalf("threshold mismatch")
	}
	m.SetThreshold(0.9)
	if m.Threshold() != 0.9 {
		t.Fatalf("threshold not updated")
	}
	m.UpdateBattery(0.8)
	if m.IsRunning() {
		t.Fatalf("expected stopped after threshold raise")
	}
}

func TestMobileMiningNodeStartError(t *testing.T) {
	m := NewMobileMiningNode("m2", 1, 0.2, 0.5)
	if err := m.Start(); err == nil {
		t.Fatalf("expected error for low battery")
	}
	if m.IsRunning() {
		t.Fatalf("should not be running")
	}
}
