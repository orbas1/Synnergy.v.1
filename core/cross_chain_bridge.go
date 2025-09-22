package core

import (
        "crypto/sha256"
        "encoding/base64"
        "encoding/hex"
        "encoding/json"
        "errors"
        "fmt"
        "sort"
        "strconv"
        "strings"
        "sync"
        "time"
)

// AssetBridge defines a token bridge between two chains with a relayer whitelist.
type AssetBridge struct {
	ID       int
	Source   string
	Target   string
	Relayers map[string]struct{}
}

// BridgeTransferRecord records a cross-chain transfer locked on this chain.
type BridgeTransferRecord struct {
        ID       int
        BridgeID int
        From     string
        To       string
        Amount   uint64
        TokenID  string
        Claimed  bool
        ClaimedAt     time.Time
        ProofChecksum string
        RelayTx       string
        Signers       []string
}

// BridgeManager manages bridges and transfer records.
type BridgeManager struct {
        mu             sync.RWMutex
        bridges        map[int]*AssetBridge
        transfers      map[int]*BridgeTransferRecord
        nextBridgeID   int
        nextTransferID int
        ledger         *Ledger
}

var ErrInvalidBridgeProof = errors.New("invalid bridge proof")

const (
        bridgeProofTTL     = 10 * time.Minute
        minProofSignatures = 1
)

type bridgeClaimProof struct {
        TransferID int      `json:"transfer_id"`
        BridgeID   int      `json:"bridge_id"`
        Recipient  string   `json:"recipient"`
        Amount     uint64   `json:"amount"`
        TokenID    string   `json:"token_id"`
        SourceTx   string   `json:"source_tx"`
        Signers    []string `json:"signers"`
        Checksum   string   `json:"checksum"`
        Timestamp  int64    `json:"timestamp"`
}

type bridgeTransferClaimProof struct {
        TransferID string   `json:"transfer_id"`
        BridgeID   string   `json:"bridge_id"`
        Recipient  string   `json:"recipient"`
        Amount     uint64   `json:"amount"`
        TokenID    string   `json:"token_id"`
        SourceTx   string   `json:"source_tx"`
        Signers    []string `json:"signers"`
        Checksum   string   `json:"checksum"`
        Timestamp  int64    `json:"timestamp"`
}

// NewBridgeManager creates a new manager using the provided ledger for balance operations.
func NewBridgeManager(l *Ledger) *BridgeManager {
	return &BridgeManager{
		bridges:   make(map[int]*AssetBridge),
		transfers: make(map[int]*BridgeTransferRecord),
		ledger:    l,
	}
}

// RegisterBridge creates a new bridge definition and returns its ID.
func (m *BridgeManager) RegisterBridge(source, target, relayer string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextBridgeID++
	id := m.nextBridgeID
	relayers := make(map[string]struct{})
	if relayer != "" {
		relayers[relayer] = struct{}{}
	}
	m.bridges[id] = &AssetBridge{ID: id, Source: source, Target: target, Relayers: relayers}
	return id
}

// ListBridges returns all registered bridges.
func (m *BridgeManager) ListBridges() []*AssetBridge {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*AssetBridge, 0, len(m.bridges))
	for _, b := range m.bridges {
		out = append(out, b)
	}
	return out
}

// GetBridge retrieves a bridge by ID.
func (m *BridgeManager) GetBridge(id int) (*AssetBridge, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bridges[id]
	if !ok {
		return nil, errors.New("bridge not found")
	}
	return b, nil
}

// AuthorizeRelayer adds an address to the bridge's relayer whitelist.
func (m *BridgeManager) AuthorizeRelayer(id int, addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bridges[id]
	if !ok {
		return errors.New("bridge not found")
	}
	b.Relayers[addr] = struct{}{}
	return nil
}

// RevokeRelayer removes an address from the bridge's whitelist.
func (m *BridgeManager) RevokeRelayer(id int, addr string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.bridges[id]
	if !ok {
		return errors.New("bridge not found")
	}
	delete(b.Relayers, addr)
	return nil
}

