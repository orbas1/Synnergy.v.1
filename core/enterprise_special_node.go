package core

import (
	"errors"
	"fmt"
	"sync"
	"time"

	synn "synnergy"
	"synnergy/internal/nodes"
)

// CombinedNodeRole enumerates the supported roles for enterprise plugins.
type CombinedNodeRole string

const (
	// CombinedRoleConsensus represents validators, full nodes and other
	// consensus-facing infrastructure.
	CombinedRoleConsensus CombinedNodeRole = "consensus"
	// CombinedRoleExecution represents execution engines and rollup
	// sequencers that supply computation.
	CombinedRoleExecution CombinedNodeRole = "execution"
	// CombinedRoleAnalytics represents monitoring, data warehousing or
	// reporting nodes layered on top of consensus data.
	CombinedRoleAnalytics CombinedNodeRole = "analytics"
	// CombinedRoleArchive captures archival and compliance oriented nodes
	// that primarily service historical queries.
	CombinedRoleArchive CombinedNodeRole = "archive"
)

// ParseCombinedNodeRole validates and normalises a textual role.
func ParseCombinedNodeRole(value string) (CombinedNodeRole, error) {
	switch CombinedNodeRole(value) {
	case CombinedRoleConsensus, CombinedRoleExecution, CombinedRoleAnalytics, CombinedRoleArchive:
		return CombinedNodeRole(value), nil
	default:
		return "", fmt.Errorf("unknown combined node role %q", value)
	}
}

// CombinedNodeMetrics summarises runtime statistics supplied by a plugin.
type CombinedNodeMetrics struct {
	MempoolSize    int `json:"mempoolSize"`
	ValidatorCount int `json:"validatorCount"`
	BlockHeight    int `json:"blockHeight"`
	PeerCount      int `json:"peerCount"`
}

// EnterpriseNodePlugin provides the hooks required to fold an external node
// into the combined enterprise surface area.
type EnterpriseNodePlugin struct {
	ID            string
	Role          CombinedNodeRole
	Labels        map[string]string
	Metrics       func() CombinedNodeMetrics
	Broadcast     func(*Transaction) error
	LedgerBalance func(string) uint64
}

// EnterprisePluginSummary exposes plugin metadata for CLI and web dashboards.
type EnterprisePluginSummary struct {
	ID      string              `json:"id"`
	Role    CombinedNodeRole    `json:"role"`
	Labels  map[string]string   `json:"labels,omitempty"`
	Metrics CombinedNodeMetrics `json:"metrics"`
}

// EnterpriseCombinedSnapshot aggregates the metrics of all attached plugins.
type EnterpriseCombinedSnapshot struct {
	Timestamp          time.Time                 `json:"timestamp"`
	NodeCount          int                       `json:"nodeCount"`
	Roles              map[CombinedNodeRole]int  `json:"roles"`
	TotalMempool       int                       `json:"totalMempool"`
	TotalValidators    int                       `json:"totalValidators"`
	HighestBlockHeight int                       `json:"highestBlockHeight"`
	TotalPeerCount     int                       `json:"totalPeerCount"`
	Plugins            []EnterprisePluginSummary `json:"plugins"`
}

// EnterpriseCombinedEventType labels emitted enterprise special node events.
type EnterpriseCombinedEventType string

const (
	enterpriseEventAttach EnterpriseCombinedEventType = "enterprise.plugin.attach"
	enterpriseEventDetach EnterpriseCombinedEventType = "enterprise.plugin.detach"
	enterpriseEventUpdate EnterpriseCombinedEventType = "enterprise.plugin.update"
)

// EnterpriseCombinedEvent captures lifecycle changes for observability.
type EnterpriseCombinedEvent struct {
	Sequence  uint64                      `json:"sequence"`
	Timestamp time.Time                   `json:"timestamp"`
	Type      EnterpriseCombinedEventType `json:"type"`
	PluginID  string                      `json:"pluginId"`
	Details   map[string]string           `json:"details,omitempty"`
}

const defaultEnterpriseEventLimit = 256

