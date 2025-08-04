package core

import (
        "bytes"
        "compress/gzip"
        "encoding/json"
        "os"
)

// ledgerSnapshot is a helper type used for serializing the ledger. It exposes
// only the fields that need to be persisted.
type ledgerSnapshot struct {
        Balances map[string]uint64 `json:"balances"`
        Blocks   []*Block          `json:"blocks"`
}

// CompressLedger returns the gzip-compressed JSON encoding of the provided ledger.
func CompressLedger(l *Ledger) ([]byte, error) {
        l.mu.RLock()
        defer l.mu.RUnlock()
        snap := ledgerSnapshot{Balances: l.balances, Blocks: l.blocks}
        var buf bytes.Buffer
        gz := gzip.NewWriter(&buf)
        if err := json.NewEncoder(gz).Encode(&snap); err != nil {
                _ = gz.Close()
                return nil, err
        }
        if err := gz.Close(); err != nil {
                return nil, err
        }
        return buf.Bytes(), nil
}

// DecompressLedger converts a gzip-compressed JSON blob back into a ledger.
func DecompressLedger(data []byte) (*Ledger, error) {
        gz, err := gzip.NewReader(bytes.NewReader(data))
        if err != nil {
                return nil, err
        }
        defer gz.Close()
        var snap ledgerSnapshot
        if err := json.NewDecoder(gz).Decode(&snap); err != nil {
                return nil, err
        }
        l := &Ledger{balances: snap.Balances, blocks: snap.Blocks}
        if l.balances == nil {
                l.balances = make(map[string]uint64)
        }
        return l, nil
}

// SaveCompressedSnapshot writes a compressed snapshot of the ledger to the given path.
func SaveCompressedSnapshot(l *Ledger, path string) error {
	data, err := CompressLedger(l)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadCompressedSnapshot reads a compressed snapshot from disk and returns the decoded ledger.
func LoadCompressedSnapshot(path string) (*Ledger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecompressLedger(data)
}
