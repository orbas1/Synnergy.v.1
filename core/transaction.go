package core

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	// ErrZeroAmount is returned when a transaction attempts to move no funds.
	ErrZeroAmount = errors.New("transaction amount must be greater than zero")
	// ErrSelfTransfer is returned when the sender and recipient are identical.
	ErrSelfTransfer = errors.New("sender and recipient must differ")
	// ErrTransactionIDMismatch is returned when the cached ID does not match the
	// computed hash. This guards against tampering or partial updates.
	ErrTransactionIDMismatch = errors.New("transaction id mismatch")
	// ErrTransactionOverflow signals that amount and fee addition overflowed.
	ErrTransactionOverflow = errors.New("transaction arithmetic overflow")
	// ErrTransactionClockSkew is raised when the timestamp deviates beyond the
	// permitted drift window.
	ErrTransactionClockSkew = errors.New("transaction timestamp outside permitted skew window")
)

// DefaultMaxClockSkew constrains transaction timestamps to a sane window so
// clients with misconfigured clocks cannot flood the network with stale or
// future-dated payloads.
const DefaultMaxClockSkew = 5 * time.Minute

// TransactionValidationConfig configures validation constraints applied to a
// transaction prior to execution.
type TransactionValidationConfig struct {
	// Now represents the validator's notion of current time. If zero the
	// timestamp check is skipped.
	Now time.Time
	// MaxClockSkew bounds how far the transaction timestamp can drift from
	// Now. A zero duration disables the check.
	MaxClockSkew time.Duration
	// MinFee enforces a network minimum fee. Zero disables the check.
	MinFee uint64
	// MaxFee bounds the total fee a transaction may declare. Zero disables
	// the check.
	MaxFee uint64
	// AllowSelfTransfer controls whether from and to may be identical.
	AllowSelfTransfer bool
}

// DefaultTransactionValidationConfig returns the standard validation policy.
func DefaultTransactionValidationConfig() TransactionValidationConfig {
	return TransactionValidationConfig{
		Now:               time.Now(),
		MaxClockSkew:      DefaultMaxClockSkew,
		MinFee:            0,
		MaxFee:            0,
		AllowSelfTransfer: false,
	}
}

// Transaction represents a transfer of Synthron between accounts.
//
// Transactions are signed payloads that move coins from one address to
// another.  They include a fee and timestamp so that they can be ordered
// deterministically by the consensus engine.
type Transaction struct {
	ID        string
	From      string
	To        string
	Amount    uint64
	Fee       uint64
	Nonce     uint64
	Timestamp int64
	Signature []byte
	Type      TransactionType
	// BiometricHash stores the hash of the biometric data used to
	// authorize this transaction. It ensures that the transaction is tied
	// to a verified identity when biometric authentication is required.
	BiometricHash []byte

	// Program holds optional bytecode instructions to be executed by the
	// Synnergy Virtual Machine.  Traditional transfer transactions will
	// leave this field nil.
	Program []Instruction
}

// NewTransaction creates a new unsigned transaction with the provided
// parameters.  The ID is derived from a hash of the core fields and can be
// reproduced deterministically prior to signing.  Transactions default to the
// Transfer type unless modified by higher level logic.
func NewTransaction(from, to string, amount, fee, nonce uint64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Fee:       fee,
		Nonce:     nonce,
		Timestamp: time.Now().Unix(),
		Type:      TxTypeTransfer,
	}
	tx.ID = tx.Hash()
	return tx
}

// Hash returns the hex-encoded hash of the transaction contents excluding the
// signature.  It is used as the message for signing and verification.
func (t *Transaction) Hash() string {
	bio := hex.EncodeToString(t.BiometricHash)
	h := sha256.Sum256([]byte(fmt.Sprintf(
		"%s%s%d%d%d%d%d%s%s",
		t.From,
		t.To,
		t.Amount,
		t.Fee,
		t.Nonce,
		t.Timestamp,
		t.Type,
		bio,
		t.programDigest(),
	)))
	return hex.EncodeToString(h[:])
}

