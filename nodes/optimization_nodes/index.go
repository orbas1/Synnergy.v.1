package optimizationnodes

// Transaction represents the minimal transaction information required for
// optimisation decisions. Only fields relevant to ordering are included.
type Transaction struct {
	Hash string // unique transaction identifier
	Fee  uint64 // total fee offered by the transaction
	Size int    // approximate size in bytes
}
