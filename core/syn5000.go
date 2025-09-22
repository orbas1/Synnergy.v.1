package core

import (
	"crypto/ed25519"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"
)

const syn5000BasisPoints = 10000

var (
	// ErrBetNotFound indicates that the requested bet does not exist in the
	// token's in-memory ledger.
	ErrBetNotFound = errors.New("syn5000: bet not found")
	// ErrBetFinalized signals that a bet has already been settled or
	// cancelled.
	ErrBetFinalized = errors.New("syn5000: bet already finalized")
	// ErrExposureLimit is returned when the configured exposure caps would
	// be exceeded by a placement.
	ErrExposureLimit = errors.New("syn5000: exposure limit exceeded")
	// ErrInvalidBetAmount captures attempts to place wagers below the
	// minimum or above the maximum configured limits.
	ErrInvalidBetAmount = errors.New("syn5000: invalid bet amount")
	// ErrInvalidOdds indicates that the provided odds fall outside the
	// permitted range.
	ErrInvalidOdds = errors.New("syn5000: invalid odds")
	// ErrMissingBettor marks placements that do not declare a bettor ID.
	ErrMissingBettor = errors.New("syn5000: bettor required")
	// ErrMissingGame marks placements that do not target a game bucket.
	ErrMissingGame = errors.New("syn5000: game required")
	// ErrDuplicateNonce guards against replay of signed bet placements.
	ErrDuplicateNonce = errors.New("syn5000: duplicate placement nonce")
	// ErrSignatureRequired indicates the placement must be signed.
	ErrSignatureRequired = errors.New("syn5000: signature required")
	// ErrUnknownParticipant is returned when a bet references an
	// unregistered bettor for signature validation.
	ErrUnknownParticipant = errors.New("syn5000: unknown participant")
	// ErrInvalidSignature indicates that the supplied signature does not
	// validate against the stored participant key.
	ErrInvalidSignature = errors.New("syn5000: invalid signature")
	// ErrMissingNonce marks placements that require an explicit nonce.
	ErrMissingNonce = errors.New("syn5000: nonce required")
)

// BetState captures the lifecycle state of a gambling bet.
type BetState string

const (
	// BetStatePending indicates an active bet awaiting resolution.
	BetStatePending BetState = "pending"
	// BetStateWon indicates a winning bet where the bettor receives the
	// calculated payout.
	BetStateWon BetState = "won"
	// BetStateLost indicates a losing bet.
	BetStateLost BetState = "lost"
	// BetStateCancelled indicates a bet voided by compliance or operator
	// intervention.
	BetStateCancelled BetState = "cancelled"
)

// BetPlacement describes a bet request submitted by a participant.
type BetPlacement struct {
	Bettor          string
	Amount          uint64
	OddsBasisPoints uint32
	Game            string
	Nonce           string
	Metadata        string
}

// signingDigest returns a deterministic digest of the placement that is used
// for signature validation.
func (p BetPlacement) signingDigest() [32]byte {
	payload := fmt.Sprintf("%s|%d|%d|%s|%s|%s", p.Bettor, p.Amount, p.OddsBasisPoints, p.Game, p.Nonce, p.Metadata)
	return sha256.Sum256([]byte(payload))
}

// betNonceKey creates the nonce key used to detect replayed signed requests.
func (p BetPlacement) betNonceKey() string {
	return fmt.Sprintf("%s|%s", p.Bettor, p.Nonce)
}

// BetRecord stores betting activity for SYN5000 tokens.
type BetRecord struct {
	ID              uint64
	Placement       BetPlacement
	PotentialPayout uint64
	Liability       uint64
	State           BetState
	CreatedAt       time.Time
	ResolvedAt      time.Time
	ResolutionNote  string
	Signature       []byte
}

func (b *BetRecord) clone() BetRecord {
	cp := *b
	if len(b.Signature) > 0 {
		cp.Signature = append([]byte(nil), b.Signature...)
	}
	return cp
}

