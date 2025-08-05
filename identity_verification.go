package synnergy

import (
	"errors"
	"sync"
	"time"
)

// IdentityInfo contains basic identity metadata.
type IdentityInfo struct {
	Name        string
	DateOfBirth string
	Nationality string
}

// VerificationLog records a verification attempt.
type VerificationLog struct {
	Method    string
	Timestamp time.Time
}

// IdentityService manages verified addresses on the ledger.
type IdentityService struct {
	mu         sync.RWMutex
	identities map[string]IdentityInfo
	logs       map[string][]VerificationLog
}

// NewIdentityService creates a new IdentityService instance.
func NewIdentityService() *IdentityService {
	return &IdentityService{
		identities: make(map[string]IdentityInfo),
		logs:       make(map[string][]VerificationLog),
	}
}

// Register stores identity information for an address.
func (s *IdentityService) Register(addr, name, dob, nationality string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.identities[addr]; exists {
		return errors.New("identity already registered")
	}
	s.identities[addr] = IdentityInfo{Name: name, DateOfBirth: dob, Nationality: nationality}
	return nil
}

// Verify records a verification method for an address.
func (s *IdentityService) Verify(addr, method string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.identities[addr]; !exists {
		return errors.New("identity not registered")
	}
	log := VerificationLog{Method: method, Timestamp: time.Now()}
	s.logs[addr] = append(s.logs[addr], log)
	return nil
}

// Info retrieves identity information for an address.
func (s *IdentityService) Info(addr string) (IdentityInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.identities[addr]
	return info, ok
}

// Logs returns verification logs for an address.
func (s *IdentityService) Logs(addr string) []VerificationLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entries := s.logs[addr]
	out := make([]VerificationLog, len(entries))
	copy(out, entries)
	return out
}
