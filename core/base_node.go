package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"
	"time"

	"synnergy/internal/nodes"
)

// BaseNode wraps a NodeInterface and exposes common networking behaviour.
type BaseNode struct {
	id               nodes.Address
	peers            map[nodes.Address]peerRecord
	running          bool
	mu               sync.RWMutex
	maxPeers         int
	peerTTL          time.Duration
	failureThreshold int
}

type peerRecord struct {
	lastSeen   time.Time
	failures   int
	persistent bool
}

const (
	defaultMaxPeers         = 1024
	defaultPeerTTL          = 15 * time.Minute
	defaultFailureThreshold = 3
)

// NewBaseNode constructs a BaseNode with the provided identifier.
func NewBaseNode(id nodes.Address) *BaseNode {
	return &BaseNode{
		id:               id,
		peers:            make(map[nodes.Address]peerRecord),
		maxPeers:         defaultMaxPeers,
		peerTTL:          defaultPeerTTL,
		failureThreshold: defaultFailureThreshold,
	}
}

// ID returns the node identifier.
func (n *BaseNode) ID() nodes.Address { return n.id }

// Start marks the node as running.
func (n *BaseNode) Start() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.running {
		return nil
	}
	n.running = true
	return nil
}

// Stop halts node operations.
func (n *BaseNode) Stop() error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return nil
	}
	n.running = false
	return nil
}

// IsRunning reports whether the node is currently active.
func (n *BaseNode) IsRunning() bool {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.running
}

// Peers returns the list of known peers.
func (n *BaseNode) Peers() []nodes.Address {
	snaps := n.PeerSnapshots()
	out := make([]nodes.Address, len(snaps))
	for i, snap := range snaps {
		out[i] = snap.Address
	}
	return out
}

// DialSeed records a connection to a seed peer.
func (n *BaseNode) DialSeed(addr nodes.Address) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return fmt.Errorf("node not running")
	}
	now := time.Now()
	if !n.ensureCapacityLocked(now, addr) {
		return fmt.Errorf("peer capacity reached")
	}
	rec := n.peers[addr]
	rec.lastSeen = now
	rec.failures = 0
	n.peers[addr] = rec
	return nil
}

// DialSeedSigned records a connection to a seed peer after verifying the
// provided signature matches the peer's address.
func (n *BaseNode) DialSeedSigned(addr nodes.Address, sig []byte, pub ed25519.PublicKey) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.running {
		return fmt.Errorf("node not running")
	}
	if hex.EncodeToString(pub) != string(addr) {
		return fmt.Errorf("address mismatch")
	}
	if !ed25519.Verify(pub, []byte(addr), sig) {
		return fmt.Errorf("invalid signature")
	}
	now := time.Now()
	if !n.ensureCapacityLocked(now, addr) {
		return fmt.Errorf("peer capacity reached")
	}
	rec := n.peers[addr]
	rec.lastSeen = now
	rec.failures = 0
	n.peers[addr] = rec
	return nil
}

// RecordPeerHeartbeat updates activity metadata for the provided peer. It
// returns true if the peer is known to the node.
func (n *BaseNode) RecordPeerHeartbeat(addr nodes.Address) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	rec, ok := n.peers[addr]
	if !ok {
		return false
	}
	rec.lastSeen = time.Now()
	rec.failures = 0
	n.peers[addr] = rec
	return true
}

// RecordPeerFailure increments the failure counter for the given peer. When the
// failure threshold is exceeded the peer is removed from the registry and false
// is returned.
func (n *BaseNode) RecordPeerFailure(addr nodes.Address) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	rec, ok := n.peers[addr]
	if !ok {
		return false
	}
	rec.failures++
	if n.failureThreshold > 0 && rec.failures >= n.failureThreshold {
		delete(n.peers, addr)
		return false
	}
	n.peers[addr] = rec
	return true
}

// RemovePeer deletes a peer from the registry.
func (n *BaseNode) RemovePeer(addr nodes.Address) {
	n.mu.Lock()
	delete(n.peers, addr)
	n.mu.Unlock()
}

