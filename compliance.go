package synnergy

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AuditEntry records a single compliance event for an address.
type AuditEntry struct {
	Address   string            // Address associated with the event
	Event     string            // Name of the event
	Metadata  map[string]string // Optional metadata about the event
	Timestamp time.Time         // Time the event occurred
}

// KYCRecord stores the commitment for an address and when it was recorded.
type KYCRecord struct {
	Commitment string    // Hex encoded commitment of the KYC document
	Timestamp  time.Time // When the commitment was recorded
}

// FraudSignal represents a fraud warning with a severity score.
type FraudSignal struct {
	Severity  int       // Severity score for the fraud signal
	Timestamp time.Time // When the signal was recorded
}

// ComplianceTransaction represents a minimal transaction for compliance checks.
// It uses float amounts because compliance monitoring may involve fractional
// values independent of on-chain integer amounts.
type ComplianceTransaction struct {
	ID     string  // Identifier of the transaction
	From   string  // Sender address
	To     string  // Receiver address
	Amount float64 // Amount involved in the transaction
}

// ComplianceService manages KYC records and fraud risk scores.
type ComplianceService struct {
	mu         sync.RWMutex
	kycs       map[string]KYCRecord
	frauds     map[string][]FraudSignal
	audits     map[string][]AuditEntry
	riskScores map[string]int
}

// NewComplianceService creates a new ComplianceService instance.
func NewComplianceService() *ComplianceService {
	return &ComplianceService{
		kycs:       make(map[string]KYCRecord),
		frauds:     make(map[string][]FraudSignal),
		audits:     make(map[string][]AuditEntry),
		riskScores: make(map[string]int),
	}
}

// ValidateKYC validates and stores a KYC document commitment for an address.
func (s *ComplianceService) ValidateKYC(address string, kycData []byte) (string, error) {
	if address == "" {
		return "", errors.New("address required")
	}
	hash := sha256.Sum256(kycData)
	commitment := hex.EncodeToString(hash[:])
	s.mu.Lock()
	s.kycs[address] = KYCRecord{Commitment: commitment, Timestamp: time.Now()}
	s.appendAudit(address, "kyc_validated", nil)
	s.mu.Unlock()
	return commitment, nil
}

// EraseKYC removes stored KYC data for an address.
func (s *ComplianceService) EraseKYC(address string) {
	s.mu.Lock()
	delete(s.kycs, address)
	s.appendAudit(address, "kyc_erased", nil)
	s.mu.Unlock()
}

// RecordFraud records a fraud signal and updates the risk score.
func (s *ComplianceService) RecordFraud(address string, severity int) {
	s.mu.Lock()
	sig := FraudSignal{Severity: severity, Timestamp: time.Now()}
	s.frauds[address] = append(s.frauds[address], sig)
	s.riskScores[address] += severity
	meta := map[string]string{"severity": fmt.Sprintf("%d", severity)}
	s.appendAudit(address, "fraud_signal", meta)
	s.mu.Unlock()
}

// RiskScore returns the accumulated fraud risk score for an address.
func (s *ComplianceService) RiskScore(address string) int {
	s.mu.RLock()
	score := s.riskScores[address]
	s.mu.RUnlock()
	return score
}

// AuditTrail returns a copy of the audit trail for an address.
func (s *ComplianceService) AuditTrail(address string) []AuditEntry {
	s.mu.RLock()
	entries := s.audits[address]
	out := make([]AuditEntry, len(entries))
	copy(out, entries)
	s.mu.RUnlock()
	return out
}

// MonitorTransaction performs a simple anomaly detection on a transaction amount.
// It operates on ComplianceTransaction to avoid colliding with any ledger specific
// transaction type.
func (s *ComplianceService) MonitorTransaction(tx ComplianceTransaction, threshold float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.appendAudit(tx.From, "tx_monitored", map[string]string{"id": tx.ID})
	if tx.Amount > threshold {
		s.appendAudit(tx.From, "tx_anomaly", map[string]string{"id": tx.ID})
		return true
	}
	return false
}

// VerifyZKP verifies a simple commitment for demonstration purposes.
func (s *ComplianceService) VerifyZKP(blob []byte, commitmentHex, proofHex string) bool {
	hash := sha256.Sum256(blob)
	return hex.EncodeToString(hash[:]) == commitmentHex
}

func (s *ComplianceService) appendAudit(addr, event string, metadata map[string]string) {
	s.audits[addr] = append(s.audits[addr], AuditEntry{
		Address:   addr,
		Event:     event,
		Metadata:  metadata,
		Timestamp: time.Now(),
	})
}