// EnterpriseSpecialNode allows operators to combine several node types under a
// single enterprise control surface. Plugins provide metrics, transaction
// broadcast hooks and ledger accessors so the node can collate data without
// dictating concrete implementations.
type EnterpriseSpecialNode struct {
	*BaseNode

	mu       sync.RWMutex
	plugins  map[string]EnterpriseNodePlugin
	events   []EnterpriseCombinedEvent
	eventSeq uint64
	evLimit  int

	watchers  map[uint64]chan EnterpriseCombinedEvent
	watcherID uint64
}

var enterpriseSpecialGasSchedule = map[string]uint64{
	"EnterpriseSpecialAttach":    110,
	"EnterpriseSpecialDetach":    55,
	"EnterpriseSpecialBroadcast": 145,
	"EnterpriseSpecialSnapshot":  40,
	"EnterpriseSpecialLedger":    30,
}

// NewEnterpriseSpecialNode constructs a combined enterprise node and ensures the
// Stage 79 gas schedule is synchronised.
func NewEnterpriseSpecialNode(id nodes.Address) *EnterpriseSpecialNode {
	if _, err := synn.EnsureGasSchedule(enterpriseSpecialGasSchedule); err != nil {
		// The schedule must be present during Stage 79 rollouts. Bubble the
		// panic so CI catches mismatches between code and documentation.
		panic(fmt.Sprintf("ensure enterprise special gas schedule: %v", err))
	}

	return &EnterpriseSpecialNode{
		BaseNode: NewBaseNode(id),
		plugins:  make(map[string]EnterpriseNodePlugin),
		evLimit:  defaultEnterpriseEventLimit,
		watchers: make(map[uint64]chan EnterpriseCombinedEvent),
	}
}

// AttachPlugin registers a plugin with the enterprise node.
func (e *EnterpriseSpecialNode) AttachPlugin(plugin EnterpriseNodePlugin) error {
	if plugin.ID == "" {
		return errors.New("plugin id required")
	}
	if plugin.Role == "" {
		return errors.New("plugin role required")
	}
	if _, err := ParseCombinedNodeRole(string(plugin.Role)); err != nil {
		return err
	}

	copyLabels := func(src map[string]string) map[string]string {
		if len(src) == 0 {
			return nil
		}
		out := make(map[string]string, len(src))
		for k, v := range src {
			out[k] = v
		}
		return out
	}

	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.plugins[plugin.ID]; exists {
		return fmt.Errorf("plugin %s already registered", plugin.ID)
	}
	plugin.Labels = copyLabels(plugin.Labels)
	e.plugins[plugin.ID] = plugin
	e.recordEventLocked(enterpriseEventAttach, plugin.ID, map[string]string{"role": string(plugin.Role)})
	return nil
}

// DetachPlugin removes a plugin by identifier. It returns true when a plugin was
// removed so callers can detect stale identifiers.
func (e *EnterpriseSpecialNode) DetachPlugin(id string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.plugins[id]; !ok {
		return false
	}
	delete(e.plugins, id)
	e.recordEventLocked(enterpriseEventDetach, id, nil)
	return true
}

// UpdatePluginLabels replaces the metadata attached to a plugin. An error is
// returned when the plugin is not registered.
func (e *EnterpriseSpecialNode) UpdatePluginLabels(id string, labels map[string]string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	plugin, ok := e.plugins[id]
	if !ok {
		return fmt.Errorf("plugin %s not registered", id)
	}
	plugin.Labels = nil
	if len(labels) > 0 {
		plugin.Labels = make(map[string]string, len(labels))
		for k, v := range labels {
			plugin.Labels[k] = v
		}
	}
	e.plugins[id] = plugin
	e.recordEventLocked(enterpriseEventUpdate, id, map[string]string{"labels": fmt.Sprintf("%d", len(labels))})
	return nil
}

// Plugins returns a copy of registered plugins for inspection.
func (e *EnterpriseSpecialNode) Plugins() []EnterpriseNodePlugin {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]EnterpriseNodePlugin, 0, len(e.plugins))
	for _, p := range e.plugins {
		dup := p
		if len(p.Labels) > 0 {
			dup.Labels = make(map[string]string, len(p.Labels))
			for k, v := range p.Labels {
				dup.Labels[k] = v
			}
		}
		out = append(out, dup)
	}
	return out
}

