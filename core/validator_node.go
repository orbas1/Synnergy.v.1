package core

import "context"

// ValidatorNode bundles a base node with validator management and quorum tracking.
type ValidatorNode struct {
	*Node
	Validators *ValidatorManager
	Quorum     *QuorumTracker
}

// NewValidatorNode creates a validator node with the given minimum stake and
// quorum requirement.
func NewValidatorNode(id, addr string, ledger *Ledger, minStake uint64, quorum int) *ValidatorNode {
	vn := &ValidatorNode{
		Node:       NewNode(id, addr, ledger),
		Validators: NewValidatorManager(minStake),
		Quorum:     NewQuorumTracker(quorum),
	}
	return vn
}

// AddValidator registers a validator within the node's manager.
func (vn *ValidatorNode) AddValidator(addr string, stake uint64) error {
	if err := vn.Validators.Add(context.Background(), addr, stake); err != nil {
		return err
	}
	vn.Quorum.Join(addr)
	return nil
}

// RemoveValidator removes a validator from the set and quorum tracker.
func (vn *ValidatorNode) RemoveValidator(addr string) {
	vn.Validators.Remove(context.Background(), addr)
	vn.Quorum.Leave(addr)
}

// SlashValidator penalises a validator and removes it from quorum if needed.
func (vn *ValidatorNode) SlashValidator(addr string) {
	vn.Validators.Slash(context.Background(), addr)
	vn.Quorum.Leave(addr)
}

// HasQuorum reports whether the active validator set meets the quorum
// requirement.
func (vn *ValidatorNode) HasQuorum() bool {
	return vn.Quorum.Reached()
}
