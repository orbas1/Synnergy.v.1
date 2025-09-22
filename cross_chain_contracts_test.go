package synnergy

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type stubConnectionSource struct {
	mu          sync.RWMutex
	connections map[string]*ChainConnection
}

func newStubConnectionSource() *stubConnectionSource {
	return &stubConnectionSource{connections: make(map[string]*ChainConnection)}
}

func (s *stubConnectionSource) GetConnection(id string) (*ChainConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conn, ok := s.connections[id]
	if !ok {
		return nil, false
	}
	return cloneConnection(conn), true
}

func (s *stubConnectionSource) setConnection(conn *ChainConnection) {
	s.mu.Lock()
	s.connections[conn.ID] = conn
	s.mu.Unlock()
}

type deterministicClock struct {
	mu      sync.Mutex
	current time.Time
	step    time.Duration
}

func newDeterministicClock(start time.Time) *deterministicClock {
	return &deterministicClock{current: start, step: time.Second}
}

func (c *deterministicClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.current
	c.current = c.current.Add(c.step)
	return now
}

func newTestContractLinkManager(t *testing.T) (*ContractLinkManager, *deterministicClock, *stubConnectionSource) {
	t.Helper()
	clock := newDeterministicClock(time.Unix(1_700_000_000, 0))
	resolver := newStubConnectionSource()
	resolver.setConnection(&ChainConnection{
		ID:     "conn-1",
		Status: ConnectionStatusActive,
		Spec: ConnectionSpec{
			LocalChain:     "synnergy-main",
			RemoteChain:    "ally-main",
			LocalEndpoint:  "https://synnergy.local",
			RemoteEndpoint: "https://ally.local",
			GasProfile:     "syn-default",
		},
		CreatedAt: clock.Now(),
		UpdatedAt: clock.Now(),
		OpenedAt:  clock.Now(),
	})
	manager := NewContractLinkManager(
		WithContractLinkResolver(resolver),
		WithContractLinkClock(clock.Now),
	)
	return manager, clock, resolver
}

func baseContractLinkSpec() ContractLinkSpec {
	return ContractLinkSpec{
		LocalChain:    "synnergy-main",
		LocalAddress:  "0xabc",
		RemoteChain:   "ally-main",
		RemoteAddress: "0xdef",
		ConnectionID:  "conn-1",
		Capabilities:  []string{"invoke", "query"},
		GasProfile:    "syn-default",
		Metadata: map[string]string{
			"department": "treasury",
		},
		AccessPolicy: AccessPolicy{
			AllowedApprovers:  []string{"alice", "bob"},
			RequiredApprovals: 0,
			PrivacyLevel:      "confidential",
			EncryptionScheme:  "aes-gcm-256",
		},
		AuditTrailHint: "ops",
	}
}

func TestContractLinkRegisterValidation(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	spec := baseContractLinkSpec()
	spec.LocalAddress = ""
	if _, err := manager.Register(spec); !errors.Is(err, ErrContractLinkInvalidSpec) {
		t.Fatalf("expected ErrContractLinkInvalidSpec, got %v", err)
	}

	spec = baseContractLinkSpec()
	spec.AccessPolicy.RequiredApprovals = 1
	spec.AccessPolicy.AllowedApprovers = nil
	if _, err := manager.Register(spec); !errors.Is(err, ErrContractLinkInvalidSpec) {
		t.Fatalf("expected ErrContractLinkInvalidSpec for missing approvers, got %v", err)
	}
}

