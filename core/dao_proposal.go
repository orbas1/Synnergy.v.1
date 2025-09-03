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

// CreateProposal adds a new proposal to a DAO.
func (pm *ProposalManager) CreateProposal(dao *DAO, creator, desc string) *DAOProposal {
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
	return p
}

// Vote casts a vote on a proposal with given weight.
func (pm *ProposalManager) Vote(id, voter string, weight uint64, support bool) error {
	pm.mu.RLock()
	p, ok := pm.proposals[id]
	pm.mu.RUnlock()
	if !ok {
		return errors.New("proposal not found")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Executed {
		return errors.New("proposal executed")
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
		return 0, 0, errors.New("proposal not found")
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

// Execute marks a proposal as executed.
func (pm *ProposalManager) Execute(id string) error {
	pm.mu.RLock()
	p, ok := pm.proposals[id]
	pm.mu.RUnlock()
	if !ok {
		return errors.New("proposal not found")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.Executed {
		return errors.New("proposal executed")
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
		return nil, errors.New("proposal not found")
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
