package core

import (
	"errors"
	"fmt"
)

// DAO represents a decentralised autonomous organisation.
type DAO struct {
	ID        string
	Name      string
	Creator   string
	Members   map[string]string // addr -> role
	Proposals []*DAOProposal
}

// DAOManager tracks all DAOs.
type DAOManager struct {
	daos   map[string]*DAO
	nextID int
}

// NewDAOManager creates an empty DAO manager.
func NewDAOManager() *DAOManager {
	return &DAOManager{daos: make(map[string]*DAO), nextID: 1}
}

// Create initialises a new DAO.
func (m *DAOManager) Create(name, creator string) *DAO {
	id := fmt.Sprintf("%d", m.nextID)
	m.nextID++
	dao := &DAO{ID: id, Name: name, Creator: creator, Members: map[string]string{creator: "admin"}}
	m.daos[id] = dao
	return dao
}

// Join adds an address to the DAO with member role.
func (m *DAOManager) Join(id, addr string) error {
	dao, ok := m.daos[id]
	if !ok {
		return errors.New("dao not found")
	}
	if _, exists := dao.Members[addr]; exists {
		return errors.New("member already exists")
	}
	dao.Members[addr] = "member"
	return nil
}

// Leave removes an address from the DAO.
func (m *DAOManager) Leave(id, addr string) error {
	dao, ok := m.daos[id]
	if !ok {
		return errors.New("dao not found")
	}
	delete(dao.Members, addr)
	return nil
}

// Info returns details about a DAO.
func (m *DAOManager) Info(id string) (*DAO, error) {
	dao, ok := m.daos[id]
	if !ok {
		return nil, errors.New("dao not found")
	}
	return dao, nil
}

// List returns all DAOs.
func (m *DAOManager) List() []*DAO {
	out := make([]*DAO, 0, len(m.daos))
	for _, dao := range m.daos {
		out = append(out, dao)
	}
	return out
}