func TestContractLinkRegisterRequiresActiveConnection(t *testing.T) {
	clock := newDeterministicClock(time.Unix(1_700_010_000, 0))
	resolver := newStubConnectionSource()
	resolver.setConnection(&ChainConnection{ID: "conn-1", Status: ConnectionStatusPending})
	manager := NewContractLinkManager(WithContractLinkResolver(resolver), WithContractLinkClock(clock.Now))
	if _, err := manager.Register(baseContractLinkSpec()); !errors.Is(err, ErrContractLinkConnectionInactive) {
		t.Fatalf("expected ErrContractLinkConnectionInactive, got %v", err)
	}

	resolverEmpty := newStubConnectionSource()
	manager = NewContractLinkManager(WithContractLinkResolver(resolverEmpty), WithContractLinkClock(clock.Now))
	if _, err := manager.Register(baseContractLinkSpec()); !errors.Is(err, ErrContractLinkConnectionUnknown) {
		t.Fatalf("expected ErrContractLinkConnectionUnknown, got %v", err)
	}
}

func TestContractLinkRegisterAutoActive(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	link, err := manager.Register(baseContractLinkSpec())
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if link.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status, got %s", link.Status)
	}
	if link.ActivatedAt.IsZero() {
		t.Fatalf("expected activation timestamp to be set")
	}
	if link.Version != 1 {
		t.Fatalf("expected version 1, got %d", link.Version)
	}
}

func TestContractLinkApprovalFlow(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	spec := baseContractLinkSpec()
	spec.LocalAddress = "0xaaa"
	spec.AccessPolicy.RequiredApprovals = 2
	spec.AccessPolicy.AllowedApprovers = []string{"alice", "bob", "carol"}
	link, err := manager.Register(spec)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if link.Status != ContractLinkStatusPending {
		t.Fatalf("expected pending status, got %s", link.Status)
	}
	if link.Version != 1 {
		t.Fatalf("expected version 1 after registration, got %d", link.Version)
	}

	updated, err := manager.RecordApproval(link.ID, "alice")
	if err != nil {
		t.Fatalf("record approval failed: %v", err)
	}
	if updated.Status != ContractLinkStatusPending {
		t.Fatalf("expected pending after first approval, got %s", updated.Status)
	}
	if updated.Version != 2 {
		t.Fatalf("expected version 2 after approval, got %d", updated.Version)
	}

	if _, err = manager.RecordApproval(link.ID, "ALICE"); !errors.Is(err, ErrContractLinkApprovalDuplicate) {
		t.Fatalf("expected duplicate approval error, got %v", err)
	}

	activated, err := manager.RecordApproval(link.ID, "bob")
	if err != nil {
		t.Fatalf("second approval failed: %v", err)
	}
	if activated.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status after threshold, got %s", activated.Status)
	}
	if len(activated.Approvals) != 2 {
		t.Fatalf("expected two approvals, got %d", len(activated.Approvals))
	}
	if activated.Version != 3 {
		t.Fatalf("expected version 3 after activation, got %d", activated.Version)
	}
}

func TestContractLinkSuspendResumeAndFailure(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	link, err := manager.Register(baseContractLinkSpec())
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	failed, err := manager.ReportFailure(link.ID, "SYNC_TIMEOUT", "remote acknowledgement missing")
	if err != nil {
		t.Fatalf("report failure failed: %v", err)
	}
	if failed.Status != ContractLinkStatusFailed {
		t.Fatalf("expected failed status, got %s", failed.Status)
	}
	if len(failed.Failures) != 1 {
		t.Fatalf("expected one failure record, got %d", len(failed.Failures))
	}

	resumed, err := manager.Resume(link.ID)
	if err != nil {
		t.Fatalf("resume failed: %v", err)
	}
	if resumed.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status after resume, got %s", resumed.Status)
	}
	last := resumed.Failures[len(resumed.Failures)-1]
	if !last.Resolved {
		t.Fatalf("expected failure to be resolved after resume")
	}

	suspended, err := manager.Suspend(link.ID, "maintenance window")
	if err != nil {
		t.Fatalf("suspend failed: %v", err)
	}
	if suspended.Status != ContractLinkStatusSuspended {
		t.Fatalf("expected suspended status, got %s", suspended.Status)
	}

	resumed, err = manager.Resume(link.ID)
	if err != nil {
		t.Fatalf("resume after suspension failed: %v", err)
	}
	if resumed.Status != ContractLinkStatusActive {
		t.Fatalf("expected active status after resume, got %s", resumed.Status)
	}

	if _, err := manager.Suspend(link.ID, ""); err == nil {
		t.Fatalf("expected error when suspending without reason")
	}
}