// PromotePeer marks a peer as persistent so it will not be evicted during
// pruning.
func (n *BaseNode) PromotePeer(addr nodes.Address) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	rec, ok := n.peers[addr]
	if !ok {
		return false
	}
	rec.persistent = true
	n.peers[addr] = rec
	return true
}

// DemotePeer removes the persistent flag from a peer allowing it to be evicted
// by TTL pruning.
func (n *BaseNode) DemotePeer(addr nodes.Address) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	rec, ok := n.peers[addr]
	if !ok {
		return false
	}
	rec.persistent = false
	n.peers[addr] = rec
	return true
}

// SetPeerRetention updates the peer retention policy. Values less than or equal
// to zero leave the existing configuration unchanged.
func (n *BaseNode) SetPeerRetention(maxPeers int, ttl time.Duration) {
	n.mu.Lock()
	if maxPeers > 0 {
		n.maxPeers = maxPeers
	}
	if ttl > 0 {
		n.peerTTL = ttl
	}
	n.enforceCapacityLocked(time.Now())
	n.mu.Unlock()
}

// SetFailureThreshold adjusts how many consecutive failures are tolerated
// before a peer is pruned.
func (n *BaseNode) SetFailureThreshold(threshold int) {
	n.mu.Lock()
	if threshold > 0 {
		n.failureThreshold = threshold
	}
	n.mu.Unlock()
}

// PruneExpiredPeers removes any peers that have not been seen within the
// configured TTL. The number of removed peers is returned.
func (n *BaseNode) PruneExpiredPeers() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.pruneExpiredLocked(time.Now())
}

// PeerSnapshot captures metadata about a tracked peer.
type PeerSnapshot struct {
	Address    nodes.Address
	LastSeen   time.Time
	Failures   int
	Persistent bool
}

// PeerSnapshots returns a stable, sorted view of peer metadata. Peers are
// ordered by most recent activity.
func (n *BaseNode) PeerSnapshots() []PeerSnapshot {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.pruneExpiredLocked(time.Now())
	snaps := make([]PeerSnapshot, 0, len(n.peers))
	for addr, rec := range n.peers {
		snaps = append(snaps, PeerSnapshot{
			Address:    addr,
			LastSeen:   rec.lastSeen,
			Failures:   rec.failures,
			Persistent: rec.persistent,
		})
	}
	sort.Slice(snaps, func(i, j int) bool {
		if snaps[i].LastSeen.Equal(snaps[j].LastSeen) {
			return snaps[i].Address < snaps[j].Address
		}
		return snaps[i].LastSeen.After(snaps[j].LastSeen)
	})
	return snaps
}

func (n *BaseNode) ensureCapacityLocked(now time.Time, candidate nodes.Address) bool {
	n.pruneExpiredLocked(now)
	if _, exists := n.peers[candidate]; exists {
		return true
	}
	if len(n.peers) < n.maxPeers {
		return true
	}
	if evicted := n.evictOldestLocked(); evicted {
		return true
	}
	return false
}

func (n *BaseNode) enforceCapacityLocked(now time.Time) {
	n.pruneExpiredLocked(now)
	for len(n.peers) > n.maxPeers {
		if !n.evictOldestLocked() {
			break
		}
	}
}

func (n *BaseNode) evictOldestLocked() bool {
	var oldestAddr nodes.Address
	var oldestTime time.Time
	found := false
	for addr, rec := range n.peers {
		if rec.persistent {
			continue
		}
		if !found || rec.lastSeen.Before(oldestTime) {
			found = true
			oldestAddr = addr
			oldestTime = rec.lastSeen
		}
	}
	if found {
		delete(n.peers, oldestAddr)
	}
	return found
}

func (n *BaseNode) pruneExpiredLocked(now time.Time) int {
	if n.peerTTL <= 0 {
		return 0
	}
	removed := 0
	for addr, rec := range n.peers {
		if rec.persistent {
			continue
		}
		if now.Sub(rec.lastSeen) > n.peerTTL {
			delete(n.peers, addr)
			removed++
		}
	}
	return removed
}

var _ nodes.NodeInterface = (*BaseNode)(nil)