func (t *Transaction) programDigest() string {
	if len(t.Program) == 0 {
		return ""
	}
	payload := make([]byte, 0, len(t.Program)*16)
	for _, inst := range t.Program {
		payload = append(payload, []byte(fmt.Sprintf("%d:%d;", inst.Op, inst.Value))...)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:])
}

// Verify checks the transaction's signature against the provided public key.
func (t *Transaction) Verify(pub *ecdsa.PublicKey) bool {
	return VerifySignature(t, t.Signature, pub)
}

// AttachBiometric verifies the provided biometric data for the given user
// through the supplied BiometricService. If verification succeeds the biometric
// hash is attached to the transaction and the transaction ID is recalculated to
// include the biometric. This prevents replay or tampering with biometric data
// after signing.
func (t *Transaction) AttachBiometric(userID string, biometric []byte, sig []byte, svc *BiometricService) error {
	if svc == nil {
		return errors.New("biometric service not available")
	}
	if !svc.Verify(userID, biometric, sig) {
		return errors.New("biometric verification failed")
	}
	h := sha256.Sum256(biometric)
	t.BiometricHash = h[:]
	t.ID = t.Hash()
	return nil
}

// Clone performs a deep copy of the transaction and recalculates the hash so
// downstream components can mutate the copy without affecting the original.
func (t *Transaction) Clone() *Transaction {
	if t == nil {
		return nil
	}
	cp := *t
	if len(t.Signature) > 0 {
		cp.Signature = append([]byte(nil), t.Signature...)
	}
	if len(t.BiometricHash) > 0 {
		cp.BiometricHash = append([]byte(nil), t.BiometricHash...)
	}
	if len(t.Program) > 0 {
		cp.Program = append([]Instruction(nil), t.Program...)
	}
	cp.ID = cp.Hash()
	return &cp
}

// TotalCost returns amount+fee, guarding against overflow.
func (t *Transaction) TotalCost() (uint64, error) {
	if t.Amount > math.MaxUint64-t.Fee {
		return 0, ErrTransactionOverflow
	}
	return t.Amount + t.Fee, nil
}

// ValidateBasic applies structural validation using the provided config. It is
// safe to call with a zero config which applies the defaults.
func (t *Transaction) ValidateBasic(cfgs ...TransactionValidationConfig) error {
	cfg := DefaultTransactionValidationConfig()
	if len(cfgs) > 0 {
		cfg = cfgs[0]
		if cfg.Now.IsZero() {
			cfg.Now = time.Now()
		}
	}
	if t == nil {
		return ErrNilTransaction
	}
	if !cfg.AllowSelfTransfer && t.From == t.To {
		return ErrSelfTransfer
	}
	if t.From == "" || t.To == "" {
		return ErrEmptyAddress
	}
	if t.Amount == 0 {
		return ErrZeroAmount
	}
	if cfg.MinFee > 0 && t.Fee < cfg.MinFee {
		return fmt.Errorf("fee below minimum %d", cfg.MinFee)
	}
	if cfg.MaxFee > 0 && t.Fee > cfg.MaxFee {
		return fmt.Errorf("fee exceeds maximum %d", cfg.MaxFee)
	}
	if _, err := t.TotalCost(); err != nil {
		return err
	}
	if cfg.MaxClockSkew > 0 && cfg.Now.Unix() != 0 {
		ts := time.Unix(t.Timestamp, 0)
		if ts.Before(cfg.Now.Add(-cfg.MaxClockSkew)) || ts.After(cfg.Now.Add(cfg.MaxClockSkew)) {
			return ErrTransactionClockSkew
		}
	}
	if expected := t.Hash(); t.ID == "" {
		t.ID = expected
	} else if expected != t.ID {
		return ErrTransactionIDMismatch
	}
	for idx, inst := range t.Program {
		if err := inst.Validate(); err != nil {
			return fmt.Errorf("program[%d]: %w", idx, err)
		}
	}
	return nil
}
