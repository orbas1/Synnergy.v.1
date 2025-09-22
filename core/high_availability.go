package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	security "synnergy/internal/security"
	"synnergy/internal/telemetry"
)

// FailoverOption configures the failover manager at construction time.
type FailoverOption func(*FailoverManager)

// FailoverNode captures metadata for a node participating in the failover ring.
type FailoverNode struct {
	ID                string
	Role              string
	Region            string
	PublicKey         []byte
	LastHeartbeat     time.Time
	LastLatency       time.Duration
	SignatureVerified bool
	LastSignature     []byte
	LastPayload       []byte
	Healthy           bool
}

// FailoverAudit records lifecycle events for compliance and troubleshooting.
type FailoverAudit struct {
	Timestamp time.Time `json:"timestamp"`
	NodeID    string    `json:"nodeId"`
	Event     string    `json:"event"`
	Details   string    `json:"details,omitempty"`
}

// HeartbeatProof captures signed heartbeat data for verification.
type HeartbeatProof struct {
	ID        string
	Payload   []byte
	Signature []byte
	Latency   time.Duration
}

// FailoverNodeStatus is a read-only view exported via diagnostics.
type FailoverNodeStatus struct {
	ID                string    `json:"id"`
	Role              string    `json:"role"`
	Region            string    `json:"region"`
	Healthy           bool      `json:"healthy"`
	LastHeartbeat     time.Time `json:"lastHeartbeat"`
	LatencyMS         int64     `json:"latencyMs"`
	SignatureVerified bool      `json:"signatureVerified"`
}

// ResilienceReport summarises Stage 77 high-availability health.
type ResilienceReport struct {
	Timestamp         time.Time            `json:"timestamp"`
	ActiveNode        string               `json:"activeNode"`
	ActiveRole        string               `json:"activeRole"`
	ActiveRegion      string               `json:"activeRegion"`
	Backups           []FailoverNodeStatus `json:"backups"`
	ConsensusNetworks int                  `json:"consensusNetworks"`
	AuthorityNodes    int                  `json:"authorityNodes"`
	WalletReady       bool                 `json:"walletReady"`
	WalletAddress     string               `json:"walletAddress"`
	VMHealthy         bool                 `json:"vmHealthy"`
	VMMode            string               `json:"vmMode"`
	VMConcurrency     int                  `json:"vmConcurrency"`
	LedgerHeight      int                  `json:"ledgerHeight"`
	RegionDiversity   bool                 `json:"regionDiversity"`
	ScalabilityScore  float64              `json:"scalabilityScore"`
	GovernanceReady   bool                 `json:"governanceReady"`
	Interoperability  bool                 `json:"interoperability"`
	PrivacyTier       string               `json:"privacyTier"`
	Compliance        []string             `json:"compliance"`
	LastFailover      time.Time            `json:"lastFailover,omitempty"`
	AuditTrail        []FailoverAudit      `json:"auditTrail,omitempty"`
	Signature         string               `json:"signature,omitempty"`
	SigningKey        string               `json:"signingKey,omitempty"`
	SigningKeyVersion int                  `json:"signingKeyVersion,omitempty"`
}

// FailoverManager tracks node heartbeats to provide high availability through
// automatic promotion of backup nodes when the primary becomes unresponsive.
type FailoverManager struct {
	mu           sync.RWMutex
	primary      string
	nodes        map[string]*FailoverNode
	timeout      time.Duration
	lastFailover time.Time
	signer       *security.KeyManager
	vm           *SimpleVM
	consensus    *ConsensusNetworkManager
	wallet       *Wallet
	ledger       *Ledger
	registry     *AuthorityNodeRegistry
	audit        []FailoverAudit
}

// WithFailoverSigner provides a custom signer for report attestation.
func WithFailoverSigner(km *security.KeyManager) FailoverOption {
	return func(m *FailoverManager) {
		if km != nil {
			m.signer = km
		}
	}
}

// WithFailoverVirtualMachine integrates VM status into diagnostics.
func WithFailoverVirtualMachine(vm *SimpleVM) FailoverOption {
	return func(m *FailoverManager) {
		m.vm = vm
	}
}

