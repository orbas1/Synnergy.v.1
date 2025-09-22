package synnergy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"sort"
	"sync"
	"time"
)

const aes256KeySize = 32

var (
	ErrInvalidKeyLength     = errors.New("private tx: key must be 32 bytes")
	ErrInvalidNonce         = errors.New("private tx: nonce size mismatch")
	ErrMissingTransactionID = errors.New("private tx: id required")
	ErrDuplicateTransaction = errors.New("private tx: duplicate identifier")
	ErrSignatureMismatch    = errors.New("private tx: signature verification failed")
)

// Envelope captures the raw encrypted payload produced by EncryptWithAAD.  It is
// useful for persisting data alongside metadata (such as sender/recipient) so
// that replay protection and auditing can be layered on top by callers.
type Envelope struct {
	Nonce          []byte
	Ciphertext     []byte
	AssociatedData []byte
}

func cloneBytes(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func newGCM(key []byte) (cipher.AEAD, error) {
	if len(key) != aes256KeySize {
		return nil, ErrInvalidKeyLength
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// Encrypt encrypts plaintext using AES-GCM with the provided key.  The returned
// slice contains nonce||ciphertext to remain backward compatible with earlier
// CLI tooling.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	env, err := EncryptWithAAD(key, plaintext, nil)
	if err != nil {
		return nil, err
	}
	combined := make([]byte, len(env.Nonce)+len(env.Ciphertext))
	copy(combined, env.Nonce)
	copy(combined[len(env.Nonce):], env.Ciphertext)
	return combined, nil
}

// Decrypt decrypts data produced by Encrypt.
func Decrypt(key, data []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, io.ErrUnexpectedEOF
	}
	env := Envelope{
		Nonce:      cloneBytes(data[:nonceSize]),
		Ciphertext: cloneBytes(data[nonceSize:]),
	}
	return decryptWithGCM(gcm, env)
}

// EncryptWithAAD encrypts plaintext and binds the supplied associated data to
// the ciphertext.  The associated data is not encrypted but is authenticated by
// AES-GCM, allowing consensus and regulatory layers to assert metadata integrity.
func EncryptWithAAD(key, plaintext, aad []byte) (Envelope, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return Envelope{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return Envelope{}, err
	}
	ciphertext := gcm.Seal(nil, nonce, plaintext, aad)
	return Envelope{
		Nonce:          nonce,
		Ciphertext:     ciphertext,
		AssociatedData: cloneBytes(aad),
	}, nil
}

// DecryptWithAAD decrypts an Envelope produced by EncryptWithAAD.
func DecryptWithAAD(key []byte, env Envelope) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	return decryptWithGCM(gcm, env)
}

func decryptWithGCM(gcm cipher.AEAD, env Envelope) ([]byte, error) {
	if len(env.Nonce) != gcm.NonceSize() {
		return nil, ErrInvalidNonce
	}
	return gcm.Open(nil, env.Nonce, env.Ciphertext, env.AssociatedData)
}

// PrivateTransaction holds an encrypted payload along with metadata required for
// regulatory reviews and replay protection.
type PrivateTransaction struct {
	ID             string
	Sender         string
	Recipient      string
	Payload        []byte
	Nonce          []byte
	AssociatedData []byte
	Signature      []byte
	Timestamp      time.Time
}

// Envelope returns a copy of the encrypted envelope data for the transaction.
func (pt PrivateTransaction) Envelope() Envelope {
	return Envelope{
		Nonce:          cloneBytes(pt.Nonce),
		Ciphertext:     cloneBytes(pt.Payload),
		AssociatedData: cloneBytes(pt.AssociatedData),
	}
}

// ApplyEnvelope stores the encrypted envelope on the transaction.  It is used by
// helpers that manage encryption lifecycle separately from metadata enrichment.
func (pt *PrivateTransaction) ApplyEnvelope(env Envelope) {
	if pt == nil {
		return
	}
	pt.Nonce = cloneBytes(env.Nonce)
	pt.Payload = cloneBytes(env.Ciphertext)
	pt.AssociatedData = cloneBytes(env.AssociatedData)
}

// Clone returns a deep copy of the transaction.
func (pt PrivateTransaction) Clone() PrivateTransaction {
	dup := pt
	dup.Payload = cloneBytes(pt.Payload)
	dup.Nonce = cloneBytes(pt.Nonce)
	dup.AssociatedData = cloneBytes(pt.AssociatedData)
	dup.Signature = cloneBytes(pt.Signature)
	return dup
}

