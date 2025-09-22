package core

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	synn "synnergy"
)

var (
	// ErrGrantSignatureInvalid is returned when a wallet proof fails verification.
	ErrGrantSignatureInvalid = errors.New("invalid grant signature")
	// ErrGrantConsensusOffline indicates the consensus engine is unavailable for coordination.
	ErrGrantConsensusOffline = errors.New("consensus unavailable")
	// ErrGrantAuthorityUnknown reports missing authority enrolment for a wallet address.
	ErrGrantAuthorityUnknown = errors.New("authority node unknown")
)

// WalletProof captures an authenticated wallet approval for a grant action.
// Proofs embed the original message to prevent replay across different operations.
//
// Message must be signed via Wallet.SignMessage so VerifyMessageSignature can validate it.
type WalletProof struct {
	Address   string
	PublicKey *ecdsa.PublicKey
	Signature []byte
	Message   []byte
}

// NewWalletProof builds a WalletProof from a wallet and message payload.
func NewWalletProof(wallet *Wallet, message []byte) (*WalletProof, error) {
	if wallet == nil {
		return nil, errors.New("wallet required")
	}
	if len(message) == 0 {
		return nil, errors.New("message required")
	}
	sig, err := wallet.SignMessage(message)
	if err != nil {
		return nil, err
	}
	return &WalletProof{
		Address:   wallet.Address,
		PublicKey: &wallet.PrivateKey.PublicKey,
		Signature: append([]byte(nil), sig...),
		Message:   append([]byte(nil), message...),
	}, nil
}

// Verify checks the signature for the embedded message.
func (p *WalletProof) Verify() bool {
	if p == nil {
		return false
	}
	return VerifyMessageSignature(p.Message, p.Signature, p.PublicKey)
}

func (p *WalletProof) clone() *WalletProof {
	if p == nil {
		return nil
	}
	cp := &WalletProof{Address: p.Address}
	if p.PublicKey != nil {
		cp.PublicKey = &ecdsa.PublicKey{Curve: p.PublicKey.Curve, X: new(big.Int).Set(p.PublicKey.X), Y: new(big.Int).Set(p.PublicKey.Y)}
	}
	if len(p.Signature) > 0 {
		cp.Signature = append([]byte(nil), p.Signature...)
	}
	if len(p.Message) > 0 {
		cp.Message = append([]byte(nil), p.Message...)
	}
	return cp
}

// GrantCreationRequest encapsulates inputs for provisioning a new grant.
type GrantCreationRequest struct {
	Beneficiary string
	Name        string
	Amount      uint64
	Authorizers []string
	Creator     *WalletProof
}

// GrantAuthorizationRequest captures data required to authorise a wallet for a grant.
type GrantAuthorizationRequest struct {
	GrantID uint64
	Proof   *WalletProof
}

// GrantReleaseRequest represents a disbursement request for a grant.
type GrantReleaseRequest struct {
	GrantID uint64
	Amount  uint64
	Note    string
	Proof   *WalletProof
}

// GrantCreateMessage returns the canonical message used when signing grant creation requests.
func GrantCreateMessage(beneficiary, name string, amount uint64) []byte {
	return []byte(fmt.Sprintf("grant:create:%s:%s:%d", beneficiary, name, amount))
}

// GrantAuthorizeMessage returns the canonical message for grant authorisation.
func GrantAuthorizeMessage(id uint64) []byte {
	return []byte(fmt.Sprintf("grant:authorize:%d", id))
}

// GrantReleaseMessage returns the canonical message for grant disbursements.
func GrantReleaseMessage(id uint64, amount uint64, note string) []byte {
	return []byte(fmt.Sprintf("grant:release:%d:%d:%s", id, amount, note))
}

// GrantOrchestrator coordinates the registry, consensus engine, VM opcodes and authority directory.
type GrantOrchestrator struct {
	registry    *GrantRegistry
	vm          *synn.SimpleVM
	consensus   *SynnergyConsensus
	authorities *AuthorityNodeRegistry

	mu       sync.RWMutex
	watchers map[uint64][]chan<- GrantEvent

	vmOnce sync.Once
}