// IsRelayerAuthorized checks if an address is authorized to relay for the bridge.
// It returns false if the bridge does not exist or the relayer is not whitelisted.
func (m *BridgeManager) IsRelayerAuthorized(id int, addr string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.bridges[id]
	if !ok {
		return false
	}
	_, authorized := b.Relayers[addr]
	return authorized
}

// RemoveBridge deletes a bridge definition from the manager.
// It returns an error if the bridge cannot be found.
func (m *BridgeManager) RemoveBridge(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bridges[id]; !ok {
		return errors.New("bridge not found")
	}
	delete(m.bridges, id)
	return nil
}

// Deposit locks assets on the source chain creating a transfer record.
func (m *BridgeManager) Deposit(bridgeID int, from, to string, amount uint64, tokenID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.bridges[bridgeID]; !ok {
		return 0, errors.New("bridge not found")
	}
	if m.ledger == nil {
		return 0, errors.New("ledger not configured")
	}
	tx := &Transaction{From: from, To: "bridge_escrow", Amount: amount}
	if err := m.ledger.ApplyTransaction(tx); err != nil {
		return 0, err
	}
	m.nextTransferID++
	id := m.nextTransferID
	m.transfers[id] = &BridgeTransferRecord{ID: id, BridgeID: bridgeID, From: from, To: to, Amount: amount, TokenID: tokenID}
	return id, nil
}

// Claim releases locked assets to the recipient using a proof placeholder.
// The relayer must be authorized for the bridge that locked the funds.
func (m *BridgeManager) Claim(transferID int, relayer, proof string) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        t, ok := m.transfers[transferID]
        if !ok {
                return errors.New("transfer not found")
        }
        b, ok := m.bridges[t.BridgeID]
        if !ok {
                return errors.New("bridge not found")
        }
        if _, authorized := b.Relayers[relayer]; !authorized {
                return errors.New("relayer not authorized")
        }
        if t.Claimed {
                return errors.New("transfer already claimed")
        }
        payload := bridgeClaimProof{}
        if err := decodeProofPayload(proof, &payload); err != nil {
                return err
        }
        checksum, signers, err := payload.validate(t, relayer, b.Relayers)
        if err != nil {
                return err
        }
        if m.ledger == nil {
                return errors.New("ledger not configured")
        }
        m.ledger.Credit(t.To, t.Amount)
        t.Claimed = true
        t.ClaimedAt = time.Now()
        t.ProofChecksum = checksum
        t.RelayTx = payload.SourceTx
        t.Signers = append([]string(nil), signers...)
        return nil
}

// GetTransfer returns a transfer record by ID.
func (m *BridgeManager) GetTransfer(id int) (*BridgeTransferRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transfers[id]
	if !ok {
		return nil, errors.New("transfer not found")
	}
	return t, nil
}

// ListTransfers returns all transfer records.
func (m *BridgeManager) ListTransfers() []*BridgeTransferRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]*BridgeTransferRecord, 0, len(m.transfers))
	for _, t := range m.transfers {
		out = append(out, t)
	}
	return out
}

// BridgeTransfer represents a transfer record used by BridgeTransferManager with string IDs.
type BridgeTransfer struct {
        ID       string
        BridgeID string
        From     string
        To       string
        Amount   uint64
        TokenID  string
        Status   string
        ProofChecksum string
        ReleasedAt    time.Time
        SourceTx      string
        Signers       []string
}

// BridgeTransferManager manages cross-chain transfer records with string identifiers.
type BridgeTransferManager struct {
	mu        sync.RWMutex
	seq       int
	transfers map[string]*BridgeTransfer
}

// NewBridgeTransferManager creates a new manager.
func NewBridgeTransferManager() *BridgeTransferManager {
	return &BridgeTransferManager{transfers: make(map[string]*BridgeTransfer)}
}

// Deposit locks assets for bridging and records the transfer.
func (m *BridgeTransferManager) Deposit(bridgeID, from, to string, amount uint64, tokenID string) (*BridgeTransfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.seq++
	id := fmt.Sprintf("transfer-%d", m.seq)
	t := &BridgeTransfer{
		ID:       id,
		BridgeID: bridgeID,
		From:     from,
		To:       to,
		Amount:   amount,
		TokenID:  tokenID,
		Status:   "locked",
	}
	m.transfers[id] = t
	return t, nil
}

