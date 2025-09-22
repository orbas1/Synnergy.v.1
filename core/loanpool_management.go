package core

// LoanPoolStats summarises loan pool state for administrators.
type LoanPoolStats struct {
	Treasury       uint64
	ProposalCount  int
	ApprovedCount  int
	DisbursedCount int
}

// LoanPoolManager provides administrative helpers around LoanPool.
type LoanPoolManager struct {
	Pool *LoanPool
}

// NewLoanPoolManager creates a new manager for a pool.
func NewLoanPoolManager(p *LoanPool) *LoanPoolManager {
	return &LoanPoolManager{Pool: p}
}

// Pause stops new proposals from being submitted.
func (m *LoanPoolManager) Pause() {
	if m.Pool != nil {
		m.Pool.SetPaused(true)
	}
}

// Resume allows proposal submissions again.
func (m *LoanPoolManager) Resume() {
	if m.Pool != nil {
		m.Pool.SetPaused(false)
	}
}

// Stats returns a summary of the pool's state.
func (m *LoanPoolManager) Stats() LoanPoolStats {
	if m.Pool == nil {
		return LoanPoolStats{}
	}
	m.Pool.mu.RLock()
	defer m.Pool.mu.RUnlock()
	stats := LoanPoolStats{Treasury: m.Pool.Treasury}
	for _, p := range m.Pool.Proposals {
		stats.ProposalCount++
		if p.Approved {
			stats.ApprovedCount++
		}
		if p.Disbursed {
			stats.DisbursedCount++
		}
	}
	return stats
}