// BetFilter filters bet listings for CLI and API consumers.
type BetFilter struct {
	Bettor string
	Game   string
	State  BetState
	Limit  int
}

func (f BetFilter) matches(rec *BetRecord) bool {
	if f.Bettor != "" && rec.Placement.Bettor != f.Bettor {
		return false
	}
	if f.Game != "" && rec.Placement.Game != f.Game {
		return false
	}
	if f.State != "" && rec.State != f.State {
		return false
	}
	return true
}

// SYN5000Snapshot summarises outstanding liabilities for operations dashboards
// and CLI status queries.
type SYN5000Snapshot struct {
	Name           string
	Symbol         string
	TotalBets      int
	Pending        int
	Won            int
	Lost           int
	Cancelled      int
	GlobalExposure uint64
	ExposureByGame map[string]uint64
	UpdatedAt      time.Time
}

// SYN5000Config controls risk, compliance and signature validation behaviour.
type SYN5000Config struct {
	MinBet             uint64
	MaxBet             uint64
	MinOddsBasisPoints uint32
	MaxOddsBasisPoints uint32
	MaxExposurePerGame uint64
	MaxGlobalExposure  uint64
	RequireSignature   bool
}

func normalizeSYN5000Config(cfg *SYN5000Config) {
	if cfg.MinBet == 0 {
		cfg.MinBet = 1
	}
	if cfg.MaxBet != 0 && cfg.MaxBet < cfg.MinBet {
		cfg.MaxBet = cfg.MinBet
	}
	if cfg.MinOddsBasisPoints == 0 {
		cfg.MinOddsBasisPoints = syn5000BasisPoints
	}
	if cfg.MaxOddsBasisPoints != 0 && cfg.MaxOddsBasisPoints < cfg.MinOddsBasisPoints {
		cfg.MaxOddsBasisPoints = cfg.MinOddsBasisPoints
	}
}

// DefaultSYN5000Config returns the default enterprise configuration used by new
// tokens unless overridden.
func DefaultSYN5000Config() SYN5000Config {
	cfg := SYN5000Config{
		MinBet:             1,
		MaxBet:             0,
		MinOddsBasisPoints: syn5000BasisPoints,
		MaxOddsBasisPoints: syn5000BasisPoints * 200, // up to 200x multipliers
		MaxExposurePerGame: 0,
		MaxGlobalExposure:  0,
		RequireSignature:   false,
	}
	normalizeSYN5000Config(&cfg)
	return cfg
}

// SYN5000Option customises a new token instance.
type SYN5000Option func(*SYN5000Token)

// WithSYN5000Config overrides the default configuration for a token.
func WithSYN5000Config(cfg SYN5000Config) SYN5000Option {
	return func(t *SYN5000Token) {
		normalizeSYN5000Config(&cfg)
		t.config = cfg
	}
}

// WithSYN5000Participant pre-registers a participant public key for signature
// validation.
func WithSYN5000Participant(bettor string, pubKey ed25519.PublicKey) SYN5000Option {
	return func(t *SYN5000Token) {
		if t.bettors == nil {
			t.bettors = make(map[string]ed25519.PublicKey)
		}
		cp := append(ed25519.PublicKey(nil), pubKey...)
		t.bettors[bettor] = cp
	}
}

// SYN5000Token implements the GamblingToken interface with comprehensive risk
// and compliance controls.
type SYN5000Token struct {
	Name     string
	Symbol   string
	Decimals uint8

	mu             sync.RWMutex
	nextBetID      uint64
	bets           map[uint64]*BetRecord
	exposureByGame map[string]uint64
	globalExposure uint64
	config         SYN5000Config
	bettors        map[string]ed25519.PublicKey
	usedNonces     map[string]struct{}
	nonceCounter   uint64
}