// Claim releases assets when provided with a valid proof.
func (m *BridgeTransferManager) Claim(id, proof string) error {
        m.mu.Lock()
        defer m.mu.Unlock()
        t, ok := m.transfers[id]
        if !ok {
                return fmt.Errorf("transfer %s not found", id)
        }
        if t.Status != "locked" {
                return fmt.Errorf("transfer %s already claimed", id)
        }
        payload := bridgeTransferClaimProof{}
        if err := decodeProofPayload(proof, &payload); err != nil {
                return err
        }
        checksum, signers, err := payload.validate(t)
        if err != nil {
                return err
        }
        t.Status = "released"
        t.ProofChecksum = checksum
        t.Signers = append([]string(nil), signers...)
        t.SourceTx = payload.SourceTx
        t.ReleasedAt = time.Now()
        return nil
}

// GetTransfer retrieves a transfer by ID.
func (m *BridgeTransferManager) GetTransfer(id string) (*BridgeTransfer, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transfers[id]
	return t, ok
}

// ListTransfers lists all transfer records.
func (m *BridgeTransferManager) ListTransfers() []*BridgeTransfer {
        m.mu.RLock()
        defer m.mu.RUnlock()
        out := make([]*BridgeTransfer, 0, len(m.transfers))
        for _, t := range m.transfers {
                out = append(out, t)
        }
        return out
}

func decodeProofPayload(proof string, target interface{}) error {
        raw := strings.TrimSpace(proof)
        if raw == "" {
                return fmt.Errorf("%w: proof payload empty", ErrInvalidBridgeProof)
        }
        var data []byte
        if decoded, err := base64.StdEncoding.DecodeString(raw); err == nil {
                data = decoded
        } else {
                if !strings.HasPrefix(raw, "{") {
                        return fmt.Errorf("%w: invalid encoding", ErrInvalidBridgeProof)
                }
                data = []byte(raw)
        }
        if err := json.Unmarshal(data, target); err != nil {
                return fmt.Errorf("%w: %v", ErrInvalidBridgeProof, err)
        }
        return nil
}

func computeClaimChecksum(parts ...string) string {
        sanitized := make([]string, len(parts))
        for i, part := range parts {
                sanitized[i] = strings.TrimSpace(part)
        }
        sum := sha256.Sum256([]byte(strings.Join(sanitized, "|")))
        return hex.EncodeToString(sum[:])
}

func normaliseSigners(signers []string) ([]string, error) {
        if len(signers) < minProofSignatures {
                return nil, fmt.Errorf("%w: insufficient signatures", ErrInvalidBridgeProof)
        }
        seen := make(map[string]struct{}, len(signers))
        for _, signer := range signers {
                signer = strings.TrimSpace(signer)
                if signer == "" {
                        return nil, fmt.Errorf("%w: signer cannot be empty", ErrInvalidBridgeProof)
                }
                seen[signer] = struct{}{}
        }
        out := make([]string, 0, len(seen))
        for signer := range seen {
                out = append(out, signer)
        }
        sort.Strings(out)
        return out, nil
}

func (p *bridgeClaimProof) expectedChecksum(record *BridgeTransferRecord) string {
        return computeClaimChecksum(
                strconv.Itoa(p.TransferID),
                strconv.Itoa(p.BridgeID),
                p.SourceTx,
                record.To,
                strconv.FormatUint(record.Amount, 10),
                record.TokenID,
        )
}