// WithFailoverConsensus integrates consensus topology awareness.
func WithFailoverConsensus(consensus *ConsensusNetworkManager) FailoverOption {
	return func(m *FailoverManager) {
		m.consensus = consensus
	}
}

// WithFailoverWallet integrates wallet readiness into diagnostics.
func WithFailoverWallet(wallet *Wallet) FailoverOption {
	return func(m *FailoverManager) {
		m.wallet = wallet
	}
}

// WithFailoverLedger integrates ledger height into diagnostics.
func WithFailoverLedger(ledger *Ledger) FailoverOption {
	return func(m *FailoverManager) {
		m.ledger = ledger
	}
}

// WithFailoverRegistry integrates authority registry coverage.
func WithFailoverRegistry(registry *AuthorityNodeRegistry) FailoverOption {
	return func(m *FailoverManager) {
		m.registry = registry
	}
}

// NewFailoverManager creates a FailoverManager with a primary node identifier
// and a timeout indicating how long a node may miss heartbeats before being
// considered offline.
func NewFailoverManager(primary string, timeout time.Duration, opts ...FailoverOption) *FailoverManager {
	if primary == "" {
		primary = "primary"
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	fm := &FailoverManager{
		primary: primary,
		nodes:   make(map[string]*FailoverNode),
		timeout: timeout,
		signer:  security.NewKeyManager(),
	}
	for _, opt := range opts {
		opt(fm)
	}
	if fm.signer == nil {
		fm.signer = security.NewKeyManager()
	}
	if fm.signer != nil {
		if _, _, _, err := fm.signer.SigningKey(security.PurposeStateSigning); err != nil {
			if _, _, err := fm.signer.GenerateSigningKey(security.PurposeStateSigning, "failover-manager"); err != nil {
				fallback := security.NewKeyManager()
				if _, _, genErr := fallback.GenerateSigningKey(security.PurposeStateSigning, "failover-manager"); genErr == nil {
					fm.signer = fallback
				} else {
					fm.signer = nil
				}
			}
		}
	}
	now := time.Now().UTC()
	fm.nodes[primary] = &FailoverNode{
		ID:            primary,
		Role:          "primary",
		Region:        "global",
		LastHeartbeat: now,
		Healthy:       true,
	}
	fm.audit = append(fm.audit, FailoverAudit{
		Timestamp: now,
		NodeID:    primary,
		Event:     "bootstrap",
		Details:   "initial primary registered",
	})
	return fm
}

// RegisterBackup adds a new backup node to the manager.
func (m *FailoverManager) RegisterBackup(id string) {
	m.RegisterNode(FailoverNode{ID: id, Role: "validator", Region: "global"})
}

// RegisterNode registers or updates metadata for a failover participant.
func (m *FailoverManager) RegisterNode(node FailoverNode) {
	if node.ID == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	existing := m.nodes[node.ID]
	if node.Role == "" {
		if existing != nil && existing.Role != "" {
			node.Role = existing.Role
		} else {
			node.Role = "validator"
		}
	}
	if node.Region == "" {
		if existing != nil && existing.Region != "" {
			node.Region = existing.Region
		} else {
			node.Region = "global"
		}
	}
	now := time.Now().UTC()
	record := &FailoverNode{
		ID:            node.ID,
		Role:          node.Role,
		Region:        node.Region,
		PublicKey:     append([]byte(nil), node.PublicKey...),
		LastHeartbeat: now,
		Healthy:       true,
	}
	if existing != nil {
		if !existing.LastHeartbeat.IsZero() {
			record.LastHeartbeat = existing.LastHeartbeat
		}
		record.LastLatency = existing.LastLatency
		if len(record.PublicKey) == 0 {
			record.PublicKey = append([]byte(nil), existing.PublicKey...)
		}
		record.LastPayload = append([]byte(nil), existing.LastPayload...)
		record.LastSignature = append([]byte(nil), existing.LastSignature...)
		record.SignatureVerified = existing.SignatureVerified
		record.Healthy = existing.Healthy
	}
	m.nodes[node.ID] = record
	m.audit = append(m.audit, FailoverAudit{
		Timestamp: now,
		NodeID:    node.ID,
		Event:     "register",
		Details:   fmt.Sprintf("role=%s region=%s", record.Role, record.Region),
	})
}

// Heartbeat records a heartbeat for the specified node.
func (m *FailoverManager) Heartbeat(id string) {
	_ = m.RecordHeartbeat(HeartbeatProof{ID: id})
}

// RecordHeartbeat records a heartbeat together with optional latency and proof data.
func (m *FailoverManager) RecordHeartbeat(proof HeartbeatProof) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	node, ok := m.nodes[proof.ID]
	if !ok {
		return fmt.Errorf("failover: node %s not registered", proof.ID)
	}
	now := time.Now().UTC()
	node.LastHeartbeat = now
	if proof.Latency > 0 {
		node.LastLatency = proof.Latency
	}
	if len(proof.Payload) > 0 {
		node.LastPayload = append([]byte(nil), proof.Payload...)
	}
	if len(proof.Signature) > 0 {
		node.LastSignature = append([]byte(nil), proof.Signature...)
		verified := false
		if len(node.PublicKey) > 0 {
			if pub, err := decodePublicKey(node.PublicKey); err == nil {
				verified = VerifyMessage(proof.Payload, proof.Signature, pub)
			}
		}
		node.SignatureVerified = verified
	}
	node.Healthy = true
	m.nodes[proof.ID] = node
	m.audit = append(m.audit, FailoverAudit{
		Timestamp: now,
		NodeID:    proof.ID,
		Event:     "heartbeat",
		Details:   fmt.Sprintf("latency=%s verified=%t", proof.Latency, node.SignatureVerified),
	})
	return nil
}

