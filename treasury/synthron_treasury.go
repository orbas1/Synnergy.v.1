package treasury

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	synn "synnergy"
	"synnergy/core"
)

var (
	treasuryOnce sync.Once
	treasuryInst *SynthronTreasury
	treasuryErr  error
)

var synthronTreasuryGas = map[string]uint64{
	"SynthronCoin_Issue":             90,
	"SynthronCoin_Transfer":          55,
	"SynthronCoin_Burn":              45,
	"SynthronCoin_Reconcile":         110,
	"SynthronCoin_AuthoritySeal":     70,
	"SynthronCoin_ConsensusBridge":   85,
	"SynthronCoin_Telemetry":         25,
	"SynthronCoin_AuthorizeOperator": 40,
	"SynthronCoin_RevokeOperator":    35,
}

var synthronTreasuryDescriptions = map[string]string{
	"SynthronCoin_Issue":             "Mint Synthron to an audited account with dual ledger + VM hooks",
	"SynthronCoin_Transfer":          "Transfer Synthron with wallet-signed verification",
	"SynthronCoin_Burn":              "Retire Synthron supply with ledger burn semantics",
	"SynthronCoin_Reconcile":         "Reconcile treasury state across ledger, VM and gas table",
	"SynthronCoin_AuthoritySeal":     "Register or reseal an authority node linked to the treasury",
	"SynthronCoin_ConsensusBridge":   "Authorise cross-consensus relays for treasury oversight",
	"SynthronCoin_Telemetry":         "Emit telemetry snapshots for CLI and web dashboards",
	"SynthronCoin_AuthorizeOperator": "Grant treasury operator permissions with dual-signature auditing",
	"SynthronCoin_RevokeOperator":    "Revoke operator access while preserving root guardians",
}

type operatorContextKey struct{}

// ErrUnauthorizedOperator indicates the caller does not have permission to execute
// a privileged treasury action.
var ErrUnauthorizedOperator = errors.New("treasury: unauthorized operator")

// ErrOperatorNotFound indicates an operator record does not exist.
var ErrOperatorNotFound = errors.New("treasury: operator not found")

// ErrProtectedOperator is returned when attempting to revoke the built-in
// treasury guardian operator.
var ErrProtectedOperator = errors.New("treasury: protected operator")

// WithOperator annotates the context with an operator address so privileged
// treasury methods can enforce permissioned access. When omitted, the
// orchestrator treats requests as originating from the root treasury wallet.
func WithOperator(ctx context.Context, operator string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, operatorContextKey{}, operator)
}

