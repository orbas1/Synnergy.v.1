package p2p

// DiscoveryService provides basic peer discovery using bootstrap nodes.
type DiscoveryService struct {
	manager   *Manager
	bootstrap []Peer
}

// NewDiscoveryService creates a new discovery service.
func NewDiscoveryService(m *Manager, bootstrap []Peer) *DiscoveryService {
	return &DiscoveryService{manager: m, bootstrap: bootstrap}
}

// DiscoverPeers returns known peers, seeding from bootstrap nodes if necessary.
func (d *DiscoveryService) DiscoverPeers() []Peer {
	peers := d.manager.ListPeers()
	if len(peers) == 0 {
		for _, p := range d.bootstrap {
			d.manager.AddPeer(p)
		}
		peers = append([]Peer(nil), d.bootstrap...)
	}
	return peers
}
