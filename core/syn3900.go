package core

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

// BenefitStatus represents the lifecycle of a SYN3900 benefit allocation.
type BenefitStatus string

const (
	BenefitStatusRegistered BenefitStatus = "registered"
	BenefitStatusApproved   BenefitStatus = "approved"
	BenefitStatusClaimed    BenefitStatus = "claimed"
	BenefitStatusCancelled  BenefitStatus = "cancelled"
)

// BenefitAuditEvent tracks governance actions for benefit records.
type BenefitAuditEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Actor     Address   `json:"actor"`
	Action    string    `json:"action"`
	Note      string    `json:"note,omitempty"`
}

// BenefitRecord holds metadata for a government benefit token issuance.
type BenefitRecord struct {
	ID        uint64                `json:"id"`
	Recipient string                `json:"recipient"`
	Program   string                `json:"program"`
	Amount    uint64                `json:"amount"`
	Status    BenefitStatus         `json:"status"`
	Claimed   bool                  `json:"claimed"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	Approvals map[Address]time.Time `json:"approvals"`
	ClaimedBy Address               `json:"claimed_by"`
	ClaimedAt time.Time             `json:"claimed_at"`
	Audit     []BenefitAuditEvent   `json:"audit"`
}

// BenefitRegistrySnapshot preserves registry state for Stage 73 persistence.
type BenefitRegistrySnapshot struct {
	NextID      uint64          `json:"next_id"`
	Benefits    []BenefitRecord `json:"benefits"`
	GeneratedAt time.Time       `json:"generated_at"`
}

// BenefitRegistry manages benefit records.
type BenefitRegistry struct {
	mu       sync.RWMutex
	benefits map[uint64]*BenefitRecord
	nextID   uint64
}

// NewBenefitRegistry creates a new registry.
func NewBenefitRegistry() *BenefitRegistry {
	return &BenefitRegistry{benefits: make(map[uint64]*BenefitRecord)}
}

// RegisterBenefit records a new benefit and returns its ID.
func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64, approver Address) (uint64, error) {
	if recipient == "" || program == "" {
		return 0, fmt.Errorf("recipient and program required")
	}
	if amount == 0 {
		return 0, fmt.Errorf("amount must be positive")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := r.nextID
	now := time.Now().UTC()
	record := &BenefitRecord{
		ID:        id,
		Recipient: recipient,
		Program:   program,
		Amount:    amount,
		Status:    BenefitStatusRegistered,
		CreatedAt: now,
		UpdatedAt: now,
		Approvals: make(map[Address]time.Time),
	}
	if approver != "" {
		record.Approvals[approver] = now
		record.Status = BenefitStatusApproved
	}
	record.Audit = append(record.Audit, BenefitAuditEvent{Timestamp: now, Actor: approver, Action: "register"})
	r.benefits[id] = record
	return id, nil
}

// Approve authorises the benefit for claiming.
func (r *BenefitRegistry) Approve(id uint64, actor Address) error {
	if actor == "" {
		return fmt.Errorf("approver required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Status == BenefitStatusClaimed {
		return nil
	}
	if b.Approvals == nil {
		b.Approvals = make(map[Address]time.Time)
	}
	now := time.Now().UTC()
	b.Approvals[actor] = now
	b.Status = BenefitStatusApproved
	b.UpdatedAt = now
	b.Audit = append(b.Audit, BenefitAuditEvent{Timestamp: now, Actor: actor, Action: "approve"})
	return nil
}

// Claim marks the benefit as claimed by an authorised wallet.
func (r *BenefitRegistry) Claim(id uint64, actor Address) error {
	if actor == "" {
		return fmt.Errorf("wallet required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.benefits[id]
	if !ok {
		return errors.New("benefit not found")
	}
	if b.Status != BenefitStatusApproved && b.Status != BenefitStatusRegistered {
		return fmt.Errorf("benefit not claimable")
	}
	if b.Status == BenefitStatusRegistered && len(b.Approvals) > 0 {
		// Should not happen, but guard for safety.
		b.Status = BenefitStatusApproved
	}
	if len(b.Approvals) > 0 && actor != Address(b.Recipient) {
		// If approvals exist, ensure claimer is the recipient or an approver.
		if _, ok := b.Approvals[actor]; !ok {
			return fmt.Errorf("wallet not authorised")
		}
	}
	if b.Status == BenefitStatusClaimed {
		return fmt.Errorf("benefit already claimed")
	}
	now := time.Now().UTC()
	b.Status = BenefitStatusClaimed
	b.Claimed = true
	b.ClaimedBy = actor
	b.ClaimedAt = now
	b.UpdatedAt = now
	b.Audit = append(b.Audit, BenefitAuditEvent{Timestamp: now, Actor: actor, Action: "claim"})
	return nil
}

// GetBenefit retrieves a benefit by ID.
func (r *BenefitRegistry) GetBenefit(id uint64) (*BenefitRecord, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.benefits[id]
	if !ok {
		return nil, false
	}
	cp := *b
	cp.Approvals = copyApprovals(b.Approvals)
	cp.Audit = append([]BenefitAuditEvent(nil), b.Audit...)
	cp.Claimed = b.Status == BenefitStatusClaimed
	return &cp, true
}

func copyApprovals(src map[Address]time.Time) map[Address]time.Time {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[Address]time.Time, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// ListBenefits returns all benefits.
func (r *BenefitRegistry) ListBenefits() []BenefitRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := make([]BenefitRecord, 0, len(r.benefits))
	for _, b := range r.benefits {
		cp := *b
		cp.Approvals = copyApprovals(b.Approvals)
		cp.Audit = append([]BenefitAuditEvent(nil), b.Audit...)
		cp.Claimed = b.Status == BenefitStatusClaimed
		res = append(res, cp)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// BenefitRegistryStatus summarises registry metrics.
type BenefitRegistryStatus struct {
	Total    int `json:"total"`
	Approved int `json:"approved"`
	Claimed  int `json:"claimed"`
	Pending  int `json:"pending"`
}

// StatusSummary provides aggregated metrics for dashboards.
func (r *BenefitRegistry) StatusSummary() BenefitRegistryStatus {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var status BenefitRegistryStatus
	for _, b := range r.benefits {
		status.Total++
		if len(b.Approvals) > 0 {
			status.Approved++
		}
		switch b.Status {
		case BenefitStatusClaimed:
			status.Claimed++
		case BenefitStatusApproved:
			status.Pending++
		default:
			status.Pending++
		}
	}
	return status
}

// Snapshot captures the registry state.
func (r *BenefitRegistry) Snapshot() BenefitRegistrySnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	benefits := make([]BenefitRecord, 0, len(r.benefits))
	for _, b := range r.benefits {
		cp := *b
		cp.Approvals = copyApprovals(b.Approvals)
		cp.Audit = append([]BenefitAuditEvent(nil), b.Audit...)
		cp.Claimed = b.Status == BenefitStatusClaimed
		benefits = append(benefits, cp)
	}
	sort.Slice(benefits, func(i, j int) bool { return benefits[i].ID < benefits[j].ID })
	return BenefitRegistrySnapshot{NextID: r.nextID, Benefits: benefits, GeneratedAt: time.Now().UTC()}
}

// Restore applies a snapshot to the registry.
func (r *BenefitRegistry) Restore(snapshot BenefitRegistrySnapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID = snapshot.NextID
	r.benefits = make(map[uint64]*BenefitRecord, len(snapshot.Benefits))
	for i := range snapshot.Benefits {
		benefit := snapshot.Benefits[i]
		cp := benefit
		cp.Approvals = copyApprovals(benefit.Approvals)
		cp.Audit = append([]BenefitAuditEvent(nil), benefit.Audit...)
		r.benefits[benefit.ID] = &cp
	}
}
