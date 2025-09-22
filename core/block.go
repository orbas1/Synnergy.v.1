package core

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
)

// SubBlock contains transactions ordered by PoH and validated by PoS.
type SubBlock struct {
	Transactions []*Transaction
	Validator    string
	PohHash      string
	Timestamp    int64
	Signature    []byte
	ValidatorKey []byte
	System       bool
}

const maxTimeDriftSeconds = 300 // five minutes

// NewSubBlock constructs a sub-block from the given transactions and validator.
func NewSubBlock(txs []*Transaction, validator string) *SubBlock {
	sb := &SubBlock{Transactions: txs, Validator: validator, Timestamp: time.Now().Unix()}
	sb.PohHash = sb.Hash()
	if err := SignSubBlock(sb); err != nil {
		sb.Signature = nil
	}
	return sb
}

// NewGenesisSubBlock constructs a system-signed sub-block used exclusively for
// the genesis block. The sub-block contains no transactions and is marked so
// that validation bypasses signature checks while still providing a PoH link.
func NewGenesisSubBlock(validator string) *SubBlock {
	sb := &SubBlock{Validator: validator, Timestamp: time.Now().Unix(), System: true}
	sb.PohHash = sb.Hash()
	return sb
}

// Hash generates a deterministic hash of the sub-block's contents to provide the
// PoH link.  For now a simple concatenation of transaction IDs is used.
func (sb *SubBlock) Hash() string {
	h := sha256.New()
	for _, tx := range sb.Transactions {
		h.Write([]byte(tx.ID))
	}
	h.Write([]byte(sb.Validator))
	h.Write([]byte(fmt.Sprintf("%d", sb.Timestamp)))
	return hex.EncodeToString(h.Sum(nil))
}

// VerifySignature confirms the sub-block was signed by the stated validator.
func (sb *SubBlock) VerifySignature() bool {
	if sb.System {
		return len(sb.Transactions) == 0
	}
	if len(sb.Signature) == 0 || len(sb.ValidatorKey) == 0 {
		return false
	}
	pub, err := decodePublicKey(sb.ValidatorKey)
	if err != nil {
		return false
	}
	if deriveAddress(pub) != sb.Validator {
		return false
	}
	digest, err := hex.DecodeString(sb.PohHash)
	if err != nil {
		return false
	}
	mid := len(sb.Signature) / 2
	if mid == 0 {
		return false
	}
	r := new(big.Int).SetBytes(sb.Signature[:mid])
	s := new(big.Int).SetBytes(sb.Signature[mid:])
	return ecdsa.Verify(pub, digest, r, s)
}

// Validate checks the internal consistency of the sub-block and ensures it was
// signed by the declared validator. It returns an error if any field is
// malformed or tampered with.
func (sb *SubBlock) Validate() error {
	if len(sb.Transactions) == 0 {
		if !sb.System {
			return fmt.Errorf("no transactions")
		}
	}
	seen := make(map[string]struct{})
	for _, tx := range sb.Transactions {
		if tx == nil {
			return fmt.Errorf("nil transaction")
		}
		if _, ok := seen[tx.ID]; ok {
			return fmt.Errorf("duplicate transaction %s", tx.ID)
		}
		seen[tx.ID] = struct{}{}
	}
	if !sb.System && sb.Validator == "" {
		return fmt.Errorf("validator required")
	}
	now := time.Now().Unix()
	if sb.Timestamp == 0 {
		return fmt.Errorf("timestamp required")
	}
	if sb.Timestamp > now+maxTimeDriftSeconds {
		return fmt.Errorf("timestamp in future")
	}
	if sb.PohHash != sb.Hash() {
		return fmt.Errorf("poh hash mismatch")
	}
	if !sb.System {
		if len(sb.Signature) == 0 {
			return fmt.Errorf("signature required")
		}
		if len(sb.ValidatorKey) == 0 {
			return fmt.Errorf("validator key required")
		}
		if !sb.VerifySignature() {
			return fmt.Errorf("invalid signature")
		}
	}
	return nil
}

// Block aggregates validated sub-blocks and is finalized via PoW.
type Block struct {
	SubBlocks []*SubBlock
	PrevHash  string
	Nonce     uint64
	Timestamp int64
	Hash      string
	Finalized bool
}

// NewBlock creates a block from sub-blocks and the hash of the previous block.
func NewBlock(subBlocks []*SubBlock, prevHash string) *Block {
	return &Block{SubBlocks: subBlocks, PrevHash: prevHash, Timestamp: time.Now().Unix()}
}

// HeaderHash returns the hash of the block header for a given nonce.  This is
// used as the proof-of-work target.
func (b *Block) HeaderHash(nonce uint64) string {
	h := sha256.New()
	h.Write([]byte(b.PrevHash))
	for _, sb := range b.SubBlocks {
		h.Write([]byte(sb.PohHash))
	}
	h.Write([]byte(fmt.Sprintf("%d%d", b.Timestamp, nonce)))
	return hex.EncodeToString(h.Sum(nil))
}

// SignSubBlock looks up the validator's wallet from the global registry and
// signs the sub-block's PoH hash. If no key is registered, an error is
// returned. The helper is exported so CLI tooling can opportunistically sign
// sub-blocks after loading a wallet file.
func SignSubBlock(sb *SubBlock) error {
	if sb == nil {
		return errors.New("sub-block required")
	}
	if sb.System {
		return nil
	}
	priv := validatorPrivateKey(sb.Validator)
	if priv == nil {
		return fmt.Errorf("no signing key for validator %s", sb.Validator)
	}
	sig, key, err := signSubBlock(priv, sb.PohHash)
	if err != nil {
		return err
	}
	sb.Signature = sig
	sb.ValidatorKey = key
	return nil
}

func signSubBlock(priv *ecdsa.PrivateKey, msg string) ([]byte, []byte, error) {
	if priv == nil {
		return nil, nil, errors.New("private key required")
	}
	digest, err := hex.DecodeString(msg)
	if err != nil {
		return nil, nil, err
	}
	r, s, err := ecdsa.Sign(rand.Reader, priv, digest)
	if err != nil {
		return nil, nil, err
	}
	sig := make([]byte, 64)
	rb := r.Bytes()
	sb := s.Bytes()
	copy(sig[32-len(rb):32], rb)
	copy(sig[64-len(sb):], sb)
	pub := encodePublicKey(&priv.PublicKey)
	return sig, pub, nil
}

// Validate checks that the block and its sub-blocks are internally consistent.
// For non-genesis blocks it also verifies the stored header hash matches the
// computed hash for the provided nonce.
func (b *Block) Validate() error {
	if len(b.SubBlocks) == 0 {
		return fmt.Errorf("no sub-blocks")
	}
	now := time.Now().Unix()
	if b.Timestamp == 0 {
		return fmt.Errorf("timestamp required")
	}
	if b.Timestamp > now+maxTimeDriftSeconds {
		return fmt.Errorf("timestamp in future")
	}
	for _, sb := range b.SubBlocks {
		if err := sb.Validate(); err != nil {
			return fmt.Errorf("sub-block invalid: %w", err)
		}
		if sb.Timestamp > b.Timestamp {
			return fmt.Errorf("sub-block timestamp after block timestamp")
		}
	}
	if b.PrevHash != "" {
		if b.Hash == "" {
			return fmt.Errorf("hash required")
		}
		if b.Hash != b.HeaderHash(b.Nonce) {
			return fmt.Errorf("hash mismatch")
		}
	}
	return nil
}