// RemoveNode removes a node from consideration. If the primary is removed the
// next call to Active will promote the freshest backup.
func (m *FailoverManager) RemoveNode(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.nodes, id)
	if id == m.primary {
		m.primary = ""
	}
	m.audit = append(m.audit, FailoverAudit{
		Timestamp: time.Now().UTC(),
		NodeID:    id,
		Event:     "remove",
	})
}

// Active returns the identifier of the node currently acting as primary.  If the
// existing primary has not sent a heartbeat within the timeout, the most recent
// backup node is promoted.
func (m *FailoverManager) Active() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	if node, ok := m.nodes[m.primary]; ok {
		if now.Sub(node.LastHeartbeat) <= m.timeout {
			node.Healthy = true
			return m.primary
		}
		node.Healthy = false
	}

	var candidate string
	var latest time.Time
	for id, node := range m.nodes {
		if id == m.primary {
			continue
		}
		healthy := now.Sub(node.LastHeartbeat) <= m.timeout
		node.Healthy = healthy
		if !healthy {
			continue
		}
		if node.LastHeartbeat.After(latest) {
			candidate = id
			latest = node.LastHeartbeat
		}
	}
	if candidate != "" && candidate != m.primary {
		m.lastFailover = now.UTC()
		m.audit = append(m.audit, FailoverAudit{
			Timestamp: m.lastFailover,
			NodeID:    candidate,
			Event:     "promote",
			Details:   "failover",
		})
		m.primary = candidate
	}
	return m.primary
}

