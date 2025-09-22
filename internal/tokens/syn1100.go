package tokens

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
)

// HealthRecord represents encrypted healthcare data tied to an owner.
type HealthRecord struct {
	Owner       string
	Ciphertext  []byte
	Nonce       []byte
	Access      map[string]bool
	LastUpdated time.Time
}

// HealthAuditEvent records traceable record operations.
type HealthAuditEvent struct {
	RecordID TokenID
	Actor    string
	Action   string
	Time     time.Time
}

// SYN1100Token manages healthcare records with access control.
type SYN1100Token struct {
	mu      sync.RWMutex
	records map[TokenID]*HealthRecord
	key     []byte
	audit   []HealthAuditEvent
}

// NewSYN1100Token creates an empty healthcare record store.
func NewSYN1100Token() *SYN1100Token {
	return &SYN1100Token{records: make(map[TokenID]*HealthRecord)}
}

// SetEncryptionKey configures the symmetric key used to encrypt all records.
func (t *SYN1100Token) SetEncryptionKey(key []byte) error {
	if len(key) != 32 {
		return fmt.Errorf("invalid key length: %d", len(key))
	}
	t.mu.Lock()
	t.key = make([]byte, len(key))
	copy(t.key, key)
	t.mu.Unlock()
	return nil
}

// encrypt encrypts plain data using AES-GCM.
func (t *SYN1100Token) encrypt(plain []byte) ([]byte, []byte, error) {
	if len(t.key) == 0 {
		return nil, nil, errors.New("encryption key not set")
	}
	block, err := aes.NewCipher(t.key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}
	ciphertext := gcm.Seal(nil, nonce, plain, nil)
	return ciphertext, nonce, nil
}

// decrypt decrypts ciphertext using AES-GCM.
func (t *SYN1100Token) decrypt(ciphertext, nonce []byte) ([]byte, error) {
	if len(t.key) == 0 {
		return nil, errors.New("encryption key not set")
	}
	block, err := aes.NewCipher(t.key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// AddRecord stores a new healthcare record with the given ID.
func (t *SYN1100Token) AddRecord(id TokenID, owner string, data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.records[id]; exists {
		return fmt.Errorf("record exists")
	}
	cipher, nonce, err := t.encrypt(data)
	if err != nil {
		return err
	}
	t.records[id] = &HealthRecord{
		Owner:       owner,
		Ciphertext:  cipher,
		Nonce:       nonce,
		Access:      map[string]bool{owner: true},
		LastUpdated: time.Now(),
	}
	t.audit = append(t.audit, HealthAuditEvent{RecordID: id, Actor: owner, Action: "create", Time: time.Now()})
	return nil
}

// UpdateRecord replaces the stored healthcare data.
func (t *SYN1100Token) UpdateRecord(id TokenID, actor string, data []byte) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	rec, ok := t.records[id]
	if !ok {
		return fmt.Errorf("record not found")
	}
	if actor != rec.Owner && !rec.Access[actor] {
		return fmt.Errorf("access denied")
	}
	cipher, nonce, err := t.encrypt(data)
	if err != nil {
		return err
	}
	rec.Ciphertext = cipher
	rec.Nonce = nonce
	rec.LastUpdated = time.Now()
	t.audit = append(t.audit, HealthAuditEvent{RecordID: id, Actor: actor, Action: "update", Time: rec.LastUpdated})
	return nil
}

// GrantAccess allows a grantee to read the specified record.
func (t *SYN1100Token) GrantAccess(id TokenID, grantee string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	rec, ok := t.records[id]
	if !ok {
		return fmt.Errorf("record not found")
	}
	rec.Access[grantee] = true
	t.audit = append(t.audit, HealthAuditEvent{RecordID: id, Actor: grantee, Action: "grant", Time: time.Now()})
	return nil
}

// RevokeAccess revokes a previously granted permission.
func (t *SYN1100Token) RevokeAccess(id TokenID, grantee string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	rec, ok := t.records[id]
	if !ok {
		return fmt.Errorf("record not found")
	}
	delete(rec.Access, grantee)
	t.audit = append(t.audit, HealthAuditEvent{RecordID: id, Actor: grantee, Action: "revoke", Time: time.Now()})
	return nil
}

// GetRecord returns the decrypted record data if the caller has access rights.
func (t *SYN1100Token) GetRecord(id TokenID, caller string) ([]byte, error) {
	t.mu.RLock()
	rec, ok := t.records[id]
	t.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("record not found")
	}
	if !rec.Access[caller] {
		return nil, fmt.Errorf("access denied")
	}
	plain, err := t.decrypt(rec.Ciphertext, rec.Nonce)
	if err != nil {
		return nil, err
	}
	return plain, nil
}

// AuditTrail returns the tail of the audit log for monitoring.
func (t *SYN1100Token) AuditTrail(limit int) []HealthAuditEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if limit <= 0 || limit >= len(t.audit) {
		out := make([]HealthAuditEvent, len(t.audit))
		copy(out, t.audit)
		return out
	}
	out := make([]HealthAuditEvent, limit)
	copy(out, t.audit[len(t.audit)-limit:])
	return out
}
