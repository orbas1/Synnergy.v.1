package p2p

import (
	"context"
	"testing"
	"time"

	"synnergy/internal/security"
)

type staticResolver struct {
	peers []Peer
	err   error
}

func (s *staticResolver) Discover(ctx context.Context, seed Peer) ([]Peer, error) {
	return s.peers, s.err
}

func TestDiscoveryServiceBootstrapAndResolver(t *testing.T) {
	mitigator := security.NewDDoSMitigator(security.MitigationConfig{Window: time.Second, MaxRequests: 10})
	manager := NewManager(mitigator)
	bootstrap := []Peer{{ID: "boot", Address: "10.0.0.1:9000", Capabilities: map[string]bool{"validator": true}}}
	resolver := &staticResolver{peers: []Peer{{ID: "p2", Address: "10.0.0.2:9000", Capabilities: map[string]bool{"validator": true}}}}
	svc := NewDiscoveryService(manager, bootstrap, resolver)
	svc.ConfigureQuorum(2)
	svc.WithFilter(func(p Peer) bool { return p.Capabilities["validator"] })
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	peers, err := svc.Discover(ctx)
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if len(peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(peers))
	}
	metrics := svc.Metrics()
	if metrics.Discovered == 0 || metrics.LastRun.IsZero() {
		t.Fatalf("expected metrics to be populated")
	}
}

func TestDiscoveryServiceQuorumFailure(t *testing.T) {
	manager := NewManager(nil)
	svc := NewDiscoveryService(manager, nil)
	svc.ConfigureQuorum(2)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := svc.Discover(ctx)
	if err == nil || err.Error() != "p2p: quorum not satisfied" {
		t.Fatalf("expected quorum error, got %v", err)
	}
}