// SynthronCoinEvent captures significant treasury actions for audit trails and
// subscription feeds consumed by the CLI and function web.
type SynthronCoinEvent struct {
	Timestamp time.Time         `json:"timestamp"`
	Type      string            `json:"type"`
	Amount    uint64            `json:"amount"`
	From      string            `json:"from,omitempty"`
	To        string            `json:"to,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Digest    string            `json:"digest,omitempty"`
	Signature string            `json:"signature,omitempty"`
}

// SynthronTreasuryTelemetry summarises treasury health for CLI, tests and web
// dashboards.
type SynthronTreasuryTelemetry struct {
	Timestamp         time.Time           `json:"timestamp"`
	Wallet            string              `json:"wallet"`
	VMMode            string              `json:"vmMode"`
	VMRunning         bool                `json:"vmRunning"`
	VMConcurrency     int                 `json:"vmConcurrency"`
	Minted            uint64              `json:"minted"`
	Burned            uint64              `json:"burned"`
	Circulating       uint64              `json:"circulating"`
	LedgerHeight      int                 `json:"ledgerHeight"`
	ConsensusNetworks int                 `json:"consensusNetworks"`
	AuthorityNodes    int                 `json:"authorityNodes"`
	GasCoverage       map[string]uint64   `json:"gasCoverage"`
	MissingOpcodes    []string            `json:"missingOpcodes"`
	InsertedOpcodes   []string            `json:"insertedOpcodes,omitempty"`
	Operators         []string            `json:"operators"`
	Health            TreasuryHealth      `json:"health"`
	AuditTrail        []SynthronCoinEvent `json:"auditTrail,omitempty"`
}

// TreasuryHealth surfaces subsystem status so dashboards and automation can
// detect degraded components without parsing logs.
type TreasuryHealth struct {
	VM          string `json:"vm"`
	Ledger      string `json:"ledger"`
	Wallet      string `json:"wallet"`
	Consensus   string `json:"consensus"`
	Authorities string `json:"authorities"`
}

// SynthronTreasury orchestrates the economic controls for the Synthron coin. It
// wires together the ledger, virtual machine, consensus registry, wallet and
// authority node registry so that enterprise operators have a single entry point
// for monetary policy actions.
type SynthronTreasury struct {
	mu          sync.RWMutex
	ledger      *core.Ledger
	vm          *core.SimpleVM
	wallet      *core.Wallet
	consensus   *core.ConsensusNetworkManager
	authorities *core.AuthorityNodeRegistry
	minted      uint64
	burned      uint64
	nonce       uint64
	gas         map[string]uint64
	inserted    []string
	events      []chan SynthronCoinEvent
	telemetry   SynthronTreasuryTelemetry
	operators   map[string]struct{}
	auditLog    []SynthronCoinEvent
	lastDigest  [32]byte
	hasDigest   bool
}

const treasuryAuditTrailLimit = 512

// DefaultSynthronTreasury returns the singleton treasury instance used by the
// CLI and function web. The treasury is initialised on first use to keep start-up
// lightweight for tests and tooling that do not require economic orchestration.
func DefaultSynthronTreasury(ctx context.Context) (*SynthronTreasury, error) {
	treasuryOnce.Do(func() {
		treasuryInst, treasuryErr = NewSynthronTreasury(ctx)
	})
	return treasuryInst, treasuryErr
}

// NewSynthronTreasury boots a new treasury instance. It ensures gas schedules
// are aligned, the heavy VM is running, consensus relayers are authorised and the
// treasury wallet is part of the authority registry.
func NewSynthronTreasury(ctx context.Context) (*SynthronTreasury, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	vm := core.NewSimpleVM(core.VMHeavy)
	if err := vm.Start(); err != nil {
		return nil, fmt.Errorf("start vm: %w", err)
	}

	wallet, err := core.NewWallet()
	if err != nil {
		return nil, fmt.Errorf("wallet init: %w", err)
	}

	ledger := core.NewLedger()
	ledger.Credit(wallet.Address, core.GenesisAllocation)

	consensus := core.NewConsensusNetworkManager()
	consensus.AuthorizeRelayer(wallet.Address)
	if _, err := consensus.RegisterNetwork("pow", "pos", wallet.Address); err != nil {
		return nil, fmt.Errorf("register consensus bridge: %w", err)
	}

	authorities := core.NewAuthorityNodeRegistry()
	if _, err := authorities.Register(wallet.Address, "treasury"); err != nil {
		return nil, fmt.Errorf("register treasury authority: %w", err)
	}

	inserted, err := synn.EnsureGasSchedule(synthronTreasuryGas)
	if err != nil {
		return nil, fmt.Errorf("ensure treasury gas schedule: %w", err)
	}
	for name, cost := range synthronTreasuryGas {
		desc := synthronTreasuryDescriptions[name]
		if err := synn.RegisterGasMetadata(name, cost, "treasury", desc); err != nil {
			return nil, fmt.Errorf("register gas metadata %s: %w", name, err)
		}
	}

	treasury := &SynthronTreasury{
		ledger:      ledger,
		vm:          vm,
		wallet:      wallet,
		consensus:   consensus,
		authorities: authorities,
		minted:      core.GenesisAllocation,
		gas:         synthronTreasuryGas,
		inserted:    inserted,
		operators:   map[string]struct{}{wallet.Address: {}},
	}

	vm.RegisterCallHandler("SynthronCoin_Issue", func() error {
		_, err := treasury.Issue(context.Background(), wallet.Address, core.InitialBlockReward)
		return err
	})
	vm.RegisterCallHandler("SynthronCoin_Burn", func() error {
		// Burn a minimal unit to keep telemetry active without exhausting balances.
		return treasury.Burn(context.Background(), wallet.Address, 1)
	})
	vm.RegisterCallHandler("SynthronCoin_Reconcile", func() error {
		treasury.Diagnostics(context.Background())
		return nil
	})
	vm.RegisterCallHandler("SynthronCoin_Telemetry", func() error {
		treasury.Diagnostics(context.Background())
		return nil
	})
	vm.RegisterCallHandler("SynthronCoin_AuthorizeOperator", func() error {
		return treasury.AuthorizeOperator(context.Background(), wallet.Address)
	})
	vm.RegisterCallHandler("SynthronCoin_RevokeOperator", func() error {
		return treasury.RevokeOperator(context.Background(), wallet.Address)
	})

	treasury.refreshTelemetryLocked()
	return treasury, nil
}

// Issue mints Synthron to the specified account. It updates ledger balances,
// tracks supply counters and emits an event for observers.
func (t *SynthronTreasury) Issue(ctx context.Context, addr string, amount uint64) (uint64, error) {
	if addr == "" {
		return 0, errors.New("recipient required")
	}
	if amount == 0 {
		return 0, errors.New("amount must be > 0")
	}

	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return 0, ErrUnauthorizedOperator
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if err := t.ensureVM(); err != nil {
		return 0, err
	}

	if math.MaxUint64-t.minted < amount {
		return 0, errors.New("mint overflow")
	}

	t.ledger.Credit(addr, amount)
	t.minted += amount
	cir := t.refreshTelemetryLocked().Circulating
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "issue",
		Amount:    amount,
		To:        addr,
	})
	return cir, nil
}

// Transfer moves Synthron using the provided wallet. When wallet is nil the
// treasury wallet executes the transfer. Signatures are verified before applying
// the transaction to the ledger.
func (t *SynthronTreasury) Transfer(ctx context.Context, wallet *core.Wallet, to string, amount, fee uint64) error {
	if to == "" {
		return errors.New("destination required")
	}
	if amount == 0 {
		return errors.New("amount must be > 0")
	}
	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return ErrUnauthorizedOperator
	}
	if wallet == nil {
		wallet = t.wallet
	}

	if err := t.ensureVM(); err != nil {
		return err
	}

	t.mu.Lock()
	t.nonce++
	nonce := t.nonce
	t.mu.Unlock()

	tx := core.NewTransaction(wallet.Address, to, amount, fee, nonce)
	if _, err := wallet.Sign(tx); err != nil {
		return fmt.Errorf("sign tx: %w", err)
	}
	if !core.VerifySignature(tx, tx.Signature, &wallet.PublicKey) {
		return errors.New("signature verification failed")
	}
	if err := t.ledger.ApplyTransaction(tx); err != nil {
		return err
	}

	t.mu.Lock()
	t.refreshTelemetryLocked()
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "transfer",
		Amount:    amount,
		From:      wallet.Address,
		To:        to,
	})
	t.mu.Unlock()
	return nil
}

// Burn retires Synthron supply from the specified account.
func (t *SynthronTreasury) Burn(ctx context.Context, addr string, amount uint64) error {
	if amount == 0 {
		return errors.New("amount must be > 0")
	}

	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return ErrUnauthorizedOperator
	}

	if err := t.ensureVM(); err != nil {
		return err
	}

	if err := t.ledger.Burn(addr, amount); err != nil {
		return err
	}

	t.mu.Lock()
	if t.burned > math.MaxUint64-amount {
		t.mu.Unlock()
		return errors.New("burn overflow")
	}
	t.burned += amount
	t.refreshTelemetryLocked()
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "burn",
		Amount:    amount,
		From:      addr,
	})
	t.mu.Unlock()
	return nil
}

// RegisterAuthority adds a new authority node to the registry under the provided
// role. When addr is empty the treasury wallet is reused.
func (t *SynthronTreasury) RegisterAuthority(ctx context.Context, addr, role string) (*core.AuthorityNode, error) {
	if addr == "" {
		addr = t.wallet.Address
	}
	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return nil, ErrUnauthorizedOperator
	}
	if err := t.ensureVM(); err != nil {
		return nil, err
	}

	node, err := t.authorities.Register(addr, role)
	if err != nil {
		return nil, err
	}
	t.mu.Lock()
	t.refreshTelemetryLocked()
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "authority", Amount: 0, To: addr,
		Metadata: map[string]string{"role": role},
	})
	t.mu.Unlock()
	return node, nil
}

// RegisterConsensusLink provisions a consensus bridge for treasury monitoring.
func (t *SynthronTreasury) RegisterConsensusLink(ctx context.Context, source, target string) (int, error) {
	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return 0, ErrUnauthorizedOperator
	}
	if err := t.ensureVM(); err != nil {
		return 0, err
	}

	id, err := t.consensus.RegisterNetwork(source, target, t.wallet.Address)
	if err != nil {
		return 0, err
	}
	t.mu.Lock()
	t.refreshTelemetryLocked()
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "consensus",
		Metadata:  map[string]string{"source": source, "target": target, "id": fmt.Sprintf("%d", id)},
	})
	t.mu.Unlock()
	return id, nil
}

// Diagnostics returns cached telemetry or recomputes it if the previous snapshot
// is older than half a second. The behaviour matches CLI and web polling needs
// without adding additional locks to hot paths.
func (t *SynthronTreasury) Diagnostics(ctx context.Context) SynthronTreasuryTelemetry {
	t.mu.RLock()
	diag := t.telemetry
	t.mu.RUnlock()
	if time.Since(diag.Timestamp) < 500*time.Millisecond {
		return diag
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	return t.refreshTelemetryLocked()
}

// SubscribeEvents returns a buffered channel that receives future treasury
// events. Callers should drain the channel to avoid missing updates.
func (t *SynthronTreasury) SubscribeEvents() <-chan SynthronCoinEvent {
	ch := make(chan SynthronCoinEvent, 16)
	t.mu.Lock()
	t.events = append(t.events, ch)
	t.mu.Unlock()
	return ch
}

// VirtualMachine exposes the treasury-backed VM instance for integration tests
// and monitoring.
func (t *SynthronTreasury) VirtualMachine() *core.SimpleVM {
	return t.vm
}

// AuthorizeOperator grants treasury privileges to the provided address. Only
// existing operators can add new entries.
func (t *SynthronTreasury) AuthorizeOperator(ctx context.Context, addr string) error {
	if addr == "" {
		return errors.New("address required")
	}
	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return ErrUnauthorizedOperator
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.operators == nil {
		t.operators = make(map[string]struct{})
	}
	if _, ok := t.operators[addr]; ok {
		return nil
	}
	t.operators[addr] = struct{}{}
	t.refreshTelemetryLocked()
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "authorize",
		To:        addr,
	})
	return nil
}

// RevokeOperator removes operator privileges for the provided address while
// protecting the root treasury guardian.
func (t *SynthronTreasury) RevokeOperator(ctx context.Context, addr string) error {
	if addr == "" {
		return errors.New("address required")
	}
	operator := t.resolveOperator(ctx)
	if !t.isOperator(operator) {
		return ErrUnauthorizedOperator
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, protected := t.operators[t.wallet.Address]; protected && addr == t.wallet.Address {
		return ErrProtectedOperator
	}
	if _, ok := t.operators[addr]; !ok {
		return ErrOperatorNotFound
	}
	delete(t.operators, addr)
	t.refreshTelemetryLocked()
	t.recordEventLocked(SynthronCoinEvent{
		Timestamp: time.Now().UTC(),
		Type:      "revoke",
		From:      addr,
	})
	return nil
}

func (t *SynthronTreasury) refreshTelemetryLocked() SynthronTreasuryTelemetry {
	height, _ := t.ledger.Head()
	coverage := make(map[string]uint64, len(t.gas))
	names := make([]string, 0, len(t.gas))
	for name := range t.gas {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		coverage[name] = synn.GasCost(name)
	}

	missing := make([]string, 0)
	for _, name := range names {
		if synn.SNVMOpcodeByName(name) == 0 {
			missing = append(missing, name)
		}
	}

	operators := make([]string, 0, len(t.operators))
	for addr := range t.operators {
		operators = append(operators, addr)
	}
	sort.Strings(operators)

	audit := make([]SynthronCoinEvent, len(t.auditLog))
	for i, evt := range t.auditLog {
		audit[i] = cloneEvent(evt)
	}

	diag := SynthronTreasuryTelemetry{
		Timestamp:         time.Now().UTC(),
		Wallet:            t.wallet.Address,
		VMMode:            t.vm.Mode().String(),
		VMRunning:         t.vm.Status(),
		VMConcurrency:     t.vm.Concurrency(),
		Minted:            t.minted,
		Burned:            t.burned,
		Circulating:       t.currentCirculating(),
		LedgerHeight:      height,
		ConsensusNetworks: len(t.consensus.ListNetworks()),
		AuthorityNodes:    len(t.authorities.List()),
		GasCoverage:       coverage,
		MissingOpcodes:    missing,
		InsertedOpcodes:   append([]string(nil), t.inserted...),
		Operators:         operators,
		Health:            t.computeHealthLocked(),
		AuditTrail:        audit,
	}
	t.telemetry = diag
	return diag
}

func (t *SynthronTreasury) currentCirculating() uint64 {
	if t.minted < t.burned {
		return 0
	}
	return t.minted - t.burned
}

func (t *SynthronTreasury) recordEventLocked(evt SynthronCoinEvent) {
	sealed := t.sealEventLocked(evt)
	if len(t.events) == 0 {
		return
	}
	for _, ch := range t.events {
		select {
		case ch <- cloneEvent(sealed):
		default:
		}
	}
}

func (t *SynthronTreasury) ensureVM() error {
	if t.vm.Status() {
		return nil
	}
	var lastErr error
	for i := 0; i < 3; i++ {
		if err := t.vm.Start(); err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		return nil
	}
	return fmt.Errorf("ensure vm: %w", lastErr)
}

func (t *SynthronTreasury) resolveOperator(ctx context.Context) string {
	if ctx != nil {
		if op, ok := ctx.Value(operatorContextKey{}).(string); ok && op != "" {
			return op
		}
	}
	return t.wallet.Address
}

func (t *SynthronTreasury) isOperator(addr string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.operators) == 0 {
		return addr == t.wallet.Address
	}
	_, ok := t.operators[addr]
	return ok
}

// AuditTrail returns the signed event history maintained by the treasury. The
// slice is ordered chronologically and capped to the most recent entries to
// guarantee bounded memory usage.
func (t *SynthronTreasury) AuditTrail() []SynthronCoinEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]SynthronCoinEvent, len(t.auditLog))
	for i, evt := range t.auditLog {
		out[i] = cloneEvent(evt)
	}
	return out
}

func (t *SynthronTreasury) computeHealthLocked() TreasuryHealth {
	status := TreasuryHealth{
		VM:          "stopped",
		Ledger:      "uninitialised",
		Wallet:      "uninitialised",
		Consensus:   "degraded",
		Authorities: "degraded",
	}
	if t.vm != nil {
		if t.vm.Status() {
			status.VM = "ok"
		} else {
			status.VM = "recovering"
		}
	}
	if t.ledger != nil {
		status.Ledger = "ok"
	}
	if t.wallet != nil {
		status.Wallet = "ok"
	}
	if t.consensus != nil {
		if len(t.consensus.ListNetworks()) > 0 {
			status.Consensus = "ok"
		}
	}
	if t.authorities != nil {
		if len(t.authorities.List()) > 0 {
			status.Authorities = "ok"
		}
	}
	return status
}

func (t *SynthronTreasury) sealEventLocked(evt SynthronCoinEvent) SynthronCoinEvent {
	evt.Metadata = cloneMetadata(evt.Metadata)
	payload := canonicalEventPayload(evt)
	base := []byte(payload)
	if t.hasDigest {
		base = append(base, t.lastDigest[:]...)
	}
	digest := sha256.Sum256(base)
	evt.Digest = hex.EncodeToString(digest[:])
	if sig, err := t.wallet.SignMessage(base); err == nil {
		evt.Signature = hex.EncodeToString(sig)
	} else {
		if evt.Metadata == nil {
			evt.Metadata = make(map[string]string)
		}
		evt.Metadata["signature_error"] = err.Error()
	}
	t.lastDigest = digest
	t.hasDigest = true
	t.auditLog = append(t.auditLog, cloneEvent(evt))
	if len(t.auditLog) > treasuryAuditTrailLimit {
		trimmed := make([]SynthronCoinEvent, treasuryAuditTrailLimit)
		copy(trimmed, t.auditLog[len(t.auditLog)-treasuryAuditTrailLimit:])
		t.auditLog = trimmed
	}
	return evt
}

func canonicalEventPayload(evt SynthronCoinEvent) string {
	builder := strings.Builder{}
	builder.WriteString(evt.Timestamp.UTC().Format(time.RFC3339Nano))
	builder.WriteString("|")
	builder.WriteString(evt.Type)
	builder.WriteString("|")
	builder.WriteString(strconv.FormatUint(evt.Amount, 10))
	builder.WriteString("|")
	builder.WriteString(evt.From)
	builder.WriteString("|")
	builder.WriteString(evt.To)
	if len(evt.Metadata) > 0 {
		builder.WriteString("|")
		keys := make([]string, 0, len(evt.Metadata))
		for k := range evt.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		entries := make([]string, 0, len(keys))
		for _, k := range keys {
			entries = append(entries, fmt.Sprintf("%s=%s", k, evt.Metadata[k]))
		}
		builder.WriteString(strings.Join(entries, ";"))
	}
	return builder.String()
}

func cloneMetadata(meta map[string]string) map[string]string {
	if len(meta) == 0 {
		return nil
	}
	out := make(map[string]string, len(meta))
	for k, v := range meta {
		out[k] = v
	}
	return out
}

func cloneEvent(evt SynthronCoinEvent) SynthronCoinEvent {
	evt.Metadata = cloneMetadata(evt.Metadata)
	return evt
}

// SynthronTreasurySummary renders a short, human-readable summary of the current
// treasury telemetry. It is used by the CLI to present concise diagnostics.
func SynthronTreasurySummary(diag SynthronTreasuryTelemetry) string {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "wallet: %s\n", diag.Wallet)
	fmt.Fprintf(&builder, "vm: %s running=%t concurrency=%d\n", diag.VMMode, diag.VMRunning, diag.VMConcurrency)
	fmt.Fprintf(&builder, "minted: %d burned: %d circulating: %d\n", diag.Minted, diag.Burned, diag.Circulating)
	fmt.Fprintf(&builder, "ledger height: %d consensus networks: %d authority nodes: %d\n", diag.LedgerHeight, diag.ConsensusNetworks, diag.AuthorityNodes)
	fmt.Fprintf(&builder, "operators: %d\n", len(diag.Operators))
	fmt.Fprintf(&builder, "audit trail entries: %d\n", len(diag.AuditTrail))
	if len(diag.MissingOpcodes) == 0 {
		builder.WriteString("opcode catalogue: complete\n")
	} else {
		fmt.Fprintf(&builder, "missing opcodes: %s\n", strings.Join(diag.MissingOpcodes, ", "))
	}
	return builder.String()
}
