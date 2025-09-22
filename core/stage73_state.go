package core

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// Stage73Snapshot captures the persisted state for the Stage 73 enterprise
// modules. It is written to disk so that the CLI, web interface and integration
// tests share deterministic data. The snapshot doubles as the payload exposed
// to dashboards and virtual machine helpers via the Stage73Orchestrator.
type Stage73Snapshot struct {
	Index    *stage73IndexState     `json:"index,omitempty"`
	Grants   stage73GrantSnapshot   `json:"grants"`
	Benefits stage73BenefitSnapshot `json:"benefits"`
	Charity  []stage73CharityState  `json:"charity,omitempty"`
	Legal    []stage73LegalState    `json:"legal,omitempty"`
	Utility  *stage73UtilityState   `json:"utility,omitempty"`
}

type stage73GrantSnapshot struct {
	NextID  uint64        `json:"next_id"`
	Records []GrantRecord `json:"records,omitempty"`
	Summary GrantSummary  `json:"summary"`
}

type stage73BenefitSnapshot struct {
	NextID  uint64          `json:"next_id"`
	Records []BenefitRecord `json:"records,omitempty"`
	Summary BenefitSummary  `json:"summary"`
}

type stage73IndexState struct {
	Name        string            `json:"name"`
	Symbol      string            `json:"symbol"`
	Components  []Component       `json:"components"`
	Controllers []string          `json:"controllers"`
	Audit       []IndexAuditEntry `json:"audit"`
}

type stage73CharityState struct {
	Symbol    string            `json:"symbol"`
	Purpose   string            `json:"purpose"`
	Goal      uint64            `json:"goal"`
	Raised    uint64            `json:"raised"`
	Donations map[string]uint64 `json:"donations"`
}

type stage73LegalState struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Symbol       string            `json:"symbol"`
	DocumentType string            `json:"document_type"`
	DocumentHash string            `json:"document_hash"`
	Expiry       time.Time         `json:"expiry"`
	Owner        string            `json:"owner"`
	Supply       uint64            `json:"supply"`
	Parties      []string          `json:"parties"`
	Signatures   map[string]string `json:"signatures"`
	Status       LegalTokenStatus  `json:"status"`
	Disputes     []Dispute         `json:"disputes"`
}

type stage73UtilityState struct {
	Name      string                 `json:"name"`
	Symbol    string                 `json:"symbol"`
	Owner     string                 `json:"owner"`
	Decimals  uint8                  `json:"decimals"`
	Supply    uint64                 `json:"supply"`
	Grants    map[string]ServiceTier `json:"grants"`
	Telemetry UtilityTelemetry       `json:"telemetry"`
}

// Stage73Store manages the Stage 73 modules and persists their state to a JSON
// snapshot.  It is concurrency-safe and ensures atomic writes when saving to
// disk so that crashes or abrupt process termination do not corrupt the state
// file.  The store exposes helper methods to mark the snapshot as dirty which
// are invoked by CLI commands whenever mutations occur.
type Stage73Store struct {
	mu       sync.Mutex
	path     string
	index    *SYN3700Token
	grants   *GrantRegistry
	benefits *BenefitRegistry
	charity  *SYN4200Token
	legal    *LegalTokenRegistry
	utility  *SYN500Token
	dirty    bool
	loaded   bool
}

// NewStage73Store initialises a Stage 73 store. The supplied path may be empty
// which results in an in-memory store; snapshots will only be persisted when a
// non-empty path is provided. Modules that always exist (grants, benefits,
// charity and legal registries) are initialised eagerly while optional tokens
// (SYN3700 index and SYN500 utility) remain nil until created via the CLI.
func NewStage73Store(path string) *Stage73Store {
	return &Stage73Store{
		path:     path,
		grants:   NewGrantRegistry(),
		benefits: NewBenefitRegistry(),
		charity:  NewSYN4200Token(),
		legal:    NewLegalTokenRegistry(),
	}
}

// Load hydrates the store from the on-disk snapshot. Missing files are
// tolerated so fresh environments can bootstrap without manual scaffolding. An
// error is returned when the snapshot cannot be parsed which allows callers to
// fail fast during startup.
func (s *Stage73Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.loaded {
		return nil
	}
	s.loaded = true
	if s.path == "" {
		return nil
	}
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var snap Stage73Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return fmt.Errorf("stage73: decode snapshot: %w", err)
	}
	s.applySnapshotLocked(&snap)
	return nil
}

