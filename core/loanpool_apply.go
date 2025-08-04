package core

import "errors"

// LoanApplication represents a simplified loan request reviewed by voting.
type LoanApplication struct {
	ID         uint64
	Applicant  string
	Amount     uint64
	TermMonths uint32
	Purpose    string
	Votes      map[string]bool
	Approved   bool
	Disbursed  bool
}

// LoanPoolApply manages loan applications using an underlying LoanPool treasury.
type LoanPoolApply struct {
	Pool         *LoanPool
	Applications map[uint64]*LoanApplication
	nextID       uint64
}

// NewLoanPoolApply creates a new application manager.
func NewLoanPoolApply(pool *LoanPool) *LoanPoolApply {
	return &LoanPoolApply{
		Pool:         pool,
		Applications: make(map[uint64]*LoanApplication),
		nextID:       1,
	}
}

// Submit adds a new loan application.
func (l *LoanPoolApply) Submit(applicant string, amount uint64, termMonths uint32, purpose string) uint64 {
	id := l.nextID
	l.nextID++
	l.Applications[id] = &LoanApplication{
		ID:         id,
		Applicant:  applicant,
		Amount:     amount,
		TermMonths: termMonths,
		Purpose:    purpose,
		Votes:      make(map[string]bool),
	}
	return id
}

// Vote records a vote for an application.
func (l *LoanPoolApply) Vote(voter string, id uint64) error {
	app, ok := l.Applications[id]
	if !ok {
		return errors.New("application not found")
	}
	app.Votes[voter] = true
	return nil
}

// Process finalises applications approving those with at least one vote.
func (l *LoanPoolApply) Process() {
	for _, app := range l.Applications {
		if !app.Approved && len(app.Votes) > 0 {
			app.Approved = true
		}
	}
}

// Disburse pays out an approved application.
func (l *LoanPoolApply) Disburse(id uint64) error {
	app, ok := l.Applications[id]
	if !ok {
		return errors.New("application not found")
	}
	if !app.Approved || app.Disbursed {
		return errors.New("application not approved or already disbursed")
	}
	if l.Pool.Treasury < app.Amount {
		return errors.New("insufficient treasury")
	}
	l.Pool.Treasury -= app.Amount
	app.Disbursed = true
	return nil
}

// Get returns an application by ID.
func (l *LoanPoolApply) Get(id uint64) (*LoanApplication, bool) {
	app, ok := l.Applications[id]
	return app, ok
}

// List returns all applications.
func (l *LoanPoolApply) List() []*LoanApplication {
	res := make([]*LoanApplication, 0, len(l.Applications))
	for _, a := range l.Applications {
		res = append(res, a)
	}
	return res
}
