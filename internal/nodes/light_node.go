package nodes

// BlockHeader contains minimal block information used by light nodes.
type BlockHeader struct {
	Hash       string
	Height     uint64
	ParentHash string
}