func TestContractLinkListFiltering(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	first, err := manager.Register(baseContractLinkSpec())
	if err != nil {
		t.Fatalf("register first failed: %v", err)
	}
	spec := baseContractLinkSpec()
	spec.LocalAddress = "0x999"
	spec.RemoteChain = "ally-analytics"
	spec.RemoteAddress = "0xeee"
	spec.AccessPolicy.RequiredApprovals = 1
	spec.AccessPolicy.AllowedApprovers = []string{"dora"}
	second, err := manager.Register(spec)
	if err != nil {
		t.Fatalf("register second failed: %v", err)
	}
	if second.Status != ContractLinkStatusPending {
		t.Fatalf("expected pending status for second, got %s", second.Status)
	}

	pending := manager.List(ContractLinkFilter{Statuses: []ContractLinkStatus{ContractLinkStatusPending}})
	if len(pending) != 1 || pending[0].ID != second.ID {
		t.Fatalf("expected only pending contract, got %#v", pending)
	}

	ally := manager.List(ContractLinkFilter{RemoteChain: "ally-main"})
	if len(ally) != 1 || ally[0].ID != first.ID {
		t.Fatalf("expected filter by remote chain to return first mapping")
	}

	none := manager.List(ContractLinkFilter{ConnectionID: "missing"})
	if len(none) != 0 {
		t.Fatalf("expected no results for unknown connection, got %d", len(none))
	}
}

func TestContractLinkUpdate(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	link, err := manager.Register(baseContractLinkSpec())
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	updated, err := manager.Update(link.ID, ContractLinkUpdate{
		GasProfile:   "syn-premium",
		Metadata:     map[string]string{"tier": "gold"},
		Capabilities: []string{"invoke", "audit"},
	})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Spec.GasProfile != "syn-premium" {
		t.Fatalf("expected gas profile to change, got %s", updated.Spec.GasProfile)
	}
	if updated.Spec.Metadata["tier"] != "gold" {
		t.Fatalf("expected metadata to include tier, got %#v", updated.Spec.Metadata)
	}
	if updated.Version != 2 {
		t.Fatalf("expected version 2 after update, got %d", updated.Version)
	}

	same, err := manager.Update(link.ID, ContractLinkUpdate{})
	if err != nil {
		t.Fatalf("second update failed: %v", err)
	}
	if same.Version != updated.Version {
		t.Fatalf("expected version to remain unchanged on no-op update")
	}
}

func TestContractLinkSubscribeEvents(t *testing.T) {
	manager, _, _ := newTestContractLinkManager(t)
	ch, cancel := manager.Subscribe(4)
	defer cancel()

	spec := baseContractLinkSpec()
	spec.LocalAddress = "0x123"
	spec.AccessPolicy.RequiredApprovals = 1
	spec.AccessPolicy.AllowedApprovers = []string{"alice"}
	link, err := manager.Register(spec)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	expectEvent(t, ch, ContractLinkEventRegistered, link.ID)

	if _, err := manager.RecordApproval(link.ID, "alice"); err != nil {
		t.Fatalf("approval failed: %v", err)
	}

	expectEvent(t, ch, ContractLinkEventApprovalRecorded, link.ID)
	expectEvent(t, ch, ContractLinkEventActivated, link.ID)
}

func expectEvent(t *testing.T, ch <-chan ContractLinkEvent, eventType ContractLinkEventType, linkID string) {
	t.Helper()
	select {
	case evt := <-ch:
		if evt.Type != eventType {
			t.Fatalf("expected event %s, got %s", eventType, evt.Type)
		}
		if evt.LinkID != linkID {
			t.Fatalf("expected link id %s, got %s", linkID, evt.LinkID)
		}
		if evt.Link == nil {
			t.Fatalf("expected snapshot to be included")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("timed out waiting for event %s", eventType)
	}
}