// Digest returns a deterministic hash of the transaction metadata and envelope
// contents.  Signatures are computed over this digest to protect against replay
// attacks and metadata tampering.
func (pt PrivateTransaction) Digest() []byte {
	h := sha256.New()
	h.Write([]byte(pt.ID))
	h.Write([]byte(pt.Sender))
	h.Write([]byte(pt.Recipient))
	h.Write(pt.Payload)
	h.Write(pt.Nonce)
	h.Write(pt.AssociatedData)
	ts := pt.Timestamp.UTC().UnixNano()
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(ts))
	h.Write(buf)
	return h.Sum(nil)
}

// SignTransaction signs the transaction digest with the provided Ed25519 private
// key.  The signature is stored in the transaction for later verification.
func SignTransaction(priv ed25519.PrivateKey, tx *PrivateTransaction) error {
	if tx == nil {
		return errors.New("private tx: transaction required")
	}
	if len(priv) != ed25519.PrivateKeySize {
		return errors.New("private tx: invalid ed25519 private key")
	}
	sig := ed25519.Sign(priv, tx.Digest())
	tx.Signature = cloneBytes(sig)
	return nil
}

// VerifyTransaction checks the embedded signature using the supplied public key.
func VerifyTransaction(pub ed25519.PublicKey, tx PrivateTransaction) error {
	if len(pub) != ed25519.PublicKeySize {
		return errors.New("private tx: invalid ed25519 public key")
	}
	if len(tx.Signature) == 0 {
		return ErrSignatureMismatch
	}
	if !ed25519.Verify(pub, tx.Digest(), tx.Signature) {
		return ErrSignatureMismatch
	}
	return nil
}

func deriveTransactionID(tx PrivateTransaction) string {
	h := sha256.New()
	h.Write(tx.Payload)
	h.Write(tx.Nonce)
	h.Write(tx.AssociatedData)
	return hex.EncodeToString(h.Sum(nil))
}

// PrivateTxManager manages private transactions.  It guarantees that stored
// transactions retain insertion order for auditability while providing O(1)
// lookups for CLI and web layers.
type PrivateTxManager struct {
	mu    sync.RWMutex
	txs   map[string]PrivateTransaction
	order []string
}

// NewPrivateTxManager creates a new PrivateTxManager.
func NewPrivateTxManager() *PrivateTxManager {
	return &PrivateTxManager{txs: make(map[string]PrivateTransaction)}
}

// Send adds a private transaction to the internal pool.  It falls back to
// generating an identifier when one is not provided so legacy flows continue to
// function.
func (m *PrivateTxManager) Send(tx PrivateTransaction) {
	_ = m.Store(tx)
}

// Store inserts a transaction and returns an error when duplicates are
// encountered.  Callers that require deterministic behaviour should prefer this
// helper over Send.
func (m *PrivateTxManager) Store(tx PrivateTransaction) error {
	if tx.ID == "" {
		tx.ID = deriveTransactionID(tx)
	}
	if tx.ID == "" {
		return ErrMissingTransactionID
	}
	if tx.Timestamp.IsZero() {
		tx.Timestamp = time.Now().UTC()
	}
	clone := tx.Clone()

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.txs[clone.ID]; exists {
		return ErrDuplicateTransaction
	}
	m.txs[clone.ID] = clone
	m.order = append(m.order, clone.ID)
	return nil
}

// Upsert inserts or replaces a transaction with the same identifier.
func (m *PrivateTxManager) Upsert(tx PrivateTransaction) {
	if tx.ID == "" {
		tx.ID = deriveTransactionID(tx)
	}
	if tx.Timestamp.IsZero() {
		tx.Timestamp = time.Now().UTC()
	}
	clone := tx.Clone()

	m.mu.Lock()
	if _, exists := m.txs[clone.ID]; !exists {
		m.order = append(m.order, clone.ID)
	}
	m.txs[clone.ID] = clone
	m.mu.Unlock()
}

// Get returns a transaction by identifier.
func (m *PrivateTxManager) Get(id string) (PrivateTransaction, bool) {
	m.mu.RLock()
	tx, ok := m.txs[id]
	m.mu.RUnlock()
	if !ok {
		return PrivateTransaction{}, false
	}
	return tx.Clone(), true
}

// List returns a copy of stored private transactions ordered by timestamp.
func (m *PrivateTxManager) List() []PrivateTransaction {
	m.mu.RLock()
	items := make([]PrivateTransaction, 0, len(m.txs))
	for _, id := range m.order {
		tx, ok := m.txs[id]
		if !ok {
			continue
		}
		items = append(items, tx.Clone())
	}
	m.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Timestamp.Before(items[j].Timestamp)
	})
	return items
}
