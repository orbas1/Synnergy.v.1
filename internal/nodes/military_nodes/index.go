package militarynodes

import (
	"sync"
	"time"
)

// BaseNode defines minimal functionality expected from a network node.
type BaseNode interface {
	// GetID returns the node identifier.
	GetID() string
}

// LogisticsRecord captures movement or status updates for military assets.
// It is declared in the interface package so that both interface definitions
// and concrete implementations share the same type without introducing
// circular dependencies.
type LogisticsRecord struct {
	AssetID   string
	Location  string
	Status    string
	Timestamp time.Time
}

// WarfareNode extends a base node with military specific operations.
type WarfareNode interface {
	BaseNode
	// SecureCommand executes a privileged command after verifying
	// appropriate authorization. Implementations should ensure commands
	// are transmitted using secure channels.
	SecureCommand(cmd string) error
	// TrackLogistics records movement or status changes for a military asset.
	// assetID uniquely identifies the asset, while location and status capture
	// current logistics information.
	TrackLogistics(assetID, location, status string)
	// ShareTactical distributes tactical information to allied nodes or systems.
	ShareTactical(info string)
	// Logistics returns a copy of stored logistics records for inspection or
	// auditing purposes.
	Logistics() []LogisticsRecord
}

// SimpleWarfareNode offers an in-memory implementation of the WarfareNode
// interface suitable for unit tests and demonstrations.
type SimpleWarfareNode struct {
	id   string
	mu   sync.RWMutex
	logs []LogisticsRecord
}

// NewSimpleWarfareNode creates a new warfare node with the supplied identifier.
func NewSimpleWarfareNode(id string) *SimpleWarfareNode {
	return &SimpleWarfareNode{id: id}
}

// GetID returns the node identifier.
func (n *SimpleWarfareNode) GetID() string { return n.id }

// SecureCommand records execution of a privileged command. In a real
// implementation this would include cryptographic authentication.
func (n *SimpleWarfareNode) SecureCommand(cmd string) error {
	// Placeholder for security logic.
	return nil
}

// TrackLogistics appends a logistics record for an asset.
func (n *SimpleWarfareNode) TrackLogistics(assetID, location, status string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.logs = append(n.logs, LogisticsRecord{
		AssetID:   assetID,
		Location:  location,
		Status:    status,
		Timestamp: time.Now(),
	})
}

// ShareTactical is a stub for distributing tactical information.
func (n *SimpleWarfareNode) ShareTactical(info string) {}

// Logistics returns a copy of recorded logistics data.
func (n *SimpleWarfareNode) Logistics() []LogisticsRecord {
	n.mu.RLock()
	defer n.mu.RUnlock()
	cp := make([]LogisticsRecord, len(n.logs))
	copy(cp, n.logs)
	return cp
}
