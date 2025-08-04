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
	mu     sync.RWMutex
	daos   map[string]*DAO
	nextID int
}

// NewDAOManager creates an empty DAO manager.
func NewDAOManager() *DAOManager {
	return &DAOManager{daos: make(map[string]*DAO), nextID: 1}
}

// Create initialises a new DAO.
func (m *DAOManager) Create(name, creator string) *DAO {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := fmt.Sprintf("%d", m.nextID)
	m.nextID++
	dao := &DAO{ID: id, Name: name, Creator: creator, Members: map[string]string{creator: "admin"}}
	m.daos[id] = dao
	return dao
}

// Join adds an address to the DAO with member role.
func (m *DAOManager) Join(id, addr string) error {
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
	m.mu.RLock()
	dao, ok := m.daos[id]
	m.mu.RUnlock()
	if !ok {
		return errors.New("dao not found")
	}
	dao.mu.Lock()
	defer dao.mu.Unlock()
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
