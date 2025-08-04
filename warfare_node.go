package synnergy

import (
	"errors"
	"sync"
	"time"

	militarynodes "synnergy/nodes/military_nodes"
)

// WarfareNode provides military focused extensions on top of a base node.
// Logistics are kept in-memory and protected by a RWMutex to allow
// concurrent reads while serialising writes.
type WarfareNode struct {
	id        string
	mu        sync.RWMutex
	logistics []militarynodes.LogisticsRecord
}

// NewWarfareNode constructs a new WarfareNode with the given identifier.
func NewWarfareNode(id string) *WarfareNode {
	return &WarfareNode{id: id}
}

// GetID satisfies militarynodes.BaseNode.
func (w *WarfareNode) GetID() string { return w.id }

// SecureCommand executes a privileged command after validating input.
// In this stub implementation it simply ensures the command is non-empty.
func (w *WarfareNode) SecureCommand(cmd string) error {
	if cmd == "" {
		return errors.New("empty command")
	}
	return nil
}

// TrackLogistics records a logistics update for a military asset.
func (w *WarfareNode) TrackLogistics(assetID, location, status string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	rec := militarynodes.LogisticsRecord{
		AssetID:   assetID,
		Location:  location,
		Status:    status,
		Timestamp: time.Now().UTC(),
	}
	w.logistics = append(w.logistics, rec)
}

// ShareTactical distributes tactical information. This stub stores no state
// but acts as a hook for future broadcasting logic.
func (w *WarfareNode) ShareTactical(info string) {
	_ = info
}

// Logistics returns a copy of stored logistics records.
func (w *WarfareNode) Logistics() []militarynodes.LogisticsRecord {
	w.mu.RLock()
	defer w.mu.RUnlock()
	cp := make([]militarynodes.LogisticsRecord, len(w.logistics))
	copy(cp, w.logistics)
	return cp
}

// LogisticsByAsset filters stored logistics records for a specific asset ID.
// The returned slice is a copy and may be safely modified by the caller.
func (w *WarfareNode) LogisticsByAsset(assetID string) []militarynodes.LogisticsRecord {
	w.mu.RLock()
	defer w.mu.RUnlock()
	var res []militarynodes.LogisticsRecord
	for _, r := range w.logistics {
		if r.AssetID == assetID {
			res = append(res, r)
		}
	}
	cp := make([]militarynodes.LogisticsRecord, len(res))
	copy(cp, res)
	return cp
}

// ensure interface compliance
var _ militarynodes.WarfareNode = (*WarfareNode)(nil)