// applySnapshotLocked rehydrates the store from the provided snapshot. The
// caller must hold s.mu.
func (s *Stage73Store) applySnapshotLocked(snap *Stage73Snapshot) {
	if snap == nil {
		return
	}
	if snap.Index != nil {
		idx := NewSYN3700Token(snap.Index.Name, snap.Index.Symbol)
		idx.mu.Lock()
		idx.components = make(map[string]*componentState, len(snap.Index.Components))
		for _, comp := range snap.Index.Components {
			c := comp
			idx.components[c.Token] = &componentState{Component: c, lastRebalance: time.Now().UTC()}
		}
		idx.controllers = make(map[string]struct{}, len(snap.Index.Controllers))
		for _, ctrl := range snap.Index.Controllers {
			if ctrl == "" {
				continue
			}
			idx.controllers[ctrl] = struct{}{}
		}
		idx.audit = append([]IndexAuditEntry(nil), snap.Index.Audit...)
		idx.mu.Unlock()
		s.index = idx
	}

	s.grants = NewGrantRegistry()
	s.grants.mu.Lock()
	s.grants.nextID = snap.Grants.NextID
	s.grants.grants = make(map[uint64]*grantState, len(snap.Grants.Records))
	for _, rec := range snap.Grants.Records {
		cp := rec
		cp.Authorizers = append([]string(nil), rec.Authorizers...)
		cp.Events = append([]GrantEvent(nil), rec.Events...)
		cp.Notes = append([]string(nil), rec.Notes...)
		st := &grantState{
			record:      cp,
			authorizers: make(map[string]struct{}, len(cp.Authorizers)),
		}
		for _, addr := range cp.Authorizers {
			st.authorizers[addr] = struct{}{}
		}
		s.grants.grants[cp.ID] = st
	}
	s.grants.mu.Unlock()

	s.benefits = NewBenefitRegistry()
	s.benefits.mu.Lock()
	s.benefits.nextID = snap.Benefits.NextID
	s.benefits.benefits = make(map[uint64]*benefitState, len(snap.Benefits.Records))
	for _, rec := range snap.Benefits.Records {
		cp := rec
		cp.Approvers = append([]string(nil), rec.Approvers...)
		cp.Events = append([]BenefitEvent(nil), rec.Events...)
		st := &benefitState{
			record:    cp,
			approvers: make(map[string]struct{}, len(cp.Approvers)),
		}
		for _, addr := range cp.Approvers {
			st.approvers[addr] = struct{}{}
		}
		s.benefits.benefits[cp.ID] = st
	}
	s.benefits.mu.Unlock()

	s.charity = NewSYN4200Token()
	s.charity.mu.Lock()
	s.charity.campaigns = make(map[string]*CharityCampaign, len(snap.Charity))
	for _, camp := range snap.Charity {
		c := camp
		cp := &CharityCampaign{
			Symbol:    c.Symbol,
			Purpose:   c.Purpose,
			Goal:      c.Goal,
			Raised:    c.Raised,
			Donations: make(map[string]uint64, len(c.Donations)),
		}
		for k, v := range c.Donations {
			cp.Donations[k] = v
		}
		s.charity.campaigns[c.Symbol] = cp
	}
	s.charity.mu.Unlock()

	s.legal = NewLegalTokenRegistry()
	s.legal.mu.Lock()
	s.legal.tokens = make(map[string]*LegalToken, len(snap.Legal))
	for _, entry := range snap.Legal {
		tok := NewLegalToken(entry.ID, entry.Name, entry.Symbol, entry.DocumentType, entry.DocumentHash, entry.Owner, entry.Expiry, entry.Supply, entry.Parties)
		tok.mu.Lock()
		tok.Status = entry.Status
		tok.Disputes = append([]Dispute(nil), entry.Disputes...)
		tok.Signatures = make(map[string]string, len(entry.Signatures))
		for party, sig := range entry.Signatures {
			tok.Signatures[party] = sig
		}
		tok.mu.Unlock()
		s.legal.tokens[entry.ID] = tok
	}
	s.legal.mu.Unlock()

	if snap.Utility != nil {
		util := NewSYN500Token(snap.Utility.Name, snap.Utility.Symbol, snap.Utility.Owner, snap.Utility.Decimals, snap.Utility.Supply)
		util.mu.Lock()
		util.Grants = make(map[string]*ServiceTier, len(snap.Utility.Grants))
		for addr, tier := range snap.Utility.Grants {
			t := tier
			util.Grants[addr] = &t
		}
		util.mu.Unlock()
		s.utility = util
	}
}

