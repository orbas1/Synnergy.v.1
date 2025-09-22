package synnergy

import "testing"

func TestCrossChainManagerLifecycle(t *testing.T) {
	mgr := NewCrossChainManager()
	events := make(chan BridgeEvent, 32)
	mgr.RegisterListener(func(ev BridgeEvent) { events <- ev })

	id := mgr.RegisterBridge("chainA", "chainB", "relayer-1", WithBridgeMetadata("policy", "kyc"))
	if id == "" {
		t.Fatalf("expected bridge id")
	}
	bridge, ok := mgr.GetBridge(id)
	if !ok {
		t.Fatalf("bridge not found")
	}
	if !bridge.Active {
		t.Fatalf("bridge should be active")
	}
	if bridge.Metadata["policy"] != "kyc" {
		t.Fatalf("metadata missing")
	}

	if err := mgr.AuthorizeBridgeRelayer(id, "relayer-2"); err != nil {
		t.Fatalf("authorize relayer: %v", err)
	}
	if err := mgr.DeactivateBridge(id); err != nil {
		t.Fatalf("deactivate: %v", err)
	}
	if err := mgr.ActivateBridge(id); err != nil {
		t.Fatalf("activate: %v", err)
	}
	if err := mgr.UpdateBridgeMetadata(id, "status", "stable"); err != nil {
		t.Fatalf("update metadata: %v", err)
	}
	mgr.RevokeRelayer("relayer-1")

	if mgr.IsRelayerAuthorized("relayer-1") {
		t.Fatalf("revoked relayer should not be authorized")
	}

	metrics := mgr.Metrics()
	if metrics.Total != 1 || metrics.Active != 1 {
		t.Fatalf("unexpected metrics %+v", metrics)
	}
	if metrics.AuthorizedRelay != 2 {
		t.Fatalf("expected 2 authorized relayers, got %d", metrics.AuthorizedRelay)
	}
	if metrics.RevokedRelay != 1 {
		t.Fatalf("expected 1 revoked relayer, got %d", metrics.RevokedRelay)
	}

	close(events)
	seen := map[BridgeEventType]int{}
	for ev := range events {
		seen[ev.Type]++
	}
	required := []BridgeEventType{BridgeEventRegistered, BridgeEventRelayerAuthorized, BridgeEventDeactivated, BridgeEventActivated, BridgeEventMetadataUpdated, BridgeEventRelayerRevoked}
	for _, typ := range required {
		if seen[typ] == 0 {
			t.Fatalf("expected event %s", typ)
		}
	}
}
