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
	m.Pool.Paused = true
}

// Resume allows proposal submissions again.
func (m *LoanPoolManager) Resume() {
	m.Pool.Paused = false
}

// Stats returns a summary of the pool's state.
func (m *LoanPoolManager) Stats() LoanPoolStats {
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
