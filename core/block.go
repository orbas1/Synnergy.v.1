package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// SubBlock contains transactions ordered by PoH and validated by PoS.
type SubBlock struct {
	Transactions []*Transaction
	Validator    string
	PohHash      string
	Timestamp    int64
	Signature    string
}

// NewSubBlock constructs a sub-block from the given transactions and validator.
func NewSubBlock(txs []*Transaction, validator string) *SubBlock {
	sb := &SubBlock{Transactions: txs, Validator: validator, Timestamp: time.Now().Unix()}
	sb.PohHash = sb.Hash()
	sb.Signature = signSubBlock(validator, sb.PohHash)
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
	expected := signSubBlock(sb.Validator, sb.PohHash)
	return sb.Signature == expected
}

// Block aggregates validated sub-blocks and is finalized via PoW.
type Block struct {
	SubBlocks []*SubBlock
	PrevHash  string
	Nonce     uint64
	Timestamp int64
	Hash      string
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

func signSubBlock(validator, msg string) string {
	h := sha256.Sum256([]byte(validator + msg))
	return hex.EncodeToString(h[:])
}
