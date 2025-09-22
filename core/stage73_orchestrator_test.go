package core

import (
	"encoding/json"
	"testing"
)

type vmRecorder struct {
	handlers map[uint32]func([]byte) ([]byte, error)
}

func (r *vmRecorder) RegisterOpcode(code uint32, handler func([]byte) ([]byte, error)) {
	if r.handlers == nil {
		r.handlers = make(map[uint32]func([]byte) ([]byte, error))
	}
	r.handlers[code] = handler
}

func TestStage73OrchestratorVMBindings(t *testing.T) {
	store := NewStage73Store("")
	store.Grants().CreateGrant("beneficiary", "education", 100)
	store.MarkDirty()

	orchestrator := NewStage73Orchestrator(store, NewSynnergyConsensus())
	recorder := &vmRecorder{}
	if err := orchestrator.BindVM(recorder); err != nil {
		t.Fatalf("bind vm: %v", err)
	}
	if len(recorder.handlers) != 2 {
		t.Fatalf("expected handlers registered, got %d", len(recorder.handlers))
	}

	snapshotHandler := recorder.handlers[stage73OpcodeSnapshot]
	payload, err := snapshotHandler(nil)
	if err != nil {
		t.Fatalf("snapshot handler: %v", err)
	}
	var snap Stage73Snapshot
	if err := json.Unmarshal(payload, &snap); err != nil {
		t.Fatalf("decode snapshot: %v", err)
	}
	if snap.Grants.Summary.Total != 1 {
		t.Fatalf("unexpected snapshot summary: %+v", snap.Grants.Summary)
	}

	teleHandler := recorder.handlers[stage73OpcodeTelemetry]
	telemetryPayload, err := teleHandler(nil)
	if err != nil {
		t.Fatalf("telemetry handler: %v", err)
	}
	var telemetry map[string]any
	if err := json.Unmarshal(telemetryPayload, &telemetry); err != nil {
		t.Fatalf("decode telemetry: %v", err)
	}
	if telemetry["grants"] == nil {
		t.Fatalf("expected grants telemetry: %+v", telemetry)
	}
}

func TestStage73OrchestratorConsensus(t *testing.T) {
	store := NewStage73Store("")
	grants := store.Grants()
	id := grants.CreateGrant("beneficiary", "education", 100)
	if err := grants.AddAuthorizer(id, "wallet"); err != nil {
		t.Fatalf("add authorizer: %v", err)
	}
	if err := grants.DisburseWithActor(id, 100, "", "wallet"); err != nil {
		t.Fatalf("disburse: %v", err)
	}
	benefits := store.Benefits()
	bid := benefits.RegisterBenefit("recipient", "health", 50)
	benefits.AddApprover(bid, "wallet")
	benefits.Approve(bid, "wallet")

	consensus := NewSynnergyConsensus()
	orchestrator := NewStage73Orchestrator(store, consensus)
	orchestrator.UpdateConsensus()
	if consensus.Weights.PoW == 0 && consensus.Weights.PoS == 0 && consensus.Weights.PoH == 0 {
		t.Fatalf("consensus weights not updated")
	}
}

func TestStage73OrchestratorDigest(t *testing.T) {
	store := NewStage73Store("")
	orchestrator := NewStage73Orchestrator(store, nil)
	digest, err := orchestrator.SnapshotDigest()
	if err != nil {
		t.Fatalf("snapshot digest: %v", err)
	}
	if digest == "" {
		t.Fatalf("expected digest")
	}
	if orchestrator.LastDigest() != digest {
		t.Fatalf("last digest mismatch")
	}
}
