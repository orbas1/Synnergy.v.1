package core

// SubBlock groups transactions validated via POS and POH.
type SubBlock struct {
	Transactions []*Transaction
	Validator    string
}

// Block is composed of multiple SubBlocks and finalized via POW.
type Block struct {
	SubBlocks []SubBlock
	Nonce     uint64
}

// SynnergyConsensus encapsulates the consensus algorithms.
type SynnergyConsensus struct{}

// NewSynnergyConsensus returns a new consensus engine.
func NewSynnergyConsensus() *SynnergyConsensus { return &SynnergyConsensus{} }

// ValidateSubBlock performs POS and POH validation on a sub-block.
func (sc *SynnergyConsensus) ValidateSubBlock(sb *SubBlock) bool {
	// TODO: implement POS and POH
	return true
}

// MineBlock performs POW to finalize a block.
func (sc *SynnergyConsensus) MineBlock(b *Block) {
	// TODO: implement POW mining
}
