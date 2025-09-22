package militarynodes

import (
	"errors"
	"sync"
	"time"
)

// BaseNode defines minimal functionality expected from a network node.
type BaseNode interface {
	// GetID returns the node identifier.
	GetID() string
}

// LogisticsRecord mirrors the structure used by concrete implementations
// handling military logistics.
type LogisticsRecord struct {
	AssetID   string
	Location  string
	Status    string
	Timestamp time.Time
}

// CommandLog captures privileged actions executed through the node.
type CommandLog struct {
	Command   string
	Operator  string
	Timestamp time.Time
}

// WarfareNode extends a base node with military specific operations.
type WarfareNode interface {
	BaseNode
	// SecureCommand executes a privileged command after verifying
	// appropriate authorization. Implementations should ensure commands
	// are transmitted using secure channels.
	SecureCommand(operator, cmd string) error
	// TrackLogistics records movement or status changes for a military asset.
	// assetID uniquely identifies the asset, while location and status capture
	// current logistics information.
	TrackLogistics(assetID, location, status string)
	// ShareTactical distributes tactical information to allied nodes or systems.
	ShareTactical(info string)
	// Logistics returns a copy of stored logistics records for inspection or
	// auditing purposes.
	Logistics() []LogisticsRecord
	// CommandHistory returns previously executed privileged commands.
	CommandHistory() []CommandLog
}

type warfareConfig struct {
	maxLogistics int
	maxCommands  int
}

// WarfareOption configures optional behaviour for warfare nodes.
type WarfareOption interface {
	applyWarfareOption(*warfareConfig)
}

type warfareOptionFunc func(*warfareConfig)

func (f warfareOptionFunc) applyWarfareOption(cfg *warfareConfig) { f(cfg) }

// WithLogisticsLimit bounds the number of retained logistics records.
func WithLogisticsLimit(limit int) WarfareOption {
	return warfareOptionFunc(func(cfg *warfareConfig) {
		if limit > 0 {
			cfg.maxLogistics = limit
		}
	})
}

// WithCommandLimit bounds the number of retained command records.
func WithCommandLimit(limit int) WarfareOption {
	return warfareOptionFunc(func(cfg *warfareConfig) {
		if limit > 0 {
			cfg.maxCommands = limit
		}
	})
}

// SecureWarfareNode offers an in-memory implementation of the WarfareNode
// interface suitable for unit tests and demonstrations.
type SecureWarfareNode struct {
	id   string
	mu   sync.RWMutex
	cfg  warfareConfig
	logs []LogisticsRecord
	cmds []CommandLog
}

// NewSecureWarfareNode creates a new warfare node with the supplied identifier
// and optional configuration.
func NewSecureWarfareNode(id string, opts ...WarfareOption) *SecureWarfareNode {
	cfg := warfareConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt.applyWarfareOption(&cfg)
		}
	}
	return &SecureWarfareNode{id: id, cfg: cfg}
}

// GetID returns the node identifier.
func (n *SecureWarfareNode) GetID() string { return n.id }

// SecureCommand records execution of a privileged command.
func (n *SecureWarfareNode) SecureCommand(operator, cmd string) error {
	if operator == "" {
		return errors.New("operator is required")
	}
	if cmd == "" {
		return errors.New("command is required")
	}
	n.mu.Lock()
	n.cmds = append(n.cmds, CommandLog{Command: cmd, Operator: operator, Timestamp: time.Now().UTC()})
	if n.cfg.maxCommands > 0 && len(n.cmds) > n.cfg.maxCommands {
		n.cmds = append([]CommandLog(nil), n.cmds[len(n.cmds)-n.cfg.maxCommands:]...)
	}
	n.mu.Unlock()
	return nil
}

// TrackLogistics appends a logistics record for an asset.
func (n *SecureWarfareNode) TrackLogistics(assetID, location, status string) {
	n.mu.Lock()
	n.logs = append(n.logs, LogisticsRecord{
		AssetID:   assetID,
		Location:  location,
		Status:    status,
		Timestamp: time.Now().UTC(),
	})
	if n.cfg.maxLogistics > 0 && len(n.logs) > n.cfg.maxLogistics {
		n.logs = append([]LogisticsRecord(nil), n.logs[len(n.logs)-n.cfg.maxLogistics:]...)
	}
	n.mu.Unlock()
}

// ShareTactical is a stub for distributing tactical information.
func (n *SecureWarfareNode) ShareTactical(info string) {}

// Logistics returns a copy of recorded logistics data.
func (n *SecureWarfareNode) Logistics() []LogisticsRecord {
	n.mu.RLock()
	defer n.mu.RUnlock()
	cp := make([]LogisticsRecord, len(n.logs))
	copy(cp, n.logs)
	return cp
}

// CommandHistory returns a copy of recorded commands.
func (n *SecureWarfareNode) CommandHistory() []CommandLog {
	n.mu.RLock()
	defer n.mu.RUnlock()
	cp := make([]CommandLog, len(n.cmds))
	copy(cp, n.cmds)
	return cp
}
