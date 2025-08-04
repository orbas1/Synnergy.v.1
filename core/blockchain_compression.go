package core

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
)

// CompressLedger returns the gzip-compressed JSON encoding of the provided ledger.
func CompressLedger(l *Ledger) ([]byte, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if err := json.NewEncoder(gz).Encode(l); err != nil {
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
	var l Ledger
	if err := json.NewDecoder(gz).Decode(&l); err != nil {
		return nil, err
	}
	if l.balances == nil {
		l.balances = make(map[string]uint64)
	}
	return &l, nil
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