// Report produces a resilience snapshot suitable for CLI and web dashboards.
func (m *FailoverManager) Report(ctx context.Context) ResilienceReport {
	activeID := m.Active()
	_, span := telemetry.Tracer("core.failover").Start(ctx, "FailoverManager.Report")
	defer span.End()

	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now().UTC()
	report := ResilienceReport{
		Timestamp:     now,
		WalletReady:   m.wallet != nil && m.wallet.Address != "",
		WalletAddress: "",
	}
	if m.wallet != nil {
		report.WalletAddress = m.wallet.Address
	}
	if m.vm != nil {
		report.VMHealthy = m.vm.Status()
		report.VMMode = m.vm.Mode().String()
		report.VMConcurrency = m.vm.Concurrency()
	}
	if m.consensus != nil {
		report.ConsensusNetworks = len(m.consensus.ListNetworks())
	}
	if m.registry != nil {
		report.AuthorityNodes = len(m.registry.List())
	}
	if m.ledger != nil {
		height, _ := m.ledger.Head()
		report.LedgerHeight = height
	}
	report.Interoperability = report.ConsensusNetworks > 0
	report.GovernanceReady = report.AuthorityNodes >= 3

	var activeStatus FailoverNodeStatus
	backups := make([]FailoverNodeStatus, 0, len(m.nodes))
	regionSet := make(map[string]struct{})
	complianceSet := make(map[string]struct{})
	healthyCount := 0
	var latencyTotal time.Duration

	for id, node := range m.nodes {
		healthy := node.Healthy && time.Since(node.LastHeartbeat) <= m.timeout
		verified := node.SignatureVerified && len(node.PublicKey) > 0
		status := FailoverNodeStatus{
			ID:                id,
			Role:              node.Role,
			Region:            node.Region,
			Healthy:           healthy,
			LastHeartbeat:     node.LastHeartbeat,
			LatencyMS:         node.LastLatency.Milliseconds(),
			SignatureVerified: verified,
		}
		if status.Healthy {
			healthyCount++
			latencyTotal += node.LastLatency
		}
		if status.Region != "" {
			regionSet[status.Region] = struct{}{}
		}
		if !status.SignatureVerified {
			complianceSet["unverified-heartbeats"] = struct{}{}
		}
		if id == activeID {
			activeStatus = status
			continue
		}
		backups = append(backups, status)
	}
	sort.Slice(backups, func(i, j int) bool { return backups[i].ID < backups[j].ID })
	report.Backups = backups
	report.ActiveNode = activeStatus.ID
	report.ActiveRole = activeStatus.Role
	report.ActiveRegion = activeStatus.Region
	report.RegionDiversity = len(regionSet) > 1
	if !report.RegionDiversity {
		complianceSet["single-region"] = struct{}{}
	}
	if !report.VMHealthy {
		complianceSet["vm-offline"] = struct{}{}
	}
	if !report.WalletReady {
		complianceSet["wallet-not-ready"] = struct{}{}
	}
	if !report.Interoperability {
		complianceSet["no-cross-consensus-links"] = struct{}{}
	}
	if !report.GovernanceReady {
		complianceSet["insufficient-authority-quorum"] = struct{}{}
	}

	if healthyCount == 0 {
		report.ScalabilityScore = 0
	} else {
		avgLatency := latencyTotal
		if avgLatency <= 0 {
			avgLatency = 1 * time.Millisecond
		} else {
			avgLatency /= time.Duration(healthyCount)
		}
		ms := avgLatency.Seconds() * 1000
		if ms < 1 {
			ms = 1
		}
		report.ScalabilityScore = float64(healthyCount) * (1000.0 / ms)
	}

	if report.GovernanceReady && report.RegionDiversity {
		report.PrivacyTier = "permissioned"
	} else if report.GovernanceReady {
		report.PrivacyTier = "hybrid"
	} else {
		report.PrivacyTier = "open"
	}

	if m.lastFailover.Unix() > 0 {
		report.LastFailover = m.lastFailover
	}

	if len(m.audit) > 0 {
		limit := len(m.audit)
		if limit > 10 {
			limit = 10
		}
		start := len(m.audit) - limit
		report.AuditTrail = append([]FailoverAudit(nil), m.audit[start:]...)
	}

	if len(complianceSet) > 0 {
		report.Compliance = make([]string, 0, len(complianceSet))
		for key := range complianceSet {
			report.Compliance = append(report.Compliance, key)
		}
		sort.Strings(report.Compliance)
	}

	if m.signer != nil {
		unsigned := report
		unsigned.Signature = ""
		unsigned.SigningKey = ""
		unsigned.SigningKeyVersion = 0
		payload, err := json.Marshal(unsigned)
		if err == nil {
			sig, pub, version, err := m.signer.Sign(security.PurposeStateSigning, payload)
			if err == nil {
				report.Signature = base64.StdEncoding.EncodeToString(sig)
				report.SigningKey = base64.StdEncoding.EncodeToString(pub)
				report.SigningKeyVersion = version
			} else {
				report.Compliance = append(report.Compliance, "signature-error")
			}
		}
	}

	return report
}
