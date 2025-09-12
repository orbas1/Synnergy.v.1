package core

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// GovernanceProposal represents a proposal created using SYN300 tokens.
type GovernanceProposal struct {
	ID          uint64
	Creator     string
	Description string
	Approvals   map[string]bool
	Rejections  map[string]bool
	Executed    bool
	CreatedAt   time.Time
}

// SYN300Token provides governance features like delegation and on-chain proposals.
type SYN300Token struct {
	mu          sync.RWMutex
	balances    map[string]uint64
	delegations map[string]string
	proposals   map[uint64]*GovernanceProposal
	nextPropID  uint64
}

var (
	// ErrSelfDelegation is returned when an owner delegates to themselves.
	ErrSelfDelegation = errors.New("cannot delegate to self")
	// ErrNoDelegation is returned when revoking a missing delegation.
	ErrNoDelegation = errors.New("delegation not found")
	// ErrEmptyDescription is returned when creating a proposal without a description.
	ErrEmptyDescription = errors.New("proposal description required")
	// ErrNoVotingPower indicates the caller lacks voting power for the action.
	ErrNoVotingPower = errors.New("insufficient voting power")
	// ErrAlreadyVoted indicates the voter has already cast the same vote.
	ErrAlreadyVoted = errors.New("already cast this vote")
)

// NewSYN300Token initialises a SYN300 token with an optional map of starting balances.
func NewSYN300Token(initial map[string]uint64) *SYN300Token {
	cpy := make(map[string]uint64, len(initial))
	for k, v := range initial {
		cpy[k] = v
	}
	return &SYN300Token{
		balances:    cpy,
		delegations: make(map[string]string),
		proposals:   make(map[uint64]*GovernanceProposal),
		nextPropID:  1,
	}
}

// Delegate assigns the owner's voting power to another address. Returns an error
// if the delegation is invalid (such as delegating to self).
func (t *SYN300Token) Delegate(owner, delegate string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if owner == delegate {
		return ErrSelfDelegation
	}
	if delegate == "" {
		delete(t.delegations, owner)
	} else {
		t.delegations[owner] = delegate
	}
	return nil
}

// RevokeDelegation removes an existing delegation for the owner. Returns an
// error if no delegation exists.
func (t *SYN300Token) RevokeDelegation(owner string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.delegations[owner]; !ok {
		return ErrNoDelegation
	}
	delete(t.delegations, owner)
	return nil
}

// VotingPower returns the voting power of the specified address including delegated tokens.
func (t *SYN300Token) VotingPower(addr string) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.votingPowerLocked(addr)
}

func (t *SYN300Token) votingPowerLocked(addr string) uint64 {
	power := t.balances[addr]
	for owner, delegate := range t.delegations {
		if delegate == addr {
			power += t.balances[owner]
		}
	}
	return power
}

// CreateProposal registers a new governance proposal and returns its ID.
// The creator must have voting power and the description must be non-empty.
func (t *SYN300Token) CreateProposal(creator, description string) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if strings.TrimSpace(description) == "" {
		return 0, ErrEmptyDescription
	}
	if t.votingPowerLocked(creator) == 0 {
		return 0, ErrNoVotingPower
	}
	id := t.nextPropID
	t.nextPropID++
	t.proposals[id] = &GovernanceProposal{
		ID:          id,
		Creator:     creator,
		Description: description,
		Approvals:   make(map[string]bool),
		Rejections:  make(map[string]bool),
		CreatedAt:   time.Now(),
	}
	return id, nil
}

// Vote records a vote on a proposal from a given address.
func (t *SYN300Token) Vote(id uint64, voter string, approve bool) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.proposals[id]
	if !ok {
		return errors.New("proposal not found")
	}
	if p.Executed {
		return errors.New("proposal already executed")
	}
	if t.votingPowerLocked(voter) == 0 {
		return ErrNoVotingPower
	}
	if approve {
		if p.Approvals[voter] {
			return ErrAlreadyVoted
		}
		p.Approvals[voter] = true
		delete(p.Rejections, voter)
	} else {
		if p.Rejections[voter] {
			return ErrAlreadyVoted
		}
		p.Rejections[voter] = true
		delete(p.Approvals, voter)
	}
	return nil
}

// Execute finalises a proposal if the approval voting power meets the quorum.
func (t *SYN300Token) Execute(id uint64, quorum uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	p, ok := t.proposals[id]
	if !ok {
		return errors.New("proposal not found")
	}
	if p.Executed {
		return errors.New("proposal already executed")
	}
	var power uint64
	for voter := range p.Approvals {
		power += t.votingPowerLocked(voter)
	}
	if power < quorum {
		return errors.New("quorum not reached")
	}
	p.Executed = true
	return nil
}

// ProposalStatus returns a copy of the proposal for external inspection.
func (t *SYN300Token) ProposalStatus(id uint64) (*GovernanceProposal, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	p, ok := t.proposals[id]
	if !ok {
		return nil, errors.New("proposal not found")
	}
	// return shallow copy to prevent modification
	cp := *p
	cp.Approvals = make(map[string]bool, len(p.Approvals))
	for k, v := range p.Approvals {
		cp.Approvals[k] = v
	}
	cp.Rejections = make(map[string]bool, len(p.Rejections))
	for k, v := range p.Rejections {
		cp.Rejections[k] = v
	}
	return &cp, nil
}

// ListProposals returns all proposals currently registered.
func (t *SYN300Token) ListProposals() []*GovernanceProposal {
	t.mu.RLock()
	defer t.mu.RUnlock()
	list := make([]*GovernanceProposal, 0, len(t.proposals))
	for _, p := range t.proposals {
		cp := *p
		cp.Approvals = make(map[string]bool, len(p.Approvals))
		for k, v := range p.Approvals {
			cp.Approvals[k] = v
		}
		cp.Rejections = make(map[string]bool, len(p.Rejections))
		for k, v := range p.Rejections {
			cp.Rejections[k] = v
		}
		list = append(list, &cp)
	}
	return list
}
