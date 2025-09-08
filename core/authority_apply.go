package core

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AuthorityApplication represents a request to become an authority node.
type AuthorityApplication struct {
	ID         string
	Candidate  string
	Role       string
	Desc       string
	Approvals  map[string]bool
	Rejections map[string]bool
	ExpiresAt  time.Time
	Finalized  bool
}

// AuthorityApplicationManager tracks authority node applications.
type AuthorityApplicationManager struct {
	mu       sync.RWMutex
	apps     map[string]*AuthorityApplication
	nextID   int
	registry *AuthorityNodeRegistry
	ttl      time.Duration
}

// NewAuthorityApplicationManager creates a manager with the given registry and ttl.
func NewAuthorityApplicationManager(reg *AuthorityNodeRegistry, ttl time.Duration) *AuthorityApplicationManager {
	return &AuthorityApplicationManager{
		apps:     make(map[string]*AuthorityApplication),
		nextID:   1,
		registry: reg,
		ttl:      ttl,
	}
}

// Submit creates a new application and returns its ID.
func (m *AuthorityApplicationManager) Submit(candidate, role, desc string) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if candidate == "" || role == "" {
		return ""
	}
	id := m.nextID
	m.nextID++
	app := &AuthorityApplication{
		ID:         fmt.Sprintf("%d", id),
		Candidate:  candidate,
		Role:       role,
		Desc:       desc,
		Approvals:  make(map[string]bool),
		Rejections: make(map[string]bool),
		ExpiresAt:  time.Now().Add(m.ttl),
	}
	m.apps[app.ID] = app
	return app.ID
}

// Vote records a signed vote on an application.
// voterAddr must match pubKey and the signature must cover "id:approve".
func (m *AuthorityApplicationManager) Vote(voterAddr, id string, approve bool, sig []byte, pubKey ed25519.PublicKey) error {
	if hex.EncodeToString(pubKey) != voterAddr {
		return errors.New("voter address mismatch")
	}
	msg := fmt.Sprintf("%s:%t", id, approve)
	if !ed25519.Verify(pubKey, []byte(msg), sig) {
		return errors.New("invalid signature")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	app, ok := m.apps[id]
	if !ok {
		return errors.New("application not found")
	}
	if app.Finalized || time.Now().After(app.ExpiresAt) {
		return errors.New("application closed")
	}
	if approve {
		app.Approvals[voterAddr] = true
		delete(app.Rejections, voterAddr)
	} else {
		app.Rejections[voterAddr] = true
		delete(app.Approvals, voterAddr)
	}
	return nil
}

// Finalize concludes an application and registers the node if approved.
func (m *AuthorityApplicationManager) Finalize(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	app, ok := m.apps[id]
	if !ok {
		return errors.New("application not found")
	}
	if app.Finalized {
		return errors.New("already finalised")
	}
	app.Finalized = true
	approved := len(app.Approvals) > len(app.Rejections)
	if approved {
		if _, err := m.registry.Register(app.Candidate, app.Role); err != nil {
			return err
		}
	}
	return nil
}

// Tick removes expired applications.
func (m *AuthorityApplicationManager) Tick(now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, app := range m.apps {
		if app.Finalized || now.After(app.ExpiresAt) {
			delete(m.apps, id)
		}
	}
}

// Get returns an application by ID.
func (m *AuthorityApplicationManager) Get(id string) (*AuthorityApplication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	app, ok := m.apps[id]
	if !ok {
		return nil, errors.New("application not found")
	}
	return app, nil
}

// List returns all applications.
func (m *AuthorityApplicationManager) List() []*AuthorityApplication {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*AuthorityApplication, 0, len(m.apps))
	for _, app := range m.apps {
		out = append(out, app)
	}
	return out
}

// MarshalJSON provides deterministic output for CLI and GUI integration.
func (a *AuthorityApplication) MarshalJSON() ([]byte, error) {
	type alias AuthorityApplication
	return json.Marshal(&struct {
		Approvals  int `json:"approvals"`
		Rejections int `json:"rejections"`
		*alias
	}{
		Approvals:  len(a.Approvals),
		Rejections: len(a.Rejections),
		alias:      (*alias)(a),
	})
}
