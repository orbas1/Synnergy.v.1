package synnergy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

// AIModelMetadata describes a published model.
type AIModelMetadata struct {
	Hash        string
	CID         string
	RoyaltyBps  uint16
	PublishedAt time.Time
}

// Escrow holds pending payments for a model transaction.
type Escrow struct {
	ID        string
	ListingID string
	Buyer     string
	Renter    string
	Amount    uint64
	ExpiresAt time.Time
	Released  bool
}

// AIService provides prediction utilities and manages models.
type AIService struct {
	mu          sync.RWMutex
	models      map[string]AIModelMetadata
	marketplace *ModelMarketplace
	escrows     map[string]Escrow
	nextEscrow  uint64
}

// NewAIService constructs a new AIService instance.
func NewAIService() *AIService {
	return &AIService{
		models:      make(map[string]AIModelMetadata),
		marketplace: NewModelMarketplace(),
		escrows:     make(map[string]Escrow),
	}
}

// PredictFraud returns a pseudo probabilistic fraud score for a transaction JSON payload.
func (s *AIService) PredictFraud(txJSON []byte) (float64, error) {
	var m map[string]any
	if err := json.Unmarshal(txJSON, &m); err != nil {
		return 0, err
	}
	h := sha256.Sum256(txJSON)
	return float64(h[0]) / 255.0, nil
}

// OptimiseBaseFee suggests a base fee for the next block given network statistics.
func (s *AIService) OptimiseBaseFee(statsJSON []byte) (uint64, error) {
	var m map[string]float64
	if err := json.Unmarshal(statsJSON, &m); err != nil {
		return 0, err
	}
	fee := uint64(m["avgGasPrice"] * 1.1)
	return fee, nil
}

// ForecastVolume estimates upcoming transaction volume.
func (s *AIService) ForecastVolume(statsJSON []byte) (uint64, error) {
	var m map[string]float64
	if err := json.Unmarshal(statsJSON, &m); err != nil {
		return 0, err
	}
	vol := uint64(m["recentTxs"] * 1.05)
	return vol, nil
}

// PublishModel registers model metadata and returns its hash.
func (s *AIService) PublishModel(cid string, royaltyBps uint16) (string, error) {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s-%d-%d", cid, royaltyBps, time.Now().UnixNano())))
	hash := hex.EncodeToString(h[:])
	s.mu.Lock()
	s.models[hash] = AIModelMetadata{
		Hash:        hash,
		CID:         cid,
		RoyaltyBps:  royaltyBps,
		PublishedAt: time.Now().UTC(),
	}
	s.mu.Unlock()
	return hash, nil
}

// FetchModel returns metadata for a published model.
func (s *AIService) FetchModel(hash string) (AIModelMetadata, bool) {
	s.mu.RLock()
	m, ok := s.models[hash]
	s.mu.RUnlock()
	return m, ok
}

// ListModel creates a marketplace listing for a model and returns the listing ID.
func (s *AIService) ListModel(hash, cid, seller string, price uint64) string {
	return s.marketplace.AddListing(hash, cid, seller, price)
}

// BuyModel places funds in escrow for a listing purchase.
func (s *AIService) BuyModel(listingID, buyer string, amount uint64) (string, error) {
	listing, ok := s.marketplace.Get(listingID)
	if !ok || !listing.Active {
		return "", errors.New("listing not found")
	}
	if amount < listing.Price {
		return "", errors.New("insufficient amount")
	}
	s.mu.Lock()
	s.nextEscrow++
	id := fmt.Sprintf("escrow-%d", s.nextEscrow)
	s.escrows[id] = Escrow{
		ID:        id,
		ListingID: listingID,
		Buyer:     buyer,
		Amount:    amount,
	}
	s.mu.Unlock()
	return id, nil
}

// RentModel creates a time-limited escrow for renting a model.
func (s *AIService) RentModel(listingID, renter string, hours int, amount uint64) (string, error) {
	listing, ok := s.marketplace.Get(listingID)
	if !ok || !listing.Active {
		return "", errors.New("listing not found")
	}
	if amount < listing.Price {
		return "", errors.New("insufficient amount")
	}
	s.mu.Lock()
	s.nextEscrow++
	id := fmt.Sprintf("escrow-%d", s.nextEscrow)
	s.escrows[id] = Escrow{
		ID:        id,
		ListingID: listingID,
		Renter:    renter,
		Amount:    amount,
		ExpiresAt: time.Now().Add(time.Duration(hours) * time.Hour).UTC(),
	}
	s.mu.Unlock()
	return id, nil
}

// ReleaseEscrow releases funds to the seller for the given escrow ID.
func (s *AIService) ReleaseEscrow(id string) error {
	s.mu.Lock()
	escrow, ok := s.escrows[id]
	if !ok || escrow.Released {
		s.mu.Unlock()
		return errors.New("escrow not found or already released")
	}
	escrow.Released = true
	s.escrows[id] = escrow
	s.mu.Unlock()
	return nil
}
