package core

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	ErrLoanPoolPaused        = errors.New("loan pool is paused")
	ErrInvalidProposal       = errors.New("invalid proposal")
	errLoanProposalNotFound  = errors.New("proposal not found")
	errLoanProposalCancelled = errors.New("cannot cancel proposal")
)

// LoanPool manages loan proposals and disbursements.
type LoanPool struct {
	Treasury  uint64
	Proposals map[uint64]*LoanProposal
	nextID    uint64
	Paused    bool

	mu sync.RWMutex
}

// NewLoanPool creates a new loan pool with the given treasury balance.
func NewLoanPool(treasury uint64) *LoanPool {
	return &LoanPool{
		Treasury:  treasury,
		Proposals: make(map[uint64]*LoanProposal),
		nextID:    1,
	}
}

// SetPaused toggles whether new proposals may be submitted.
func (lp *LoanPool) SetPaused(paused bool) {
	lp.mu.Lock()
	lp.Paused = paused
	lp.mu.Unlock()
}

// SubmitProposal adds a new loan proposal to the pool.
func (lp *LoanPool) SubmitProposal(creator, recipient, typ string, amount uint64, desc string) (uint64, error) {
	if amount == 0 || recipient == "" || creator == "" {
		return 0, ErrInvalidProposal
	}
	lp.mu.Lock()
	defer lp.mu.Unlock()
	if lp.Paused {
		return 0, ErrLoanPoolPaused
	}
	id := lp.nextID
	lp.nextID++
	lp.Proposals[id] = NewLoanProposal(id, creator, recipient, typ, amount, desc, 24*time.Hour)
	return id, nil
}

// VoteProposal records a vote for a proposal.
func (lp *LoanPool) VoteProposal(voter string, id uint64) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	p, ok := lp.Proposals[id]
	if !ok {
		return errLoanProposalNotFound
	}
	if p.IsExpired(time.Now()) {
		return errors.New("voting closed")
	}
	p.Vote(voter)
	return nil
}

// Tick evaluates proposals for approval based on votes and deadlines.
func (lp *LoanPool) Tick() {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	now := time.Now()
	for _, p := range lp.Proposals {
		if !p.Approved && p.VoteCount() > 0 && !p.IsExpired(now) {
			p.Approved = true
		}
	}
}

// Disburse sends funds for an approved proposal if treasury allows.
func (lp *LoanPool) Disburse(id uint64) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	p, ok := lp.Proposals[id]
	if !ok {
		return errLoanProposalNotFound
	}
	if !p.Approved || p.Disbursed {
		return errors.New("proposal not approved or already disbursed")
	}
	if lp.Treasury < p.Amount {
		return errors.New("insufficient treasury")
	}
	lp.Treasury -= p.Amount
	p.Disbursed = true
	return nil
}

// GetProposal returns a proposal by ID.
func (lp *LoanPool) GetProposal(id uint64) (*LoanProposal, bool) {
	lp.mu.RLock()
	defer lp.mu.RUnlock()
	p, ok := lp.Proposals[id]
	if !ok {
		return nil, false
	}
	return p, true
}

// ListProposals returns proposals sorted by ID.
func (lp *LoanPool) ListProposals() []*LoanProposal {
	lp.mu.RLock()
	defer lp.mu.RUnlock()
	res := make([]*LoanProposal, 0, len(lp.Proposals))
	for _, p := range lp.Proposals {
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// CancelProposal removes an active proposal if requested by the creator.
func (lp *LoanPool) CancelProposal(creator string, id uint64) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	p, ok := lp.Proposals[id]
	if !ok {
		return errLoanProposalNotFound
	}
	if p.Creator != creator || p.Disbursed {
		return errLoanProposalCancelled
	}
	delete(lp.Proposals, id)
	return nil
}

// ExtendProposal extends the voting deadline for a proposal.
func (lp *LoanPool) ExtendProposal(creator string, id uint64, hrs int) error {
	if hrs <= 0 {
		return errors.New("extension must be positive")
	}
	lp.mu.Lock()
	defer lp.mu.Unlock()
	p, ok := lp.Proposals[id]
	if !ok {
		return errLoanProposalNotFound
	}
	if p.Creator != creator {
		return errors.New("unauthorised")
	}
	p.Deadline = p.Deadline.Add(time.Duration(hrs) * time.Hour)
	return nil
}

// TreasuryBalance returns the current treasury funds.
func (lp *LoanPool) TreasuryBalance() uint64 {
	lp.mu.RLock()
	defer lp.mu.RUnlock()
	return lp.Treasury
}

func (lp *LoanPool) withdraw(amount uint64) error {
	lp.mu.Lock()
	defer lp.mu.Unlock()
	if lp.Treasury < amount {
		return errors.New("insufficient treasury")
	}
	lp.Treasury -= amount
	return nil
}