// Save persists the current snapshot when the store has been marked as dirty.
// Writes are performed atomically by renaming a temporary file into place. A
// no-op occurs when the store has not been mutated or when no path is configured.
func (s *Stage73Store) Save() error {
	s.mu.Lock()
	dirty := s.dirty
	path := s.path
	s.dirty = false
	s.mu.Unlock()

	if !dirty || path == "" {
		return nil
	}

	snap := s.Snapshot()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(path), "stage73-*.tmp")
	if err != nil {
		return err
	}
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return err
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		return err
	}
	if err := os.Rename(tmpFile.Name(), path); err != nil {
		os.Remove(tmpFile.Name())
		return err
	}
	return nil
}

// MarkDirty flags the store as having outstanding mutations. It is safe to
// call frequently; the dirty flag is cleared automatically after Save.
func (s *Stage73Store) MarkDirty() {
	s.mu.Lock()
	s.dirty = true
	s.mu.Unlock()
}

// Index returns the SYN3700 index token if one has been initialised.
func (s *Stage73Store) Index() *SYN3700Token {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.index
}

// SetIndex stores the provided SYN3700 token and marks the snapshot as dirty.
func (s *Stage73Store) SetIndex(tok *SYN3700Token) {
	s.mu.Lock()
	s.index = tok
	s.dirty = true
	s.mu.Unlock()
}

// Grants exposes the shared grant registry.
func (s *Stage73Store) Grants() *GrantRegistry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.grants == nil {
		s.grants = NewGrantRegistry()
	}
	return s.grants
}

// Benefits exposes the shared benefit registry.
func (s *Stage73Store) Benefits() *BenefitRegistry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.benefits == nil {
		s.benefits = NewBenefitRegistry()
	}
	return s.benefits
}

// Charity exposes the SYN4200 token registry.
func (s *Stage73Store) Charity() *SYN4200Token {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.charity == nil {
		s.charity = NewSYN4200Token()
	}
	return s.charity
}

// Legal exposes the SYN4700 registry.
func (s *Stage73Store) Legal() *LegalTokenRegistry {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.legal == nil {
		s.legal = NewLegalTokenRegistry()
	}
	return s.legal
}

// Utility returns the SYN500 token if present.
func (s *Stage73Store) Utility() *SYN500Token {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.utility
}

// SetUtility stores the SYN500 token.
func (s *Stage73Store) SetUtility(tok *SYN500Token) {
	s.mu.Lock()
	s.utility = tok
	s.dirty = true
	s.mu.Unlock()
}

