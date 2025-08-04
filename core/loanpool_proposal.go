package core

import "time"

// LoanProposal represents a request for funds from the loan pool.
type LoanProposal struct {
	ID          uint64
	Creator     string
	Recipient   string
	Type        string
	Amount      uint64
	Description string
	Votes       map[string]bool
	Deadline    time.Time
	Approved    bool
	Disbursed   bool
}

// NewLoanProposal creates a new proposal instance.
func NewLoanProposal(id uint64, creator, recipient, typ string, amount uint64, desc string, duration time.Duration) *LoanProposal {
	return &LoanProposal{
		ID:          id,
		Creator:     creator,
		Recipient:   recipient,
		Type:        typ,
		Amount:      amount,
		Description: desc,
		Votes:       make(map[string]bool),
		Deadline:    time.Now().Add(duration),
	}
}

// Vote records a vote from a given address.
func (p *LoanProposal) Vote(voter string) {
	p.Votes[voter] = true
}

// VoteCount returns the number of votes cast.
func (p *LoanProposal) VoteCount() int {
	return len(p.Votes)
}

// IsExpired returns true if the proposal voting deadline has passed.
func (p *LoanProposal) IsExpired(now time.Time) bool {
	return now.After(p.Deadline)
}
