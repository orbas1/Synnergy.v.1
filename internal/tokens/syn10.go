package tokens

import (
	"errors"
	"sync"
	"time"
)

// SYN10Token represents a central bank digital currency pegged to fiat. It
// embeds BaseToken for ledger operations and adds concurrency-safe management of
// issuer metadata and fiat exchange rates.
type SYN10Token struct {
	*BaseToken
	mu             sync.RWMutex
	issuer         string
	exchangeRate   float64
	circuitBreaker bool
	participants   map[string]*CBDCParticipant
	auditTrail     []CBDCAuditRecord
}

// NewSYN10Token initialises a SYN10 token with the given metadata.
func NewSYN10Token(id TokenID, name, symbol, issuer string, rate float64, decimals uint8) *SYN10Token {
	return &SYN10Token{
		BaseToken:    NewBaseToken(id, name, symbol, decimals),
		issuer:       issuer,
		exchangeRate: rate,
		participants: make(map[string]*CBDCParticipant),
	}
}

// SetExchangeRate updates the fiat exchange rate for the CBDC.
func (t *SYN10Token) SetExchangeRate(rate float64) {
	t.mu.Lock()
	t.exchangeRate = rate
	t.mu.Unlock()
}

// ComplianceTier defines the level of due diligence applied to a participant.
type ComplianceTier string

const (
	// TierRetail represents individual retail users of the CBDC.
	TierRetail ComplianceTier = "retail"
	// TierCommercial represents commercial entities such as merchants.
	TierCommercial ComplianceTier = "commercial"
	// TierInstitution represents financial institutions and market makers.
	TierInstitution ComplianceTier = "institution"
)

// CBDCAuditRecord captures a transaction level audit trail ensuring CLI, VM and
// consensus layers can access deterministic metadata.
type CBDCAuditRecord struct {
	From      string
	To        string
	Amount    uint64
	Tier      ComplianceTier
	Timestamp time.Time
	Metadata  map[string]string
}

// CBDCParticipant describes a participant registered with the CBDC network.
type CBDCParticipant struct {
	Address       string
	Tier          ComplianceTier
	DailyLimit    uint64
	AMLFlagged    bool
	Metadata      map[string]string
	spentToday    uint64
	windowStart   time.Time
	auditMetadata map[string]string
}

var (
	// ErrCBDCParticipantUnknown indicates the caller attempted to interact
	// with an unregistered participant.
	ErrCBDCParticipantUnknown = errors.New("participant not registered")
	// ErrCBDCParticipantFlagged is returned when a flagged participant
	// attempts to send or receive funds.
	ErrCBDCParticipantFlagged = errors.New("participant is under review")
	// ErrCBDCDailyLimitExceeded is returned when a transfer would breach
	// the participant's configured risk limit.
	ErrCBDCDailyLimitExceeded = errors.New("daily transaction limit exceeded")
	// ErrCBDCCircuitBreakerActive indicates emergency controls are active.
	ErrCBDCCircuitBreakerActive = errors.New("cbdc circuit breaker active")
)

// RegisterParticipant registers or updates a participant profile.
func (t *SYN10Token) RegisterParticipant(addr string, tier ComplianceTier, limit uint64, meta map[string]string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	cp := &CBDCParticipant{
		Address:       addr,
		Tier:          tier,
		DailyLimit:    limit,
		Metadata:      make(map[string]string),
		auditMetadata: make(map[string]string),
		windowStart:   time.Now(),
	}
	for k, v := range meta {
		cp.Metadata[k] = v
	}
	t.participants[addr] = cp
}

// UpdateDailyLimit adjusts the risk budget for a participant.
func (t *SYN10Token) UpdateDailyLimit(addr string, limit uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	cp, ok := t.participants[addr]
	if !ok {
		return ErrCBDCParticipantUnknown
	}
	cp.DailyLimit = limit
	return nil
}

// FlagParticipant toggles AML review status for a participant.
func (t *SYN10Token) FlagParticipant(addr string, flagged bool, metadata map[string]string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	cp, ok := t.participants[addr]
	if !ok {
		return ErrCBDCParticipantUnknown
	}
	cp.AMLFlagged = flagged
	for k, v := range metadata {
		if cp.auditMetadata == nil {
			cp.auditMetadata = make(map[string]string)
		}
		cp.auditMetadata[k] = v
	}
	return nil
}

// SetCircuitBreaker activates or deactivates emergency controls preventing
// further state mutations until investigations are complete.
func (t *SYN10Token) SetCircuitBreaker(enabled bool) {
	t.mu.Lock()
	t.circuitBreaker = enabled
	t.mu.Unlock()
}

// TransferCBDC executes a CBDC transfer after enforcing risk and compliance
// checks. It records a deterministic audit record consumable by CLI, consensus
// and VM modules.
func (t *SYN10Token) TransferCBDC(from, to string, amount uint64, metadata map[string]string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.circuitBreaker {
		return ErrCBDCCircuitBreakerActive
	}
	sender, ok := t.participants[from]
	if !ok {
		return ErrCBDCParticipantUnknown
	}
	receiver, ok := t.participants[to]
	if !ok {
		return ErrCBDCParticipantUnknown
	}
	if sender.AMLFlagged || receiver.AMLFlagged {
		return ErrCBDCParticipantFlagged
	}

	now := time.Now()
	if now.Sub(sender.windowStart) >= 24*time.Hour {
		sender.windowStart = now
		sender.spentToday = 0
	}
	if sender.DailyLimit > 0 && sender.spentToday+amount > sender.DailyLimit {
		return ErrCBDCDailyLimitExceeded
	}

	if err := t.BaseToken.Transfer(from, to, amount); err != nil {
		return err
	}

	sender.spentToday += amount
	auditMeta := make(map[string]string)
	for k, v := range sender.auditMetadata {
		auditMeta[k] = v
	}
	for k, v := range metadata {
		auditMeta[k] = v
	}
	t.auditTrail = append(t.auditTrail, CBDCAuditRecord{
		From:      from,
		To:        to,
		Amount:    amount,
		Tier:      sender.Tier,
		Timestamp: now,
		Metadata:  auditMeta,
	})
	return nil
}

// AuditTrail returns the most recent audit records up to the provided limit.
func (t *SYN10Token) AuditTrail(limit int) []CBDCAuditRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if limit <= 0 || limit >= len(t.auditTrail) {
		cp := make([]CBDCAuditRecord, len(t.auditTrail))
		copy(cp, t.auditTrail)
		return cp
	}
	cp := make([]CBDCAuditRecord, limit)
	copy(cp, t.auditTrail[len(t.auditTrail)-limit:])
	return cp
}

// SYN10Info summarises token configuration.
type SYN10Info struct {
	Name         string
	Symbol       string
	Issuer       string
	ExchangeRate float64
	TotalSupply  uint64
	Participants int
}

// Info returns the current token information.
func (t *SYN10Token) Info() SYN10Info {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return SYN10Info{
		Name:         t.Name(),
		Symbol:       t.Symbol(),
		Issuer:       t.issuer,
		ExchangeRate: t.exchangeRate,
		TotalSupply:  t.TotalSupply(),
		Participants: len(t.participants),
	}
}
