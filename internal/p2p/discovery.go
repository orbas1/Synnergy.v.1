package p2p

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Resolver resolves new peers given a seed peer. CLI and UI components provide
// custom resolvers (DNS, REST, gossip) to extend discovery without modifying the
// core manager.
type Resolver interface {
	Discover(ctx context.Context, seed Peer) ([]Peer, error)
}

// DiscoveryMetrics provides insight into the last discovery run for dashboards
// and the CLI.
type DiscoveryMetrics struct {
	LastRun       time.Time
	LastError     error
	Discovered    int
	BootstrapUsed bool
}

// DiscoveryService performs layered peer discovery using bootstrap nodes, the
// in-memory registry and pluggable resolvers.
type DiscoveryService struct {
	manager   *Manager
	bootstrap []Peer
	resolvers []Resolver
	filter    func(Peer) bool

	mu       sync.RWMutex
	metrics  DiscoveryMetrics
	quorums  int
	required int
}

// NewDiscoveryService constructs the service. When no resolvers are provided the
// service will rely purely on bootstrap nodes and existing manager state.
func NewDiscoveryService(m *Manager, bootstrap []Peer, resolvers ...Resolver) *DiscoveryService {
	return &DiscoveryService{
		manager:   m,
		bootstrap: append([]Peer(nil), bootstrap...),
		resolvers: append([]Resolver(nil), resolvers...),
		required:  1,
	}
}

// WithFilter configures a discovery filter, allowing callers to restrict peers
// to particular capabilities or regions.
func (d *DiscoveryService) WithFilter(filter func(Peer) bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.filter = filter
}

// ConfigureQuorum instructs the discovery service to ensure at least n peers are
// available locally before returning.
func (d *DiscoveryService) ConfigureQuorum(required int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if required <= 0 {
		required = 1
	}
	d.required = required
}

// Discover performs a synchronous discovery pass using the provided context.
func (d *DiscoveryService) Discover(ctx context.Context) ([]Peer, error) {
	d.mu.RLock()
	filter := d.filter
	resolvers := append([]Resolver(nil), d.resolvers...)
	bootstrap := append([]Peer(nil), d.bootstrap...)
	required := d.required
	d.mu.RUnlock()

	discovered := make(map[string]Peer)
	peers := d.manager.ListPeers()
	metrics := DiscoveryMetrics{BootstrapUsed: len(peers) == 0}
	if len(peers) == 0 {
		for _, peer := range bootstrap {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			peer := d.manager.AddPeer(peer)
			discovered[peer.ID] = peer
		}
		peers = append(peers, bootstrap...)
	}

	var lastErr error
	for _, resolver := range resolvers {
		for _, seed := range peers {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			results, err := resolver.Discover(ctx, seed)
			if err != nil {
				lastErr = err
				continue
			}
			for _, candidate := range results {
				if filter != nil && !filter(candidate) {
					continue
				}
				if _, ok := discovered[candidate.ID]; ok {
					continue
				}
				if _, exists := d.manager.GetPeer(candidate.ID); exists {
					continue
				}
				if d.manager.ddos != nil && !d.manager.ddos.Allow(candidate.Address, time.Now().UTC()) {
					continue
				}
				peer := d.manager.AddPeer(candidate)
				discovered[peer.ID] = peer
			}
		}
	}

	metrics.Discovered = len(discovered)
	metrics.LastRun = time.Now().UTC()
	metrics.LastError = lastErr
	d.mu.Lock()
	d.metrics = metrics
	d.mu.Unlock()

	peers = d.manager.ListPeers()
	if len(peers) < required {
		return peers, errors.New("p2p: quorum not satisfied")
	}
	return peers, lastErr
}

// DiscoverPeers performs discovery with a background context, ignoring any
// error. It is used in compatibility layers.
func (d *DiscoveryService) DiscoverPeers() []Peer {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	peers, _ := d.Discover(ctx)
	return peers
}

// Metrics returns the latest discovery metrics.
func (d *DiscoveryService) Metrics() DiscoveryMetrics {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.metrics
}
