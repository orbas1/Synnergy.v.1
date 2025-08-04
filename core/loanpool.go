package core

import (
	"errors"
	"sort"
	"time"
)

// LoanPool manages loan proposals and disbursements.
type LoanPool struct {
	Treasury  uint64
	Proposals map[uint64]*LoanProposal
	nextID    uint64
	Paused    bool
}

// NewLoanPool creates a new loan pool with the given treasury balance.
func NewLoanPool(treasury uint64) *LoanPool {
	return &LoanPool{
		Treasury:  treasury,
		Proposals: make(map[uint64]*LoanProposal),
		nextID:    1,
	}
}

// SubmitProposal adds a new loan proposal to the pool.
func (lp *LoanPool) SubmitProposal(creator, recipient, typ string, amount uint64, desc string) (uint64, error) {
	if lp.Paused {
		return 0, errors.New("loan pool is paused")
	}
	id := lp.nextID
	lp.nextID++
	lp.Proposals[id] = NewLoanProposal(id, creator, recipient, typ, amount, desc, time.Hour*24)
	return id, nil
}

// VoteProposal records a vote for a proposal.
func (lp *LoanPool) VoteProposal(voter string, id uint64) error {
	p, ok := lp.Proposals[id]
	if !ok {
		return errors.New("proposal not found")
	}
	if p.IsExpired(time.Now()) {
		return errors.New("voting closed")
	}
	p.Vote(voter)
	return nil
}

// Tick evaluates proposals for approval based on votes and deadlines.
func (lp *LoanPool) Tick() {
	now := time.Now()
	for _, p := range lp.Proposals {
		if !p.Approved && p.VoteCount() > 0 && !p.IsExpired(now) {
			p.Approved = true
		}
	}
}

// Disburse sends funds for an approved proposal if treasury allows.
func (lp *LoanPool) Disburse(id uint64) error {
	p, ok := lp.Proposals[id]
	if !ok {
		return errors.New("proposal not found")
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
	p, ok := lp.Proposals[id]
	return p, ok
}

// ListProposals returns proposals sorted by ID.
func (lp *LoanPool) ListProposals() []*LoanProposal {
	res := make([]*LoanProposal, 0, len(lp.Proposals))
	for _, p := range lp.Proposals {
		res = append(res, p)
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ID < res[j].ID })
	return res
}

// CancelProposal removes an active proposal if requested by the creator.
func (lp *LoanPool) CancelProposal(creator string, id uint64) error {
	p, ok := lp.Proposals[id]
	if !ok {
		return errors.New("proposal not found")
	}
	if p.Creator != creator || p.Disbursed {
		return errors.New("cannot cancel proposal")
	}
	delete(lp.Proposals, id)
	return nil
}

// ExtendProposal extends the voting deadline for a proposal.
func (lp *LoanPool) ExtendProposal(creator string, id uint64, hrs int) error {
	p, ok := lp.Proposals[id]
	if !ok {
		return errors.New("proposal not found")
	}
	if p.Creator != creator {
		return errors.New("unauthorised")
	}
	p.Deadline = p.Deadline.Add(time.Duration(hrs) * time.Hour)
	return nil
}
