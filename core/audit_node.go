package core

import "errors"

// BootstrapNode defines the minimal behaviour required from a bootstrap
// node. The concrete implementation is expected to handle network
// initialisation for peers.
type BootstrapNode interface {
	Start() error
}

// AuditNode ties a BootstrapNode with the AuditManager to provide
// on-chain audit logging capabilities.
type AuditNode struct {
	Bootstrap BootstrapNode
	Manager   *AuditManager
}

// NewAuditNode constructs a new AuditNode instance.
func NewAuditNode(b BootstrapNode, m *AuditManager) *AuditNode {
	return &AuditNode{Bootstrap: b, Manager: m}
}

// Start launches the underlying bootstrap node.
func (n *AuditNode) Start() error {
	if n.Bootstrap == nil {
		return errors.New("bootstrap node not configured")
	}
	return n.Bootstrap.Start()
}

// LogEvent records an audit event through the underlying manager.
func (n *AuditNode) LogEvent(address, event string, metadata map[string]string) error {
	if n.Manager == nil {
		return errors.New("audit manager not configured")
	}
	return n.Manager.Log(address, event, metadata)
}

// ListEvents retrieves audit events for an address.
func (n *AuditNode) ListEvents(address string) []AuditEntry {
	if n.Manager == nil {
		return nil
	}
	return n.Manager.List(address)
}
