package core

import (
	"errors"
	"fmt"
	"sync"
)

// DAOProposal represents a proposal within a DAO.
type DAOProposal struct {
	mu       sync.RWMutex
	ID       string
	DAOID    string
	Creator  string
	Desc     string
	YesVotes map[string]uint64
	NoVotes  map[string]uint64
	Executed bool
}

// ProposalManager manages DAO proposals.
type ProposalManager struct {
	mu        sync.RWMutex
	proposals map[string]*DAOProposal
	nextID    int
}

// NewProposalManager creates an empty ProposalManager.
func NewProposalManager() *ProposalManager {
	return &ProposalManager{proposals: make(map[string]*DAOProposal), nextID: 1}
}

var (
	errProposalNotFound = errors.New("proposal not found")
	errProposalExecuted = errors.New("proposal executed")
	errNotMember        = errors.New("not a dao member")
	errNotAdmin         = errors.New("not a dao admin")
)

// CreateProposal adds a new proposal to a DAO. Only DAO members may create proposals.
func (pm *ProposalManager) CreateProposal(dao *DAO, creator, desc string) (*DAOProposal, error) {
	if !dao.IsMember(creator) {
		return nil, errNotMember
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	id := fmt.Sprintf("%d", pm.nextID)
	pm.nextID++
	p := &DAOProposal{ID: id, DAOID: dao.ID, Creator: creator, Desc: desc,
		YesVotes: make(map[string]uint64), NoVotes: make(map[string]uint64)}
	pm.proposals[id] = p
	dao.mu.Lock()
	dao.Proposals = append(dao.Proposals, p)
	dao.mu.Unlock()
	return p, nil
}

// Vote casts a vote on a proposal with given weight. Only DAO members may vote.
func (pm *ProposalManager) Vote(dao *DAO, id, voter string, weight uint64, support bool) error {
	if !dao.IsMember(voter) {
		return errNotMember
	}
	pm.mu.RLock()
	p, ok := pm.proposals[id]
	pm.mu.RUnlock()
	if !ok {
		return errProposalNotFound
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Executed {
		return errProposalExecuted
	}
	if support {
		p.YesVotes[voter] = weight
		delete(p.NoVotes, voter)
	} else {
		p.NoVotes[voter] = weight
		delete(p.YesVotes, voter)
	}
	return nil
}

// Results sums yes and no votes for a proposal.
func (pm *ProposalManager) Results(id string) (yes, no uint64, err error) {
	pm.mu.RLock()
	p, ok := pm.proposals[id]
	pm.mu.RUnlock()
	if !ok {
		return 0, 0, errProposalNotFound
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, w := range p.YesVotes {
		yes += w
	}
	for _, w := range p.NoVotes {
		no += w
	}
	return yes, no, nil
}

// Execute marks a proposal as executed. Only DAO admins may execute proposals.
func (pm *ProposalManager) Execute(dao *DAO, id, requester string) error {
	if !dao.IsAdmin(requester) {
		return errNotAdmin
	}
	pm.mu.RLock()
	p, ok := pm.proposals[id]
	pm.mu.RUnlock()
	if !ok {
		return errProposalNotFound
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Executed {
		return errProposalExecuted
	}
	p.Executed = true
	return nil
}

// Get returns a proposal by ID.
func (pm *ProposalManager) Get(id string) (*DAOProposal, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	p, ok := pm.proposals[id]
	if !ok {
		return nil, errProposalNotFound
	}
	return p, nil
}

// List returns all proposals.
func (pm *ProposalManager) List() []*DAOProposal {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	out := make([]*DAOProposal, 0, len(pm.proposals))
	for _, p := range pm.proposals {
		out = append(out, p)
	}
	return out
}