// Snapshot generates a consistent view of the Stage 73 modules. It may be used
// for persistence, orchestration or web telemetry.
func (s *Stage73Store) Snapshot() Stage73Snapshot {
	snap := Stage73Snapshot{}
	if idx := s.Index(); idx != nil {
		idx.mu.RLock()
		comps := make([]Component, 0, len(idx.components))
		for _, st := range idx.components {
			comps = append(comps, st.Component)
		}
		sort.Slice(comps, func(i, j int) bool { return comps[i].Token < comps[j].Token })
		ctrls := make([]string, 0, len(idx.controllers))
		for addr := range idx.controllers {
			ctrls = append(ctrls, addr)
		}
		sort.Strings(ctrls)
		audit := append([]IndexAuditEntry(nil), idx.audit...)
		idx.mu.RUnlock()
		snap.Index = &stage73IndexState{
			Name:        idx.Name,
			Symbol:      idx.Symbol,
			Components:  comps,
			Controllers: ctrls,
			Audit:       audit,
		}
	}

	grants := s.Grants()
	grants.mu.RLock()
	grantRecords := make([]GrantRecord, 0, len(grants.grants))
	for _, st := range grants.grants {
		rec := st.record
		rec.Authorizers = append([]string(nil), rec.Authorizers...)
		rec.Events = append([]GrantEvent(nil), rec.Events...)
		rec.Notes = append([]string(nil), rec.Notes...)
		grantRecords = append(grantRecords, rec)
	}
	sort.Slice(grantRecords, func(i, j int) bool { return grantRecords[i].ID < grantRecords[j].ID })
	summary := grants.summaryLocked()
	nextID := grants.nextID
	grants.mu.RUnlock()
	snap.Grants = stage73GrantSnapshot{NextID: nextID, Records: grantRecords, Summary: summary}

	benefits := s.Benefits()
	benefits.mu.RLock()
	benefitRecords := make([]BenefitRecord, 0, len(benefits.benefits))
	for _, st := range benefits.benefits {
		rec := st.record
		rec.Approvers = append([]string(nil), rec.Approvers...)
		rec.Events = append([]BenefitEvent(nil), rec.Events...)
		benefitRecords = append(benefitRecords, rec)
	}
	sort.Slice(benefitRecords, func(i, j int) bool { return benefitRecords[i].ID < benefitRecords[j].ID })
	bSummary := benefits.summaryLocked()
	bNext := benefits.nextID
	benefits.mu.RUnlock()
	snap.Benefits = stage73BenefitSnapshot{NextID: bNext, Records: benefitRecords, Summary: bSummary}

	charity := s.Charity()
	charity.mu.RLock()
	charityStates := make([]stage73CharityState, 0, len(charity.campaigns))
	for _, camp := range charity.campaigns {
		cs := stage73CharityState{
			Symbol:    camp.Symbol,
			Purpose:   camp.Purpose,
			Goal:      camp.Goal,
			Raised:    camp.Raised,
			Donations: make(map[string]uint64, len(camp.Donations)),
		}
		for addr, amt := range camp.Donations {
			cs.Donations[addr] = amt
		}
		charityStates = append(charityStates, cs)
	}
	charity.mu.RUnlock()
	sort.Slice(charityStates, func(i, j int) bool { return charityStates[i].Symbol < charityStates[j].Symbol })
	snap.Charity = charityStates

	legal := s.Legal()
	legal.mu.RLock()
	legalStates := make([]stage73LegalState, 0, len(legal.tokens))
	for _, tok := range legal.tokens {
		tok.mu.RLock()
		state := stage73LegalState{
			ID:           tok.ID,
			Name:         tok.Name,
			Symbol:       tok.Symbol,
			DocumentType: tok.DocumentType,
			DocumentHash: tok.DocumentHash,
			Expiry:       tok.Expiry,
			Owner:        tok.Owner,
			Supply:       tok.Supply,
			Parties:      append([]string(nil), tok.Parties...),
			Signatures:   make(map[string]string, len(tok.Signatures)),
			Status:       tok.Status,
			Disputes:     append([]Dispute(nil), tok.Disputes...),
		}
		for party, sig := range tok.Signatures {
			state.Signatures[party] = sig
		}
		tok.mu.RUnlock()
		legalStates = append(legalStates, state)
	}
	legal.mu.RUnlock()
	sort.Slice(legalStates, func(i, j int) bool { return legalStates[i].ID < legalStates[j].ID })
	snap.Legal = legalStates

	if util := s.Utility(); util != nil {
		util.mu.RLock()
		grants := make(map[string]ServiceTier, len(util.Grants))
		for addr, tier := range util.Grants {
			grants[addr] = *tier
		}
		tele := util.Telemetry()
		state := stage73UtilityState{
			Name:      util.Name,
			Symbol:    util.Symbol,
			Owner:     util.Owner,
			Decimals:  util.Decimals,
			Supply:    util.Supply,
			Grants:    grants,
			Telemetry: tele,
		}
		util.mu.RUnlock()
		snap.Utility = &state
	}

	return snap
}

// Digest returns the SHA-256 hash of the current snapshot. It is used by the
// Stage73Orchestrator to detect drift when coordinating with consensus and the
// web interfaces.
func (s *Stage73Store) Digest() (string, error) {
	snap := s.Snapshot()
	data, err := json.Marshal(snap)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum[:]), nil
}

// GrantsAreDirty allows tests to check whether mutations occurred without
// depending on internal fields.
func (s *Stage73Store) GrantsAreDirty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.dirty
}

// helper methods operating under registry locks. They allow Snapshot to reuse
// summary logic without releasing locks early.
func (r *GrantRegistry) summaryLocked() GrantSummary {
	summary := GrantSummary{Total: len(r.grants)}
	for _, state := range r.grants {
		switch state.record.Status {
		case GrantStatusPending:
			summary.Pending++
		case GrantStatusActive:
			summary.Active++
		case GrantStatusCompleted:
			summary.Completed++
		}
	}
	return summary
}

func (r *BenefitRegistry) summaryLocked() BenefitSummary {
	summary := BenefitSummary{Total: len(r.benefits)}
	for _, state := range r.benefits {
		switch state.record.Status {
		case BenefitStatusPending:
			summary.Pending++
		case BenefitStatusClaimed:
			summary.Claimed++
		case BenefitStatusApproved:
			summary.Approved++
		}
	}
	return summary
}