// NewGrantOrchestrator wires registry operations into the VM and supporting services.
func NewGrantOrchestrator(registry *GrantRegistry, vm *synn.SimpleVM, consensus *SynnergyConsensus, authorities *AuthorityNodeRegistry) (*GrantOrchestrator, error) {
	if registry == nil {
		return nil, errors.New("registry required")
	}
	if vm == nil {
		return nil, errors.New("vm required")
	}
	if !vm.Status() {
		if err := vm.Start(); err != nil {
			return nil, fmt.Errorf("start vm: %w", err)
		}
	}
	orchestrator := &GrantOrchestrator{
		registry:    registry,
		vm:          vm,
		consensus:   consensus,
		authorities: authorities,
		watchers:    make(map[uint64][]chan<- GrantEvent),
	}
	orchestrator.registerVMHandlers()
	return orchestrator, nil
}

// EnsureAuthority enrols an address into the authority directory when absent.
func (o *GrantOrchestrator) EnsureAuthority(addr, role string) {
	if o == nil || o.authorities == nil || addr == "" {
		return
	}
	if o.authorities.IsAuthorityNode(addr) {
		return
	}
	_, _ = o.authorities.Register(addr, role)
}

// CreateGrant provisions a grant after validating consensus status and wallet proof.
func (o *GrantOrchestrator) CreateGrant(ctx context.Context, req GrantCreationRequest) (*GrantRecord, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := o.ensureConsensus(); err != nil {
		return nil, err
	}
	expected := GrantCreateMessage(req.Beneficiary, req.Name, req.Amount)
	if err := o.verifyProof(req.Creator, expected); err != nil {
		return nil, err
	}
	o.EnsureAuthority(req.Creator.Address, "grant-admin")
	for _, addr := range req.Authorizers {
		if addr == "" {
			continue
		}
		o.EnsureAuthority(addr, "grant-authoriser")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	id, err := o.registry.CreateGrant(req.Beneficiary, req.Name, req.Amount, req.Authorizers...)
	if err != nil {
		return nil, err
	}
	record, ok := o.registry.GetGrant(id)
	if !ok {
		return nil, ErrGrantNotFound
	}
	if len(record.Events) > 0 {
		o.emit(record.ID, record.Events[len(record.Events)-1])
	}
	return record, nil
}

// AuthorizeGrant records an additional authorised wallet for the grant.
func (o *GrantOrchestrator) AuthorizeGrant(ctx context.Context, req GrantAuthorizationRequest) (*GrantEvent, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	expected := GrantAuthorizeMessage(req.GrantID)
	if err := o.verifyProof(req.Proof, expected); err != nil {
		return nil, err
	}
	if o.authorities != nil {
		o.EnsureAuthority(req.Proof.Address, "grant-authoriser")
	}
	evt, err := o.registry.Authorize(req.GrantID, req.Proof.Address)
	if err != nil {
		return nil, err
	}
	o.emit(req.GrantID, *evt)
	return evt, nil
}

// ReleaseGrant disburses funds from the grant after verifying the signer is authorised.
func (o *GrantOrchestrator) ReleaseGrant(ctx context.Context, req GrantReleaseRequest) (*GrantEvent, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	expected := GrantReleaseMessage(req.GrantID, req.Amount, req.Note)
	if err := o.verifyProof(req.Proof, expected); err != nil {
		return nil, err
	}
	if o.authorities != nil && !o.authorities.IsAuthorityNode(req.Proof.Address) {
		return nil, ErrGrantAuthorityUnknown
	}
	evt, err := o.registry.Disburse(req.GrantID, req.Amount, req.Note, req.Proof.Address)
	if err != nil {
		return nil, err
	}
	o.emit(req.GrantID, *evt)
	return evt, nil
}

// GetGrant returns a copy of the grant record by ID.
func (o *GrantOrchestrator) GetGrant(id uint64) (*GrantRecord, bool) {
	return o.registry.GetGrant(id)
}

// ListGrants enumerates all known grants.
func (o *GrantOrchestrator) ListGrants() []*GrantRecord {
	return o.registry.ListGrants()
}

// Audit returns the grant event trail.
func (o *GrantOrchestrator) Audit(id uint64) ([]GrantEvent, error) {
	return o.registry.Audit(id)
}

// StatusSummary aggregates lifecycle telemetry for reporting.
func (o *GrantOrchestrator) StatusSummary() GrantStatusSummary {
	return o.registry.StatusSummary()
}

// Subscribe registers a watcher for grant events. The returned function unsubscribes the watcher.
func (o *GrantOrchestrator) Subscribe(grantID uint64, ch chan<- GrantEvent) func() {
	if ch == nil {
		return func() {}
	}
	o.mu.Lock()
	o.watchers[grantID] = append(o.watchers[grantID], ch)
	o.mu.Unlock()
	return func() {
		o.mu.Lock()
		defer o.mu.Unlock()
		watchers := o.watchers[grantID]
		for i, existing := range watchers {
			if existing == ch {
				o.watchers[grantID] = append(watchers[:i], watchers[i+1:]...)
				break
			}
		}
		if len(o.watchers[grantID]) == 0 {
			delete(o.watchers, grantID)
		}
	}
}

func (o *GrantOrchestrator) emit(grantID uint64, evt GrantEvent) {
	o.mu.RLock()
	watchers := append([]chan<- GrantEvent(nil), o.watchers[grantID]...)
	o.mu.RUnlock()
	for _, ch := range watchers {
		select {
		case ch <- evt:
		default:
			go func(out chan<- GrantEvent) {
				select {
				case out <- evt:
				default:
				}
			}(ch)
		}
	}
}

func (o *GrantOrchestrator) ensureConsensus() error {
	if o == nil || o.consensus == nil {
		return nil
	}
	if o.consensus.PoWAvailable || o.consensus.PoSAvailable || o.consensus.PoHAvailable {
		return nil
	}
	return ErrGrantConsensusOffline
}

func (o *GrantOrchestrator) verifyProof(proof *WalletProof, expected []byte) error {
	if proof == nil {
		return ErrGrantSignatureInvalid
	}
	if !bytes.Equal(proof.Message, expected) {
		return ErrGrantSignatureInvalid
	}
	if !proof.Verify() {
		return ErrGrantSignatureInvalid
	}
	addr, err := AddressFromPublicKey(proof.PublicKey)
	if err != nil {
		return err
	}
	if addr != proof.Address {
		return ErrGrantSignatureInvalid
	}
	return nil
}

func (o *GrantOrchestrator) registerVMHandlers() {
	o.vmOnce.Do(func() {
		handlers := map[string]func([]byte) ([]byte, error){
			"InitGrantEngine":               o.vmInit,
			"GrantEngine":                   o.vmEngineInfo,
			"GrantToken_Create":             o.vmCreate,
			"GrantToken_Authorize":          o.vmAuthorize,
			"GrantToken_Disburse":           o.vmDisburse,
			"GrantToken_Info":               o.vmInfo,
			"GrantToken_List":               o.vmList,
			"GrantToken_Audit":              o.vmAudit,
			"GrantToken_Status":             o.vmStatus,
			"core_syn3800_NewGrantRegistry": o.vmInit,
			"core_syn3800_CreateGrant":      o.vmCreate,
			"core_syn3800_Authorize":        o.vmAuthorize,
			"core_syn3800_Disburse":         o.vmDisburse,
			"core_syn3800_GetGrant":         o.vmInfo,
			"core_syn3800_ListGrants":       o.vmList,
			"core_syn3800_Audit":            o.vmAudit,
			"core_syn3800_StatusSummary":    o.vmStatus,
		}
		for name, handler := range handlers {
			if code, ok := synn.LookupOpcode(name); ok {
				o.vm.RegisterOpcode(code, handler)
			}
		}
	})
}

func (o *GrantOrchestrator) vmInit(_ []byte) ([]byte, error) {
	return []byte("ok"), nil
}

func (o *GrantOrchestrator) vmEngineInfo(_ []byte) ([]byte, error) {
	summary := o.StatusSummary()
	return json.Marshal(summary)
}

func (o *GrantOrchestrator) vmCreate(payload []byte) ([]byte, error) {
	var req struct {
		Beneficiary string          `json:"beneficiary"`
		Name        string          `json:"name"`
		Amount      uint64          `json:"amount"`
		Authorizers []string        `json:"authorizers"`
		Creator     walletProofJSON `json:"creator"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	creator, err := decodeWalletProofJSON(req.Creator)
	if err != nil {
		return nil, err
	}
	record, err := o.CreateGrant(context.Background(), GrantCreationRequest{
		Beneficiary: req.Beneficiary,
		Name:        req.Name,
		Amount:      req.Amount,
		Authorizers: req.Authorizers,
		Creator:     creator,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(record)
}

func (o *GrantOrchestrator) vmAuthorize(payload []byte) ([]byte, error) {
	var req struct {
		ID    uint64          `json:"id"`
		Proof walletProofJSON `json:"proof"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	proof, err := decodeWalletProofJSON(req.Proof)
	if err != nil {
		return nil, err
	}
	evt, err := o.AuthorizeGrant(context.Background(), GrantAuthorizationRequest{GrantID: req.ID, Proof: proof})
	if err != nil {
		return nil, err
	}
	return json.Marshal(evt)
}

func (o *GrantOrchestrator) vmDisburse(payload []byte) ([]byte, error) {
	var req struct {
		ID     uint64          `json:"id"`
		Amount uint64          `json:"amount"`
		Note   string          `json:"note"`
		Proof  walletProofJSON `json:"proof"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	proof, err := decodeWalletProofJSON(req.Proof)
	if err != nil {
		return nil, err
	}
	evt, err := o.ReleaseGrant(context.Background(), GrantReleaseRequest{GrantID: req.ID, Amount: req.Amount, Note: req.Note, Proof: proof})
	if err != nil {
		return nil, err
	}
	return json.Marshal(evt)
}

func (o *GrantOrchestrator) vmInfo(payload []byte) ([]byte, error) {
	var req struct {
		ID uint64 `json:"id"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	record, ok := o.GetGrant(req.ID)
	if !ok {
		return nil, ErrGrantNotFound
	}
	return json.Marshal(record)
}

func (o *GrantOrchestrator) vmList(_ []byte) ([]byte, error) {
	return json.Marshal(o.ListGrants())
}

func (o *GrantOrchestrator) vmAudit(payload []byte) ([]byte, error) {
	var req struct {
		ID uint64 `json:"id"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, err
	}
	events, err := o.Audit(req.ID)
	if err != nil {
		return nil, err
	}
	return json.Marshal(events)
}

func (o *GrantOrchestrator) vmStatus(_ []byte) ([]byte, error) {
	summary := o.StatusSummary()
	return json.Marshal(summary)
}

type walletProofJSON struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Message   string `json:"message"`
	PublicKey struct {
		X string `json:"x"`
		Y string `json:"y"`
	} `json:"public_key"`
}

func decodeWalletProofJSON(src walletProofJSON) (*WalletProof, error) {
	if src.Address == "" || src.Signature == "" || src.Message == "" || src.PublicKey.X == "" || src.PublicKey.Y == "" {
		return nil, errors.New("wallet proof incomplete")
	}
	sig, err := base64.StdEncoding.DecodeString(src.Signature)
	if err != nil {
		return nil, err
	}
	msg, err := base64.StdEncoding.DecodeString(src.Message)
	if err != nil {
		return nil, err
	}
	xBytes, err := hex.DecodeString(src.PublicKey.X)
	if err != nil {
		return nil, err
	}
	yBytes, err := hex.DecodeString(src.PublicKey.Y)
	if err != nil {
		return nil, err
	}
	pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(xBytes), Y: new(big.Int).SetBytes(yBytes)}
	proof := &WalletProof{
		Address:   src.Address,
		PublicKey: pub,
		Signature: sig,
		Message:   msg,
	}
	if !proof.Verify() {
		return nil, ErrGrantSignatureInvalid
	}
	addr, err := AddressFromPublicKey(pub)
	if err != nil {
		return nil, err
	}
	if addr != proof.Address {
		return nil, ErrGrantSignatureInvalid
	}
	return proof, nil
}
