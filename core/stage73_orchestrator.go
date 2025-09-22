package core

import (
	"encoding/json"
	"errors"
	"sync"
)

// stage73VMAdapter captures the subset of the virtual machine interface used by
// the orchestrator. It mirrors the RegisterOpcode signature exposed by the
// SimpleVM so the runtime integration can wire the orchestrator without
// creating an import cycle.
type stage73VMAdapter interface {
	RegisterOpcode(code uint32, handler func([]byte) ([]byte, error))
}

const (
	stage73OpcodeSnapshot  uint32 = 0x730001
	stage73OpcodeTelemetry uint32 = 0x730002
)

// Stage73Orchestrator coordinates the Stage 73 enterprise modules with the
// consensus engine and virtual machine. It keeps a cached digest of the current
// snapshot so that callers can quickly detect whether downstream components need
// refreshing.
type Stage73Orchestrator struct {
	store     *Stage73Store
	consensus *SynnergyConsensus

	mu         sync.Mutex
	lastDigest string
}

// NewStage73Orchestrator initialises a Stage 73 orchestrator.
func NewStage73Orchestrator(store *Stage73Store, consensus *SynnergyConsensus) *Stage73Orchestrator {
	return &Stage73Orchestrator{store: store, consensus: consensus}
}

// BindVM registers Stage 73 inspection opcodes with the provided VM adapter so
// that smart contracts and dashboards can access the persisted state without
// duplicating file-system logic.
func (o *Stage73Orchestrator) BindVM(vm stage73VMAdapter) error {
	if vm == nil {
		return errors.New("stage73 orchestrator: vm required")
	}
	vm.RegisterOpcode(stage73OpcodeSnapshot, o.vmSnapshot)
	vm.RegisterOpcode(stage73OpcodeTelemetry, o.vmTelemetry)
	return nil
}

func (o *Stage73Orchestrator) vmSnapshot(_ []byte) ([]byte, error) {
	snap := o.store.Snapshot()
	return json.Marshal(snap)
}

func (o *Stage73Orchestrator) vmTelemetry(_ []byte) ([]byte, error) {
	snap := o.store.Snapshot()
	telemetry := map[string]any{
		"grants":   snap.Grants.Summary,
		"benefits": snap.Benefits.Summary,
	}
	if snap.Utility != nil {
		telemetry["utility"] = snap.Utility.Telemetry
	}
	return json.Marshal(telemetry)
}

// UpdateConsensus adjusts consensus weights based on Stage 73 telemetry. Grant
// and benefit activity contribute to network demand while approved grants and
// benefits are treated as a proxy for stake concentration.
func (o *Stage73Orchestrator) UpdateConsensus() {
	if o.consensus == nil {
		return
	}
	snap := o.store.Snapshot()
	demand := float64(snap.Grants.Summary.Total + snap.Benefits.Summary.Total)
	stake := float64(snap.Grants.Summary.Completed + snap.Benefits.Summary.Approved)
	if demand == 0 && stake == 0 {
		return
	}
	if demand > o.consensus.Dmax {
		o.consensus.Dmax = demand
	}
	if stake > o.consensus.Smax {
		o.consensus.Smax = stake
	}
	o.consensus.AdjustWeights(demand, stake)
}

// SnapshotDigest returns the cached digest of the Stage 73 snapshot. When the
// state changes a new digest is computed so callers can detect drift without
// reloading the full JSON payload.
func (o *Stage73Orchestrator) SnapshotDigest() (string, error) {
	digest, err := o.store.Digest()
	if err != nil {
		return "", err
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	o.lastDigest = digest
	return digest, nil
}

// LastDigest returns the previously calculated digest.
func (o *Stage73Orchestrator) LastDigest() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	return o.lastDigest
}
