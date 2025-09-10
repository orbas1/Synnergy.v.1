package core

import (
	"errors"
	"fmt"
	"sync"
)

// DAO represents a decentralised autonomous organisation.
type DAO struct {
	mu        sync.RWMutex
	ID        string
	Name      string
	Creator   string
	Members   map[string]string // addr -> role
	Proposals []*DAOProposal
}

// DAOManager tracks all DAOs.
type DAOManager struct {
	mu       sync.RWMutex
	daos     map[string]*DAO
	nextID   int
	relayers map[string]struct{}
}

// NewDAOManager creates an empty DAO manager.
func NewDAOManager() *DAOManager {
	return &DAOManager{daos: make(map[string]*DAO), nextID: 1, relayers: make(map[string]struct{})}
}

// AuthorizeRelayer adds an address to the DAO manager whitelist.
func (m *DAOManager) AuthorizeRelayer(addr string) {
	m.mu.Lock()
	m.relayers[addr] = struct{}{}
	m.mu.Unlock()
}

// RevokeRelayer removes an address from the whitelist.
func (m *DAOManager) RevokeRelayer(addr string) {
	m.mu.Lock()
	delete(m.relayers, addr)
	m.mu.Unlock()
}

// IsRelayerAuthorized returns true if the address is whitelisted.
func (m *DAOManager) IsRelayerAuthorized(addr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.relayers[addr]
	return ok
}

// Create initialises a new DAO. The creator must be an authorised relayer.
func (m *DAOManager) Create(name, creator string) (*DAO, error) {
	if name == "" || creator == "" {
		return nil, errors.New("invalid parameters")
	}
	if !m.IsRelayerAuthorized(creator) {
		return nil, errors.New("unauthorized relayer")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%d", m.nextID)
	m.nextID++
	dao := &DAO{ID: id, Name: name, Creator: creator, Members: map[string]string{creator: "admin"}}
	m.daos[id] = dao
	return dao, nil
}

// Join adds an address to the DAO with member role.
func (m *DAOManager) Join(id, addr string) error {
	if !m.IsRelayerAuthorized(addr) {
		return errors.New("unauthorized relayer")
	}
	m.mu.RLock()
	dao, ok := m.daos[id]
	m.mu.RUnlock()
	if !ok {
		return errors.New("dao not found")
	}
	dao.mu.Lock()
	defer dao.mu.Unlock()
	if _, exists := dao.Members[addr]; exists {
		return errors.New("member already exists")
	}
	dao.Members[addr] = "member"
	return nil
}

// Leave removes an address from the DAO.
func (m *DAOManager) Leave(id, addr string) error {
	if !m.IsRelayerAuthorized(addr) {
		return errors.New("unauthorized relayer")
	}
	m.mu.RLock()
	dao, ok := m.daos[id]
	m.mu.RUnlock()
	if !ok {
		return errors.New("dao not found")
	}
	dao.mu.Lock()
	defer dao.mu.Unlock()
	if _, ok := dao.Members[addr]; !ok {
		return errors.New("member not found")
	}
	delete(dao.Members, addr)
	return nil
}

// Info returns details about a DAO.
func (m *DAOManager) Info(id string) (*DAO, error) {
	m.mu.RLock()
	dao, ok := m.daos[id]
	m.mu.RUnlock()
	if !ok {
		return nil, errors.New("dao not found")
	}
	return dao, nil
}

// List returns all DAOs.
func (m *DAOManager) List() []*DAO {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*DAO, 0, len(m.daos))
	for _, dao := range m.daos {
		out = append(out, dao)
	}
	return out
}
