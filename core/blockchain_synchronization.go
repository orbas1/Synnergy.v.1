package core

import (
	"fmt"
	"sync"
)

// SyncManager coordinates block download and verification to keep a node's
// ledger in sync with the network. The implementation here tracks state in
// memory and exposes helpers for the CLI to control the process.
type SyncManager struct {
	mu         sync.RWMutex
	ledger     *Ledger
	running    bool
	lastHeight int
}

// NewSyncManager returns a new SyncManager bound to a ledger.
func NewSyncManager(l *Ledger) *SyncManager {
	return &SyncManager{ledger: l}
}

// Start marks the manager as running.
func (s *SyncManager) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = true
}

// Stop halts synchronization.
func (s *SyncManager) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

// Status reports whether synchronization is active and the last synced height.
func (s *SyncManager) Status() (bool, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running, s.lastHeight
}

// Once performs a single synchronization round by checking the ledger head.
func (s *SyncManager) Once() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return fmt.Errorf("synchronization not running")
	}
	h, _ := s.ledger.Head()
	s.lastHeight = h
	return nil
}
