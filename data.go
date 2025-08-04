package synnergy

// Opcode identifiers for the content distribution network (CDN) module.
// These values are reserved for operations that manage large content blobs
// referenced on-chain. They can be used by higher level components or the
// virtual machine when encoding actions.
const (
	// OpDataAdd registers new content metadata with the network.
	OpDataAdd uint8 = iota + 0x30
	// OpDataRemove removes content from the distribution index.
	OpDataRemove
	// OpDataPin marks content as pinned so that nodes retain it locally.
	OpDataPin
)
