package synnergy

import (
	"crypto/sha256"
	"errors"
	"sync"
)

// FraudResult captures fraud probability for a transaction.
type FraudResult struct {
	TxID  string
	Score float64
}

// InferenceEngine executes model inference and batch analysis.
type InferenceEngine struct {
	mu     sync.RWMutex
	models map[string][]byte
}

// NewInferenceEngine constructs a new InferenceEngine.
func NewInferenceEngine() *InferenceEngine {
	return &InferenceEngine{models: make(map[string][]byte)}
}

// LoadModel loads model data into the engine.
func (e *InferenceEngine) LoadModel(hash string, data []byte) {
	e.mu.Lock()
	e.models[hash] = data
	e.mu.Unlock()
}

// Run executes a model over input and returns a deterministic output hash.
func (e *InferenceEngine) Run(hash string, input []byte) ([]byte, error) {
	e.mu.RLock()
	_, ok := e.models[hash]
	e.mu.RUnlock()
	if !ok {
		return nil, errors.New("model not found")
	}
	h := sha256.Sum256(append([]byte(hash), input...))
	return h[:], nil
}

// Analyse evaluates transactions for fraud risk using a simple heuristic.
func (e *InferenceEngine) Analyse(txIDs []string) []FraudResult {
	results := make([]FraudResult, len(txIDs))
	for i, id := range txIDs {
		h := sha256.Sum256([]byte(id))
		score := float64(h[0]) / 255.0
		results[i] = FraudResult{TxID: id, Score: score}
	}
	return results
}
