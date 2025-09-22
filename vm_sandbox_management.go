package synnergy

import (
	"errors"
	"sync"
	"time"
)

// SandboxInfo holds runtime limits and state for a single sandboxed contract
// execution environment.
type SandboxInfo struct {
	ID           string
	ContractAddr string
	GasLimit     uint64
	MemoryLimit  uint64
	Active       bool
	LastReset    time.Time
	CreatedAt    time.Time
}

// SandboxManager manages multiple sandboxes used to isolate contract
// execution. It is safe for concurrent use.
type SandboxManager struct {
	mu        sync.RWMutex
	sandboxes map[string]*SandboxInfo
}

func cloneSandbox(sb *SandboxInfo) *SandboxInfo {
	if sb == nil {
		return nil
	}
	copy := *sb
	return &copy
}

// NewSandboxManager returns an empty manager.
func NewSandboxManager() *SandboxManager {
	return &SandboxManager{sandboxes: make(map[string]*SandboxInfo)}
}

// StartSandbox creates a new sandbox for the given contract address.
func (m *SandboxManager) StartSandbox(id, contractAddr string, gasLimit, memoryLimit uint64) (*SandboxInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.sandboxes[id]; exists {
		return nil, errors.New("sandbox already exists")
	}
	sb := &SandboxInfo{
		ID:           id,
		ContractAddr: contractAddr,
		GasLimit:     gasLimit,
		MemoryLimit:  memoryLimit,
		Active:       true,
		CreatedAt:    time.Now(),
		LastReset:    time.Now(),
	}
	m.sandboxes[id] = sb
	return cloneSandbox(sb), nil
}

// StopSandbox deactivates a sandbox.
func (m *SandboxManager) StopSandbox(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sb, ok := m.sandboxes[id]
	if !ok {
		return errors.New("sandbox not found")
	}
	sb.Active = false
	return nil
}

// DeleteSandbox removes a sandbox entirely, freeing tracking resources.
func (m *SandboxManager) DeleteSandbox(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.sandboxes[id]; !ok {
		return errors.New("sandbox not found")
	}
	delete(m.sandboxes, id)
	return nil
}

// ResetSandbox updates the LastReset timestamp for a sandbox.
func (m *SandboxManager) ResetSandbox(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	sb, ok := m.sandboxes[id]
	if !ok {
		return errors.New("sandbox not found")
	}
	sb.LastReset = time.Now()
	return nil
}

// SandboxStatus returns sandbox information by ID.
func (m *SandboxManager) SandboxStatus(id string) (*SandboxInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	sb, ok := m.sandboxes[id]
	if !ok {
		return nil, false
	}
	return cloneSandbox(sb), true
}

// ListSandboxes returns all sandboxes managed by this instance.
func (m *SandboxManager) ListSandboxes() []*SandboxInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*SandboxInfo, 0, len(m.sandboxes))
	for _, sb := range m.sandboxes {
		out = append(out, cloneSandbox(sb))
	}
	return out
}