func (p *bridgeClaimProof) validate(record *BridgeTransferRecord, relayer string, authorized map[string]struct{}) (string, []string, error) {
        if p.TransferID != record.ID {
                return "", nil, fmt.Errorf("%w: transfer id mismatch", ErrInvalidBridgeProof)
        }
        if p.BridgeID != record.BridgeID {
                return "", nil, fmt.Errorf("%w: bridge id mismatch", ErrInvalidBridgeProof)
        }
        if strings.TrimSpace(p.Recipient) != record.To {
                return "", nil, fmt.Errorf("%w: recipient mismatch", ErrInvalidBridgeProof)
        }
        if p.Amount != record.Amount {
                return "", nil, fmt.Errorf("%w: amount mismatch", ErrInvalidBridgeProof)
        }
        if strings.TrimSpace(p.TokenID) != record.TokenID {
                return "", nil, fmt.Errorf("%w: token id mismatch", ErrInvalidBridgeProof)
        }
        if strings.TrimSpace(p.SourceTx) == "" {
                return "", nil, fmt.Errorf("%w: source_tx required", ErrInvalidBridgeProof)
        }
        signers, err := normaliseSigners(p.Signers)
        if err != nil {
                return "", nil, err
        }
        hasRelayer := false
        for _, signer := range signers {
                if signer == relayer {
                        hasRelayer = true
                }
                if _, ok := authorized[signer]; !ok {
                        return "", nil, fmt.Errorf("%w: signer %s not authorized", ErrInvalidBridgeProof, signer)
                }
        }
        if !hasRelayer {
                return "", nil, fmt.Errorf("%w: relayer signature required", ErrInvalidBridgeProof)
        }
        if p.Timestamp == 0 {
                return "", nil, fmt.Errorf("%w: timestamp required", ErrInvalidBridgeProof)
        }
        proofTime := time.Unix(p.Timestamp, 0)
        now := time.Now()
        if proofTime.After(now.Add(bridgeProofTTL)) || now.Sub(proofTime) > bridgeProofTTL {
                return "", nil, fmt.Errorf("%w: proof expired", ErrInvalidBridgeProof)
        }
        checksum := p.expectedChecksum(record)
        if !strings.EqualFold(strings.TrimSpace(p.Checksum), checksum) {
                return "", nil, fmt.Errorf("%w: checksum mismatch", ErrInvalidBridgeProof)
        }
        p.Signers = signers
        return checksum, signers, nil
}

func (p *bridgeTransferClaimProof) expectedChecksum(record *BridgeTransfer) string {
        return computeClaimChecksum(
                p.TransferID,
                p.BridgeID,
                p.SourceTx,
                record.To,
                strconv.FormatUint(record.Amount, 10),
                record.TokenID,
        )
}

func (p *bridgeTransferClaimProof) validate(record *BridgeTransfer) (string, []string, error) {
        if p.TransferID != record.ID {
                return "", nil, fmt.Errorf("%w: transfer id mismatch", ErrInvalidBridgeProof)
        }
        if p.BridgeID != record.BridgeID {
                return "", nil, fmt.Errorf("%w: bridge id mismatch", ErrInvalidBridgeProof)
        }
        if strings.TrimSpace(p.Recipient) != record.To {
                return "", nil, fmt.Errorf("%w: recipient mismatch", ErrInvalidBridgeProof)
        }
        if p.Amount != record.Amount {
                return "", nil, fmt.Errorf("%w: amount mismatch", ErrInvalidBridgeProof)
        }
        if strings.TrimSpace(p.TokenID) != record.TokenID {
                return "", nil, fmt.Errorf("%w: token id mismatch", ErrInvalidBridgeProof)
        }
        if strings.TrimSpace(p.SourceTx) == "" {
                return "", nil, fmt.Errorf("%w: source_tx required", ErrInvalidBridgeProof)
        }
        signers, err := normaliseSigners(p.Signers)
        if err != nil {
                return "", nil, err
        }
        if p.Timestamp == 0 {
                return "", nil, fmt.Errorf("%w: timestamp required", ErrInvalidBridgeProof)
        }
        proofTime := time.Unix(p.Timestamp, 0)
        now := time.Now()
        if proofTime.After(now.Add(bridgeProofTTL)) || now.Sub(proofTime) > bridgeProofTTL {
                return "", nil, fmt.Errorf("%w: proof expired", ErrInvalidBridgeProof)
        }
        checksum := p.expectedChecksum(record)
        if !strings.EqualFold(strings.TrimSpace(p.Checksum), checksum) {
                return "", nil, fmt.Errorf("%w: checksum mismatch", ErrInvalidBridgeProof)
        }
        p.Signers = signers
        return checksum, signers, nil
}
