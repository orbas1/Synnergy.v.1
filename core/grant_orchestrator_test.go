package core

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"testing"

	synn "synnergy"
)

func TestGrantOrchestratorLifecycle(t *testing.T) {
	registry := NewGrantRegistry()
	vm := synn.NewSimpleVM()
	consensus := NewSynnergyConsensus()
	authorities := NewAuthorityNodeRegistry()
	orchestrator, err := NewGrantOrchestrator(registry, vm, consensus, authorities)
	if err != nil {
		t.Fatalf("new orchestrator: %v", err)
	}

	creator, err := NewWallet()
	if err != nil {
		t.Fatalf("new wallet: %v", err)
	}
	createProof, err := NewWalletProof(creator, GrantCreateMessage("alice", "research", 100))
	if err != nil {
		t.Fatalf("proof: %v", err)
	}
	record, err := orchestrator.CreateGrant(context.Background(), GrantCreationRequest{
		Beneficiary: "alice",
		Name:        "research",
		Amount:      100,
		Creator:     createProof,
	})
	if err != nil {
		t.Fatalf("create grant: %v", err)
	}
	if record.ID == 0 {
		t.Fatalf("expected non-zero id")
	}

	eventsCh := make(chan GrantEvent, 4)
	cancel := orchestrator.Subscribe(record.ID, eventsCh)
	defer cancel()

	authorizer, err := NewWallet()
	if err != nil {
		t.Fatalf("authorizer wallet: %v", err)
	}
	authProof, err := NewWalletProof(authorizer, GrantAuthorizeMessage(record.ID))
	if err != nil {
		t.Fatalf("auth proof: %v", err)
	}
	if _, err := orchestrator.AuthorizeGrant(context.Background(), GrantAuthorizationRequest{GrantID: record.ID, Proof: authProof}); err != nil {
		t.Fatalf("authorize: %v", err)
	}
	select {
	case evt := <-eventsCh:
		if evt.Type != GrantEventAuthorized {
			t.Fatalf("expected authorized event, got %v", evt.Type)
		}
	default:
		t.Fatalf("expected event emission")
	}

	rogue, err := NewWallet()
	if err != nil {
		t.Fatalf("rogue wallet: %v", err)
	}
	rogueProof, err := NewWalletProof(rogue, GrantReleaseMessage(record.ID, 50, "phase1"))
	if err != nil {
		t.Fatalf("rogue proof: %v", err)
	}
	if _, err := orchestrator.ReleaseGrant(context.Background(), GrantReleaseRequest{GrantID: record.ID, Amount: 50, Note: "phase1", Proof: rogueProof}); err == nil {
		t.Fatal("expected unauthorized release to fail")
	}

	releaseProof, err := NewWalletProof(authorizer, GrantReleaseMessage(record.ID, 50, "phase1"))
	if err != nil {
		t.Fatalf("release proof: %v", err)
	}
	if _, err := orchestrator.ReleaseGrant(context.Background(), GrantReleaseRequest{GrantID: record.ID, Amount: 50, Note: "phase1", Proof: releaseProof}); err != nil {
		t.Fatalf("release: %v", err)
	}
	summary := orchestrator.StatusSummary()
	if summary.Active != 1 || summary.Total != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}

	// Exercise VM bindings for info and status opcodes.
	infoPayload, _ := json.Marshal(struct {
		ID uint64 `json:"id"`
	}{ID: record.ID})
	infoCode, ok := synn.LookupOpcode("core_syn3800_GetGrant")
	if !ok {
		t.Fatalf("missing GrantToken_Info opcode")
	}
	wasm := []byte{byte(infoCode >> 16), byte(infoCode >> 8), byte(infoCode)}
	out, _, err := vm.Execute(wasm, "", infoPayload, 1000)
	if err != nil {
		t.Fatalf("vm info execute: %v", err)
	}
	var fetched GrantRecord
	if err := json.Unmarshal(out, &fetched); err != nil {
		t.Fatalf("decode vm info: %v", err)
	}
	if fetched.ID != record.ID {
		t.Fatalf("vm returned wrong record: %v", fetched.ID)
	}

	statusCode, ok := synn.LookupOpcode("core_syn3800_StatusSummary")
	if !ok {
		t.Fatalf("missing GrantToken_Status opcode")
	}
	wasmStatus := []byte{byte(statusCode >> 16), byte(statusCode >> 8), byte(statusCode)}
	statusOut, _, err := vm.Execute(wasmStatus, "", nil, 1000)
	if err != nil {
		t.Fatalf("vm status execute: %v", err)
	}
	var vmSummary GrantStatusSummary
	if err := json.Unmarshal(statusOut, &vmSummary); err != nil {
		t.Fatalf("decode vm status: %v", err)
	}
	if vmSummary.Active != 1 {
		t.Fatalf("vm summary incorrect: %+v", vmSummary)
	}
}

func TestDecodeWalletProofJSON(t *testing.T) {
	wallet, err := NewWallet()
	if err != nil {
		t.Fatalf("wallet: %v", err)
	}
	msg := GrantAuthorizeMessage(42)
	proof, err := NewWalletProof(wallet, msg)
	if err != nil {
		t.Fatalf("new proof: %v", err)
	}
	proofJSON := walletProofJSON{
		Address:   proof.Address,
		Signature: base64.StdEncoding.EncodeToString(proof.Signature),
		Message:   base64.StdEncoding.EncodeToString(proof.Message),
	}
	if proof.PublicKey == nil {
		t.Fatal("expected public key")
	}
	proofJSON.PublicKey.X = hex.EncodeToString(proof.PublicKey.X.Bytes())
	proofJSON.PublicKey.Y = hex.EncodeToString(proof.PublicKey.Y.Bytes())
	decoded, err := decodeWalletProofJSON(proofJSON)
	if err != nil {
		t.Fatalf("decode proof: %v", err)
	}
	if !decoded.Verify() {
		t.Fatal("decoded proof failed verification")
	}
}