// NewSYN5000Token creates a new gambling token instance with optional
// configuration.
func NewSYN5000Token(name, symbol string, decimals uint8, opts ...SYN5000Option) *SYN5000Token {
	t := &SYN5000Token{
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		bets:           make(map[uint64]*BetRecord),
		exposureByGame: make(map[string]uint64),
		usedNonces:     make(map[string]struct{}),
		bettors:        make(map[string]ed25519.PublicKey),
		config:         DefaultSYN5000Config(),
	}
	for _, opt := range opts {
		opt(t)
	}
	normalizeSYN5000Config(&t.config)
	return t
}

// RegisterParticipant records a bettor's ed25519 public key for signature
// verification.
func (t *SYN5000Token) RegisterParticipant(bettor string, pubKey ed25519.PublicKey) error {
	if bettor == "" {
		return ErrMissingBettor
	}
	if len(pubKey) != ed25519.PublicKeySize {
		return fmt.Errorf("syn5000: invalid ed25519 public key length")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	cp := append(ed25519.PublicKey(nil), pubKey...)
	t.bettors[bettor] = cp
	return nil
}

// SetConfig updates the runtime configuration, applying default values where
// omitted.
func (t *SYN5000Token) SetConfig(cfg SYN5000Config) {
	normalizeSYN5000Config(&cfg)
	t.mu.Lock()
	t.config = cfg
	t.mu.Unlock()
}

// PlaceBet records a new bet using basic parameters. Metadata, explicit nonces
// and signatures are available via PlaceBetDetailed and PlaceBetWithSignature.
func (t *SYN5000Token) PlaceBet(bettor string, amount uint64, odds float64, game string) (uint64, error) {
	oddsBP, err := ConvertOddsToBasisPoints(odds)
	if err != nil {
		return 0, err
	}
	placement := BetPlacement{
		Bettor:          bettor,
		Amount:          amount,
		OddsBasisPoints: oddsBP,
		Game:            game,
	}
	return t.placeBet(placement, nil, false)
}

// PlaceBetDetailed allows callers to specify explicit placement metadata while
// relying on the token to enforce risk checks.
func (t *SYN5000Token) PlaceBetDetailed(placement BetPlacement) (uint64, error) {
	return t.placeBet(placement, nil, false)
}

// PlaceBetWithSignature records a signed bet placement, enforcing signature
// validation regardless of global configuration.
func (t *SYN5000Token) PlaceBetWithSignature(placement BetPlacement, signature []byte) (uint64, error) {
	if len(signature) == 0 {
		return 0, ErrSignatureRequired
	}
	if placement.Nonce == "" {
		return 0, ErrMissingNonce
	}
	return t.placeBet(placement, signature, true)
}

func (t *SYN5000Token) placeBet(placement BetPlacement, signature []byte, enforceSignature bool) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if placement.Bettor == "" {
		return 0, ErrMissingBettor
	}
	if placement.Game == "" {
		return 0, ErrMissingGame
	}
	if placement.Amount == 0 {
		return 0, ErrInvalidBetAmount
	}
	if placement.Nonce == "" {
		t.nonceCounter++
		placement.Nonce = fmt.Sprintf("syn5000-%d", t.nonceCounter)
	}
	if placement.OddsBasisPoints == 0 {
		placement.OddsBasisPoints = t.config.MinOddsBasisPoints
	}
	if placement.Amount < t.config.MinBet {
		return 0, fmt.Errorf("%w: %d below minimum %d", ErrInvalidBetAmount, placement.Amount, t.config.MinBet)
	}
	if t.config.MaxBet != 0 && placement.Amount > t.config.MaxBet {
		return 0, fmt.Errorf("%w: %d exceeds maximum %d", ErrInvalidBetAmount, placement.Amount, t.config.MaxBet)
	}
	if placement.OddsBasisPoints < t.config.MinOddsBasisPoints {
		return 0, fmt.Errorf("%w: %d below minimum %d", ErrInvalidOdds, placement.OddsBasisPoints, t.config.MinOddsBasisPoints)
	}
	if t.config.MaxOddsBasisPoints != 0 && placement.OddsBasisPoints > t.config.MaxOddsBasisPoints {
		return 0, fmt.Errorf("%w: %d exceeds maximum %d", ErrInvalidOdds, placement.OddsBasisPoints, t.config.MaxOddsBasisPoints)
	}

	key := placement.betNonceKey()
	if _, exists := t.usedNonces[key]; exists {
		return 0, ErrDuplicateNonce
	}

	potential, liability, err := computePayoutAndLiability(placement.Amount, placement.OddsBasisPoints)
	if err != nil {
		return 0, err
	}
	if err := t.checkExposureLocked(placement.Game, liability); err != nil {
		return 0, err
	}

	shouldVerify := enforceSignature || t.config.RequireSignature
	if len(signature) > 0 {
		shouldVerify = true
	}

	var sigCopy []byte
	if shouldVerify {
		if len(signature) == 0 {
			return 0, ErrSignatureRequired
		}
		pub, ok := t.bettors[placement.Bettor]
		if !ok {
			return 0, ErrUnknownParticipant
		}
		digest := placement.signingDigest()
		if !ed25519.Verify(pub, digest[:], signature) {
			return 0, ErrInvalidSignature
		}
		sigCopy = append([]byte(nil), signature...)
	}

	t.nextBetID++
	betID := t.nextBetID
	record := &BetRecord{
		ID:              betID,
		Placement:       placement,
		PotentialPayout: potential,
		Liability:       liability,
		State:           BetStatePending,
		CreatedAt:       time.Now().UTC(),
		Signature:       sigCopy,
	}
	t.bets[betID] = record
	t.usedNonces[key] = struct{}{}
	t.applyExposureLocked(placement.Game, liability)
	return betID, nil
}