// Snapshot aggregates plugin metrics for dashboards.
func (e *EnterpriseSpecialNode) Snapshot() EnterpriseCombinedSnapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()

	snapshot := EnterpriseCombinedSnapshot{
		Timestamp: time.Now().UTC(),
		Roles:     make(map[CombinedNodeRole]int),
	}

	for _, plugin := range e.plugins {
		summary := EnterprisePluginSummary{ID: plugin.ID, Role: plugin.Role}
		if len(plugin.Labels) > 0 {
			summary.Labels = make(map[string]string, len(plugin.Labels))
			for k, v := range plugin.Labels {
				summary.Labels[k] = v
			}
		}
		if plugin.Metrics != nil {
			metrics := plugin.Metrics()
			summary.Metrics = metrics
			snapshot.TotalMempool += metrics.MempoolSize
			snapshot.TotalValidators += metrics.ValidatorCount
			snapshot.TotalPeerCount += metrics.PeerCount
			if metrics.BlockHeight > snapshot.HighestBlockHeight {
				snapshot.HighestBlockHeight = metrics.BlockHeight
			}
		}
		snapshot.Plugins = append(snapshot.Plugins, summary)
		snapshot.NodeCount++
		snapshot.Roles[plugin.Role]++
	}
	return snapshot
}

// BroadcastTransaction invokes the broadcast hook for every plugin. Errors are
// collected and keyed by plugin identifier.
func (e *EnterpriseSpecialNode) BroadcastTransaction(tx *Transaction) (map[string]error, error) {
	if tx == nil {
		return nil, errors.New("transaction required")
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	results := make(map[string]error, len(e.plugins))
	for id, plugin := range e.plugins {
		if plugin.Broadcast == nil {
			continue
		}
		clone := cloneTransaction(tx)
		results[id] = plugin.Broadcast(clone)
	}
	return results, nil
}

// LedgerBalance sums the balances of the provided address across all plugins.
func (e *EnterpriseSpecialNode) LedgerBalance(addr string) uint64 {
	if addr == "" {
		return 0
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	var total uint64
	for _, plugin := range e.plugins {
		if plugin.LedgerBalance == nil {
			continue
		}
		total += plugin.LedgerBalance(addr)
	}
	return total
}

// Events returns a copy of the recent event log.
func (e *EnterpriseSpecialNode) Events() []EnterpriseCombinedEvent {
	e.mu.RLock()
	defer e.mu.RUnlock()
	events := make([]EnterpriseCombinedEvent, len(e.events))
	copy(events, e.events)
	return events
}

// WatchEvents registers a subscriber for lifecycle events. The returned channel
// is closed automatically when the watcher is removed.
func (e *EnterpriseSpecialNode) WatchEvents() (uint64, <-chan EnterpriseCombinedEvent) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.watcherID++
	ch := make(chan EnterpriseCombinedEvent, 16)
	e.watchers[e.watcherID] = ch
	for _, event := range e.events {
		select {
		case ch <- event:
		default:
		}
	}
	return e.watcherID, ch
}

// StopWatching removes a previously registered watcher.
func (e *EnterpriseSpecialNode) StopWatching(id uint64) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ch, ok := e.watchers[id]; ok {
		delete(e.watchers, id)
		close(ch)
	}
}

// SetEventLimit adjusts the retained event history.
func (e *EnterpriseSpecialNode) SetEventLimit(limit int) {
	if limit <= 0 {
		limit = defaultEnterpriseEventLimit
	}
	e.mu.Lock()
	e.evLimit = limit
	if len(e.events) > limit {
		e.events = append([]EnterpriseCombinedEvent(nil), e.events[len(e.events)-limit:]...)
	}
	e.mu.Unlock()
}

