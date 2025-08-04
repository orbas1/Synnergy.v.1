package nodes

// ConsensusNodeInterface defines behaviour for nodes dedicated to a single
// consensus algorithm. Implementations typically optimise networking and
// validation logic for the returned consensus type while still satisfying the
// base NodeInterface.
type ConsensusNodeInterface interface {
	NodeInterface
	// ConsensusType reports the consensus mechanism this node is optimised for,
	// such as "pow" or "pos".
	ConsensusType() string
}