// ResolveBet settles a bet, updating exposure and returning the payout for
// winners.
func (t *SYN5000Token) ResolveBet(betID uint64, win bool) (uint64, error) {
	return t.ResolveBetWithNote(betID, win, "")
}

// ResolveBetWithNote resolves a bet, capturing an operator note for audit
// logging.
func (t *SYN5000Token) ResolveBetWithNote(betID uint64, win bool, note string) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	bet, ok := t.bets[betID]
	if !ok {
		return 0, ErrBetNotFound
	}
	if bet.State != BetStatePending {
		return 0, ErrBetFinalized
	}

	t.releaseExposureLocked(bet.Placement.Game, bet.Liability)
	bet.ResolvedAt = time.Now().UTC()
	bet.ResolutionNote = note

	if win {
		bet.State = BetStateWon
		return bet.PotentialPayout, nil
	}

	bet.State = BetStateLost
	return 0, nil
}

// CancelBet voids a bet and releases its reserved exposure.
func (t *SYN5000Token) CancelBet(betID uint64, note string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	bet, ok := t.bets[betID]
	if !ok {
		return ErrBetNotFound
	}
	if bet.State != BetStatePending {
		return ErrBetFinalized
	}

	t.releaseExposureLocked(bet.Placement.Game, bet.Liability)
	bet.State = BetStateCancelled
	bet.ResolvedAt = time.Now().UTC()
	bet.ResolutionNote = note
	return nil
}

// GetBet returns a bet record by ID.
func (t *SYN5000Token) GetBet(betID uint64) (*BetRecord, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	b, ok := t.bets[betID]
	if !ok {
		return nil, false
	}
	cp := b.clone()
	return &cp, true
}