// recordEventLocked appends an event to the log. The caller must hold e.mu.
func (e *EnterpriseSpecialNode) recordEventLocked(eventType EnterpriseCombinedEventType, pluginID string, details map[string]string) {
	e.eventSeq++
	evt := EnterpriseCombinedEvent{
		Sequence:  e.eventSeq,
		Timestamp: time.Now().UTC(),
		Type:      eventType,
		PluginID:  pluginID,
	}
	if len(details) > 0 {
		evt.Details = make(map[string]string, len(details))
		for k, v := range details {
			evt.Details[k] = v
		}
	}
	e.events = append(e.events, evt)
	if len(e.events) > e.evLimit {
		e.events = append([]EnterpriseCombinedEvent(nil), e.events[len(e.events)-e.evLimit:]...)
	}
	for id, ch := range e.watchers {
		select {
		case ch <- evt:
		default:
			// Drop events on back pressure to keep the node responsive.
			close(ch)
			delete(e.watchers, id)
		}
	}
}

// cloneTransaction creates a deep copy of the transaction for plugin isolation.
func cloneTransaction(tx *Transaction) *Transaction {
	if tx == nil {
		return nil
	}
	cp := *tx
	if len(tx.Program) > 0 {
		cp.Program = append([]Instruction(nil), tx.Program...)
	}
	if len(tx.Signature) > 0 {
		cp.Signature = append([]byte(nil), tx.Signature...)
	}
	if len(tx.BiometricHash) > 0 {
		cp.BiometricHash = append([]byte(nil), tx.BiometricHash...)
	}
	return &cp
}

// EnterpriseNodePluginFromNode constructs an EnterpriseNodePlugin from a
// standard core.Node instance.
func EnterpriseNodePluginFromNode(id string, role CombinedNodeRole, node *Node, labels map[string]string) EnterpriseNodePlugin {
	plugin := EnterpriseNodePlugin{ID: id, Role: role}
	if len(labels) > 0 {
		plugin.Labels = make(map[string]string, len(labels))
		for k, v := range labels {
			plugin.Labels[k] = v
		}
	}
	if node != nil {
		plugin.Metrics = func() CombinedNodeMetrics {
			node.mu.Lock()
			mempool := len(node.Mempool)
			height := len(node.Blockchain)
			node.mu.Unlock()
			validators := 0
			if node.Validators != nil {
				validators = len(node.Validators.Eligible())
			}
			return CombinedNodeMetrics{
				MempoolSize:    mempool,
				ValidatorCount: validators,
				BlockHeight:    height,
			}
		}
		plugin.Broadcast = func(tx *Transaction) error {
			if tx == nil {
				return errors.New("transaction required")
			}
			return node.AddTransaction(tx)
		}
		plugin.LedgerBalance = func(addr string) uint64 {
			if node.Ledger == nil {
				return 0
			}
			return node.Ledger.GetBalance(addr)
		}
	}
	return plugin
}

// EnterpriseSpecialAttach is an exported helper used by VM opcode bindings and
// documentation tooling to register a plugin with the provided node.
func EnterpriseSpecialAttach(node *EnterpriseSpecialNode, plugin EnterpriseNodePlugin) error {
	if node == nil {
		return errors.New("enterprise special node required")
	}
	return node.AttachPlugin(plugin)
}

// EnterpriseSpecialDetach removes a plugin from the provided node.
func EnterpriseSpecialDetach(node *EnterpriseSpecialNode, id string) bool {
	if node == nil {
		return false
	}
	return node.DetachPlugin(id)
}

// EnterpriseSpecialBroadcast proxies to BroadcastTransaction for opcode usage.
func EnterpriseSpecialBroadcast(node *EnterpriseSpecialNode, tx *Transaction) (map[string]error, error) {
	if node == nil {
		return nil, errors.New("enterprise special node required")
	}
	return node.BroadcastTransaction(tx)
}

// EnterpriseSpecialSnapshot returns the aggregated metrics for the supplied node.
func EnterpriseSpecialSnapshot(node *EnterpriseSpecialNode) EnterpriseCombinedSnapshot {
	if node == nil {
		return EnterpriseCombinedSnapshot{Timestamp: time.Now().UTC()}
	}
	return node.Snapshot()
}

// EnterpriseSpecialLedger returns the aggregated ledger balance for the address.
func EnterpriseSpecialLedger(node *EnterpriseSpecialNode, addr string) uint64 {
	if node == nil {
		return 0
	}
	return node.LedgerBalance(addr)
}
