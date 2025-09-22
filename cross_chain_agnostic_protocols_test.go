package synnergy

import "testing"

func TestProtocolRegistryLifecycle(t *testing.T) {
	reg := NewProtocolRegistry()
	events := make(chan ProtocolEvent, 32)
	reg.RegisterListener(func(ev ProtocolEvent) { events <- ev })

	id := reg.RegisterProtocol("syn-bridge", WithProtocolVersion("1.0.0"), WithProtocolMetadata("purpose", "bridge"))
	if id == "" {
		t.Fatalf("expected protocol id")
	}
	proto, ok := reg.GetProtocol(id)
	if !ok {
		t.Fatalf("protocol not found")
	}
	if proto.Version != "1.0.0" {
		t.Fatalf("version mismatch")
	}
	if proto.Metadata["purpose"] != "bridge" {
		t.Fatalf("metadata missing")
	}

	if _, ok := reg.UpdateProtocol(id, WithProtocolMetadata("region", "eu")); !ok {
		t.Fatalf("update failed")
	}
	if !reg.DeactivateProtocol(id) {
		t.Fatalf("deactivate failed")
	}
	if !reg.ActivateProtocol(id) {
		t.Fatalf("activate failed")
	}

	metrics := reg.Metrics()
	if metrics.Total != 1 || metrics.Active != 1 || metrics.Updates == 0 {
		t.Fatalf("unexpected metrics %+v", metrics)
	}

	close(events)
	seen := map[ProtocolEventType]int{}
	for ev := range events {
		seen[ev.Type]++
	}
	required := []ProtocolEventType{ProtocolEventRegistered, ProtocolEventUpdated, ProtocolEventDeactivated, ProtocolEventActivated}
	for _, typ := range required {
		if seen[typ] == 0 {
			t.Fatalf("expected event %s", typ)
		}
	}
}