// ListBets returns bets matching the supplied filter, ordered by bet ID.
func (t *SYN5000Token) ListBets(filter BetFilter) []BetRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]BetRecord, 0, len(t.bets))
	for _, bet := range t.bets {
		if !filter.matches(bet) {
			continue
		}
		out = append(out, bet.clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	if filter.Limit > 0 && len(out) > filter.Limit {
		out = out[:filter.Limit]
	}
	return out
}

// Snapshot summarises pending exposure and bet lifecycle counts for the token.
func (t *SYN5000Token) Snapshot() SYN5000Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()

	snap := SYN5000Snapshot{
		Name:           t.Name,
		Symbol:         t.Symbol,
		ExposureByGame: make(map[string]uint64, len(t.exposureByGame)),
		GlobalExposure: t.globalExposure,
		TotalBets:      len(t.bets),
		UpdatedAt:      time.Now().UTC(),
	}
	for game, value := range t.exposureByGame {
		snap.ExposureByGame[game] = value
	}
	for _, bet := range t.bets {
		switch bet.State {
		case BetStatePending:
			snap.Pending++
		case BetStateWon:
			snap.Won++
		case BetStateLost:
			snap.Lost++
		case BetStateCancelled:
			snap.Cancelled++
		}
	}
	return snap
}

func (t *SYN5000Token) checkExposureLocked(game string, liability uint64) error {
	if liability == 0 {
		return nil
	}
	if t.config.MaxExposurePerGame != 0 {
		current := t.exposureByGame[game]
		if willOverflow(current, liability) || current+liability > t.config.MaxExposurePerGame {
			return ErrExposureLimit
		}
	}
	if t.config.MaxGlobalExposure != 0 {
		if willOverflow(t.globalExposure, liability) || t.globalExposure+liability > t.config.MaxGlobalExposure {
			return ErrExposureLimit
		}
	}
	return nil
}

func (t *SYN5000Token) applyExposureLocked(game string, liability uint64) {
	if liability == 0 {
		return
	}
	t.exposureByGame[game] += liability
	t.globalExposure += liability
}

func (t *SYN5000Token) releaseExposureLocked(game string, liability uint64) {
	if liability == 0 {
		return
	}
	if current := t.exposureByGame[game]; current >= liability {
		t.exposureByGame[game] = current - liability
	} else {
		t.exposureByGame[game] = 0
	}
	if t.globalExposure >= liability {
		t.globalExposure -= liability
	} else {
		t.globalExposure = 0
	}
}

func computePayoutAndLiability(amount uint64, oddsBasisPoints uint32) (uint64, uint64, error) {
	if oddsBasisPoints < syn5000BasisPoints {
		return 0, 0, ErrInvalidOdds
	}
	payout, err := mulDiv(amount, uint64(oddsBasisPoints), syn5000BasisPoints)
	if err != nil {
		return 0, 0, err
	}
	if payout < amount {
		return 0, 0, ErrInvalidOdds
	}
	liability := payout - amount
	return payout, liability, nil
}

func mulDiv(amount uint64, multiplier uint64, divisor uint64) (uint64, error) {
	if divisor == 0 {
		return 0, errors.New("syn5000: division by zero")
	}
	var a big.Int
	a.SetUint64(amount)
	var m big.Int
	m.SetUint64(multiplier)
	a.Mul(&a, &m)
	var d big.Int
	d.SetUint64(divisor)
	a.Div(&a, &d)
	if !a.IsUint64() {
		return 0, fmt.Errorf("syn5000: value exceeds uint64 range")
	}
	return a.Uint64(), nil
}

func willOverflow(a, b uint64) bool {
	return math.MaxUint64-a < b
}

// ConvertOddsToBasisPoints converts floating odds to basis points used internally.
func ConvertOddsToBasisPoints(odds float64) (uint32, error) {
	return oddsToBasisPoints(odds)
}

func oddsToBasisPoints(odds float64) (uint32, error) {
	if math.IsNaN(odds) || math.IsInf(odds, 0) {
		return 0, ErrInvalidOdds
	}
	if odds <= 0 {
		return 0, ErrInvalidOdds
	}
	bp := math.Round(odds * syn5000BasisPoints)
	if bp <= 0 || bp > math.MaxUint32 {
		return 0, ErrInvalidOdds
	}
	return uint32(bp), nil
}

var _ GamblingToken = (*SYN5000Token)(nil)
