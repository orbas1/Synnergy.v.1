package nodes

// BlockHeader contains minimal block information used by light nodes.
type BlockHeader struct {
	Hash       string
	Height     uint64
	ParentHash string
}

// LightNode implements a minimal client that tracks block headers and embeds
// the behaviour of a BasicNode.  It can be used by wallets or monitoring tools
// that do not require full blockchain state.
type LightNode struct {
	*BasicNode
	headers []BlockHeader
}

// NewLightNode returns a new light node with the provided identifier.
func NewLightNode(id Address) *LightNode {
	return &LightNode{BasicNode: NewBasicNode(id)}
}

// AddHeader records a new block header.
func (n *LightNode) AddHeader(h BlockHeader) {
	n.headers = append(n.headers, h)
}

// Headers returns a copy of all headers known to the node.
func (n *LightNode) Headers() []BlockHeader {
	cp := make([]BlockHeader, len(n.headers))
	copy(cp, n.headers)
	return cp
}

// LatestHeader returns the most recently added header.  The boolean return
// value is false if the node has not observed any blocks yet.
func (n *LightNode) LatestHeader() (BlockHeader, bool) {
	if len(n.headers) == 0 {
		return BlockHeader{}, false
	}
	return n.headers[len(n.headers)-1], true
}
