package core

import "time"

// LoanProposalView provides a serialisable representation of a LoanProposal.
type LoanProposalView struct {
	ID          uint64    `json:"id"`
	Creator     string    `json:"creator"`
	Recipient   string    `json:"recipient"`
	Type        string    `json:"type"`
	Amount      uint64    `json:"amount"`
	Description string    `json:"description"`
	Votes       int       `json:"votes"`
	Deadline    time.Time `json:"deadline"`
	Approved    bool      `json:"approved"`
	Disbursed   bool      `json:"disbursed"`
}

// NewLoanProposalView constructs a view from a proposal.
func NewLoanProposalView(p *LoanProposal) LoanProposalView {
	return LoanProposalView{
		ID:          p.ID,
		Creator:     p.Creator,
		Recipient:   p.Recipient,
		Type:        p.Type,
		Amount:      p.Amount,
		Description: p.Description,
		Votes:       len(p.Votes),
		Deadline:    p.Deadline,
		Approved:    p.Approved,
		Disbursed:   p.Disbursed,
	}
}

// ProposalInfo returns a view for a proposal by ID if it exists.
func (lp *LoanPool) ProposalInfo(id uint64) (LoanProposalView, bool) {
	lp.mu.RLock()
	defer lp.mu.RUnlock()
	p, ok := lp.Proposals[id]
	if !ok {
		return LoanProposalView{}, false
	}
	return NewLoanProposalView(p), true
}

// ProposalViews lists all proposals in view form.
func (lp *LoanPool) ProposalViews() []LoanProposalView {
	lp.mu.RLock()
	defer lp.mu.RUnlock()
	views := make([]LoanProposalView, 0, len(lp.Proposals))
	for _, p := range lp.Proposals {
		views = append(views, NewLoanProposalView(p))
	}
	return views
}

// LoanApplicationView provides a serialisable representation of a LoanApplication.
type LoanApplicationView struct {
	ID         uint64 `json:"id"`
	Applicant  string `json:"applicant"`
	Amount     uint64 `json:"amount"`
	TermMonths uint32 `json:"term_months"`
	Purpose    string `json:"purpose"`
	Votes      int    `json:"votes"`
	Approved   bool   `json:"approved"`
	Disbursed  bool   `json:"disbursed"`
}

// NewLoanApplicationView constructs a view from an application.
func NewLoanApplicationView(a *LoanApplication) LoanApplicationView {
	return LoanApplicationView{
		ID:         a.ID,
		Applicant:  a.Applicant,
		Amount:     a.Amount,
		TermMonths: a.TermMonths,
		Purpose:    a.Purpose,
		Votes:      len(a.Votes),
		Approved:   a.Approved,
		Disbursed:  a.Disbursed,
	}
}

// ApplicationInfo returns a view for an application by ID if it exists.
func (l *LoanPoolApply) ApplicationInfo(id uint64) (LoanApplicationView, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	app, ok := l.Applications[id]
	if !ok {
		return LoanApplicationView{}, false
	}
	return NewLoanApplicationView(app), true
}

// ApplicationViews lists all applications in view form.
func (l *LoanPoolApply) ApplicationViews() []LoanApplicationView {
	l.mu.RLock()
	defer l.mu.RUnlock()
	views := make([]LoanApplicationView, 0, len(l.Applications))
	for _, a := range l.Applications {
		views = append(views, NewLoanApplicationView(a))
	}
	return views
}
