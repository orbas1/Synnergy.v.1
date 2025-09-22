package core

import (
	"context"
	"errors"

	ilog "synnergy/internal/log"
)

// GenesisStats contains summary information about the chain after genesis
// block creation.
type GenesisStats struct {
	Height      uint64
	Hash        string
	Circulating uint64
	Remaining   uint64
	Weights     ConsensusWeights
}

// InitGenesis creates the genesis block using the node's Synnergy consensus.
// It credits the creator wallet with the GenesisAllocation and mines the first
// block. An error is returned if a block already exists.
func (n *Node) InitGenesis(wallets GenesisWallets) (GenesisStats, *Block, error) {
	if len(n.Blockchain) != 0 {
		return GenesisStats{}, nil, errors.New("genesis already exists")
	}
	if err := wallets.Validate(); err != nil {
		return GenesisStats{}, nil, err
	}
	n.Ledger.Credit(wallets.CreatorWallet, GenesisAllocation)
	sb := NewSubBlock(nil, wallets.Genesis)
	block := NewBlock([]*SubBlock{sb}, "")
	if err := n.Consensus.MineBlock(context.Background(), block, 1); err != nil {
		return GenesisStats{}, nil, err
	}
	n.Blockchain = append(n.Blockchain, block)
	if err := n.Ledger.AddBlock(block); err != nil {
		return GenesisStats{}, nil, err
	}
	stats := GenesisStats{
		Height:      1,
		Hash:        block.Hash,
		Circulating: CirculatingSupply(0),
		Remaining:   RemainingSupply(0),
		Weights:     n.Consensus.Weights,
	}
	ilog.Info("genesis_init", "hash", block.Hash, "circulating", stats.Circulating, "remaining", stats.Remaining)
	return stats, block, nil
}
